package mashery

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

func resourceMasheryPackage() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		ReadContext:   packageRead,
		CreateContext: packageCreate,
		UpdateContext: packageUpdate,
		DeleteContext: packageDelete,
		// Schema
		Schema: mashschema.PackageSchema,
		// Importer by ID
		Importer: &schema.ResourceImporter{
			StateContext: importMasheryPackage,
		},
	}
}

func importMasheryPackage(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	mashV3Cl := m.(v3client.Client)

	if pack, err := mashV3Cl.GetPackage(ctx, d.Id()); err != nil {
		return []*schema.ResourceData{}, errwrap.Wrapf("Failed to import this package: {{err}}", err)
	} else if pack == nil {
		return []*schema.ResourceData{}, errors.New("No such service")
	} else {
		mashschema.PackageMapper.PersistTyped(ctx, pack, d)
		return []*schema.ResourceData{d}, nil
	}
}

func packageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	packId := d.Id()

	v3Client := m.(v3client.Client)
	if rv, err := v3Client.GetPackage(ctx, packId); err != nil {
		return diag.FromErr(err)
	} else {
		if rv != nil {
			return mashschema.PackageMapper.PersistTyped(ctx, rv, d)
		} else {
			d.SetId("")
		}
	}

	return nil
}

func packageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Client := m.(v3client.Client)
	upsert, _ := mashschema.PackageMapper.UpsertableTyped(d)

	doLogJson("Creating JSON using upsertable dataset", &upsert)

	if rv, err := v3Client.CreatePackage(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(rv.Id)

		// Warn that we've missed the Id on creation.
		if len(rv.Id) == 0 {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Missing id attribute",
				Detail:        "An attempt to create package succeeded, but resulting structure did not provide Id attribute",
				AttributePath: nil,
			}}
		}

		return mashschema.PackageMapper.PersistTyped(ctx, rv, d)
	}
}

func packageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	upsert, _ := mashschema.PackageMapper.UpsertableTyped(d)

	v3Client := m.(v3client.Client)
	if rv, err := v3Client.UpdatePackage(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		return mashschema.PackageMapper.PersistTyped(ctx, rv, d)
	}
}

func packageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	packId := d.Id()

	v3Client := m.(v3client.Client)

	// Count the plans still left in this package at the moment when the package is deleted.
	plansLeft, err := v3Client.CountPlans(ctx, packId)
	if err != nil {
		return diag.FromErr(err)
	} else if plansLeft > 0 {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Offending foreign objects",
			Detail:        fmt.Sprintf("Package %s still contains %d plans which are not managed in this state", packId, plansLeft),
			AttributePath: nil,
		}}
	}

	err = v3Client.DeletePackage(ctx, packId)
	if err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId("")
	}

	// Delete went okay.
	return nil
}
