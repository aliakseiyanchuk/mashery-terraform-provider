package mashery

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryService() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		ReadContext:   serviceRead,
		CreateContext: serviceCreate,
		UpdateContext: serviceUpdate,
		DeleteContext: serviceDelete,
		// Schema
		Schema: ServiceSchema,
		// Importer by ID
		Importer: &schema.ResourceImporter{
			StateContext: importMasheryService,
		},
	}
}

func importMasheryService(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	mashV3Cl := m.(v3client.Client)

	if service, err := mashV3Cl.GetService(ctx, d.Id()); err != nil {
		return []*schema.ResourceData{}, errwrap.Wrapf("Failed to import this service: {{err}}", err)
	} else if service == nil {
		return []*schema.ResourceData{}, errors.New("No such service")
	} else {
		V3ServiceToTerraform(service, d)

		roleDiags := serviceReadRoles(ctx, d, mashV3Cl)
		if roleDiags.HasError() {
			doLogf("FAILURE while trying to read roles during import of the service")
			err := errors.New(fmt.Sprintf("service roles were not read successfully with %d diagnostic messages returned", len(roleDiags)))
			doLogJson("Returned diagnostics", roleDiags)
			return nil, err

		}
		return []*schema.ResourceData{d}, nil
	}
}

func serviceReadRoles(ctx context.Context, d *schema.ResourceData, v3cl v3client.Client) diag.Diagnostics {
	if rv, err := v3cl.GetServiceRoles(ctx, d.Id()); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error returned while retrieving roles",
			Detail:   err.Error(),
		}}
	} else {
		if rv != nil {
			return V3ServiceRolesToTerraform(rv, d)
		} else {
			return diag.Diagnostics{}
		}
	}
}

func serviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mashSvcUpsert := V3ServiceUpsertable(d, true, true)

	doLogJson("Will attempt to create new service with this upsertable", mashSvcUpsert)

	mashV3Cl := m.(v3client.Client)
	if rv, err := mashV3Cl.CreateService(ctx, mashSvcUpsert); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(rv.Id)

		opDiagnostic := V3ServiceToTerraform(rv, d)

		// After the service has been created, portal access groups need to be pushed. Otherwise default
		// pre-populated list needs to be read.
		roleDiags := trySetServiceRoles(ctx, d, mashV3Cl)
		if len(roleDiags) > 0 {
			opDiagnostic = append(roleDiags)
		}

		return opDiagnostic
	}
}

func trySetServiceRoles(ctx context.Context, d *schema.ResourceData, mashV3Cl v3client.Client) diag.Diagnostics {
	opDiagnostic := diag.Diagnostics{}
	if getSetLength(d.Get(MashSvcInteractiveDocsRoles)) > 0 {
		roles := V3ServiceRolePermissionUpsertable(d)
		doLogJson("Will attempt to set service roles with this upsertable", roles)

		err := mashV3Cl.SetServiceRoles(ctx, d.Id(), roles)
		if err != nil {
			doLogf("Returned error: %s", err)
			opDiagnostic = append(opDiagnostic, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to set service roles after service creation",
				Detail:   err.Error(),
			})
		}
	}

	return opDiagnostic
}

func serviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mashV3Cl := m.(v3client.Client)
	if rv, err := mashV3Cl.GetService(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	} else {
		servDiags := V3ServiceToTerraform(rv, d)
		rolesDiags := serviceReadRoles(ctx, d, mashV3Cl)

		return append(servDiags, rolesDiags...)
	}
}

func serviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mashV3Cl := m.(v3client.Client)
	updateDiagnostic := diag.Diagnostics{}

	if MashSvcHasDirectUpsertableModifications(d) {
		mashServ := V3ServiceUpsertable(d, false, false)

		if rv, err := mashV3Cl.UpdateService(ctx, mashServ); err != nil {
			return diag.FromErr(err)
		} else {
			upd := V3ServiceToTerraform(rv, d)
			if len(upd) > 0 {
				updateDiagnostic = append(updateDiagnostic, upd...)
			}
		}
	}

	if d.HasChange(MashSvcOAuth) {
		curSet := d.Get(MashSvcOAuth)
		if curSet == nil {
			if err := mashV3Cl.DeleteServiceOAuthSecurityProfile(ctx, d.Id()); err != nil {
				updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "could not delete oauth security profile",
					Detail:   err.Error(),
				})
			} else {
				// OAuth profile has been deleted.
				clearDiags := ClearServiceOAuthProfile(d)
				if len(clearDiags) > 0 {
					updateDiagnostic = append(updateDiagnostic, clearDiags...)
				}
			}
		} else {
			requestedProfile := V3SecurityProfileUpsertable(d)
			if actualOAuth, err := mashV3Cl.UpdateServiceOAuthSecurityProfile(ctx, d.Id(), *requestedProfile.OAuth); err != nil {
				updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "could not update oauth security profile",
					Detail:   err.Error(),
				})
			} else {
				oauthSave := V3ServiceOAuthProfileToTerraform(actualOAuth, d)
				if len(oauthSave) > 0 {
					updateDiagnostic = append(updateDiagnostic, oauthSave...)
				}
			}
		}
	}

	if d.HasChange(MashSvcCacheTtl) {
		serviceCache := masherytypes.MasheryServiceCache{CacheTtl: d.Get(MashSvcCacheTtl).(int)}
		if _, err := mashV3Cl.UpdateServiceCache(ctx, d.Id(), serviceCache); err != nil {
			updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "could not update cache ttl",
				Detail:   err.Error(),
			})
		}
	}

	if d.HasChange(MashSvcInteractiveDocsRoles) {
		roleUpsert := V3ServiceRolePermissionUpsertable(d)
		if err := mashV3Cl.SetServiceRoles(ctx, d.Id(), roleUpsert); err != nil {
			updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "could not emit roles",
				Detail:   err.Error(),
			})
		}
	}

	return updateDiagnostic
}

func serviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(v3client.Client)
	svcId := d.Id()

	// Verify that it's safe to delete this Mashery service, that it doesn't have
	// endpoints left to be associated with it, possibly unmanaged.
	endpointsLeft, err := c.CountEndpointsOf(ctx, svcId)
	if err != nil {
		return diag.FromErr(err)
	} else if endpointsLeft > 0 {
		return diag.Errorf("there are still %d endpoints left for this service", endpointsLeft)
	}

	err = c.DeleteService(ctx, svcId)
	if err != nil {
		return diag.FromErr(err)
	} else {
		// Clear the ID, the object has been deleted
		d.SetId("")
		return diag.Diagnostics{}
	}
}
