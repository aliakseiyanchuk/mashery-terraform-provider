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
	"terraform-provider-mashery/mashschema"
)

func resourceMasheryService() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		ReadContext:   serviceRead,
		CreateContext: ServiceCreate,
		UpdateContext: serviceUpdate,
		DeleteContext: serviceDelete,
		// schema
		Schema: mashschema.ServiceSchema,
		// Importer by ID
		Importer: &schema.ResourceImporter{
			StateContext: importMasheryService,
		},
	}
}

func importMasheryService(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	mashV3Cl := m.(v3client.Client)
	svcId := masherytypes.ServiceIdentifier{ServiceId: d.Id()}

	if service, err := mashV3Cl.GetService(ctx, svcId); err != nil {
		return []*schema.ResourceData{}, errwrap.Wrapf("Failed to import this service: {{err}}", err)
	} else if service == nil {
		return []*schema.ResourceData{}, errors.New("no such service")
	} else {
		mashschema.ServiceMapper.PersistTyped(*service, d)

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
	svcId, dg := mashschema.ServiceMapper.V3Identity(d)
	if len(dg) > 0 {
		return dg
	}

	if rv, err := v3cl.GetServiceRoles(ctx, svcId.(masherytypes.ServiceIdentifier)); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error returned while retrieving roles",
			Detail:   err.Error(),
		}}
	} else {
		return mashschema.ServiceMapper.PersisRoles(rv, d)
	}
}

func ServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mashSvcUpsert, _, _ := mashschema.ServiceMapper.UpsertableTyped(d)

	doLogJson("Will attempt to create new service with this upsertable", mashSvcUpsert)

	mashV3Cl := m.(v3client.Client)
	if rv, err := mashV3Cl.CreateService(ctx, mashSvcUpsert); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(rv.Id)

		opDiagnostic := mashschema.ServiceMapper.PersistTyped(*rv, d)

		// After the service has been created, portal access groups need to be pushed. Otherwise, default
		// pre-populated list needs to be read.
		roleDiags := trySetServiceRoles(ctx, d, mashV3Cl)
		if len(roleDiags) > 0 {
			opDiagnostic = append(roleDiags)
		}

		return opDiagnostic
	}
}

func trySetServiceRoles(ctx context.Context, d *schema.ResourceData, mashV3Cl v3client.Client) diag.Diagnostics {
	svcId, dg := mashschema.ServiceMapper.V3Identity(d)
	if len(dg) > 0 {
		return dg
	}

	opDiagnostic := diag.Diagnostics{}

	if mashschema.ServiceMapper.IODocsRolesDefined(d) {
		roles := mashschema.ServiceMapper.UpsertableServiceRoles(d)
		doLogJson("Will attempt to set service roles with this upsertable", roles)

		err := mashV3Cl.SetServiceRoles(ctx, svcId.(masherytypes.ServiceIdentifier), *roles)
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
	svcId := masherytypes.ServiceIdentifier{ServiceId: d.Id()}

	if rv, err := mashV3Cl.GetService(ctx, svcId); err != nil {
		return diag.FromErr(err)
	} else {
		servDiags := mashschema.ServiceMapper.PersistTyped(*rv, d)
		rolesDiags := serviceReadRoles(ctx, d, mashV3Cl)

		return append(servDiags, rolesDiags...)
	}
}

func serviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mashV3Cl := m.(v3client.Client)

	ident, updateDiagnostic := mashschema.ServiceMapper.V3IdentityTyped(d)
	if len(updateDiagnostic) > 0 {
		return updateDiagnostic
	}

	if mashschema.ServiceMapper.HasDirectUpsertableChanges(d) {
		mashServ, _ := mashschema.ServiceMapper.DirectlyUpdateable(d)

		if rv, err := mashV3Cl.UpdateService(ctx, mashServ); err != nil {
			return diag.FromErr(err)
		} else {
			upd := mashschema.ServiceMapper.PersistTyped(*rv, d)
			if len(upd) > 0 {
				updateDiagnostic = append(updateDiagnostic, upd...)
			}
		}
	}

	if mashschema.ServiceMapper.OAuthProfileChanged(d) {
		curSet := mashschema.ServiceMapper.UpsertableSecurityProfile(d)
		if curSet == nil {
			if err := mashV3Cl.DeleteServiceOAuthSecurityProfile(ctx, ident); err != nil {
				updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "could not delete oauth security profile",
					Detail:   err.Error(),
				})
			} else {
				// OAuth profile has been deleted.
				_ = mashschema.ServiceMapper.ClearServiceOAuthProfile(d)
			}
		} else {
			requestedProfile := mashschema.ServiceMapper.UpsertableSecurityProfile(d)
			if actualOAuth, err := mashV3Cl.UpdateServiceOAuthSecurityProfile(ctx, *requestedProfile.OAuth); err != nil {
				updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "could not update oauth security profile",
					Detail:   err.Error(),
				})
			} else {
				oauthSave := mashschema.ServiceMapper.PersistOAuthProfile(actualOAuth, d)
				if len(oauthSave) > 0 {
					updateDiagnostic = append(updateDiagnostic, oauthSave...)
				}
			}
		}
	}

	if mashschema.ServiceMapper.CacheTTLChanged(d) {
		serviceCache := mashschema.ServiceMapper.CacheUpsertable(d)
		serviceIdent := masherytypes.ServiceIdentifier{
			ServiceId: d.Id(),
		}
		if serviceCache != nil {
			if _, err := mashV3Cl.UpdateServiceCache(ctx, serviceIdent, *serviceCache); err != nil {
				updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "could not update cache ttl",
					Detail:   err.Error(),
				})
			}
		} else {
			if err := mashV3Cl.DeleteServiceCache(ctx, serviceIdent); err != nil {
				updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "could not delete cache ttl",
					Detail:   err.Error(),
				})
			}
		}
	}

	if mashschema.ServiceMapper.IODocsRolesChanged(d) {
		if roleUpsert := mashschema.ServiceMapper.UpsertableServiceRoles(d); roleUpsert != nil {
			if err := mashV3Cl.SetServiceRoles(ctx, ident, *roleUpsert); err != nil {
				updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "could not emit roles",
					Detail:   err.Error(),
				})
			}
		} else {
			if err := mashV3Cl.DeleteServiceRoles(ctx, ident); err != nil {
				updateDiagnostic = append(updateDiagnostic, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "could not delete service roles",
					Detail:   err.Error(),
				})
			}

		}
	}

	return updateDiagnostic
}

func serviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(v3client.Client)
	svcId := masherytypes.ServiceIdentifier{ServiceId: d.Id()}

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
