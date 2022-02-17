package mashery

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mashschema "terraform-provider-mashery/mashschema"
)

func resourceMasheryPackageKey() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		ReadContext:   packageKeyRead,
		CreateContext: packageKeyCreate,
		UpdateContext: packageKeyUpdate,
		DeleteContext: packageKeyDelete,
		// Schema
		Schema: mashschema.PackageKeyMapper.TerraformSchema(),
		// Importer by ID
		//Importer: &mashschema.ResourceImporter{
		//	StateContext: mashschema.ImportStatePassthroughContext,
		//},
	}
}

func packageKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	keyIdent := mashschema.PackageKeyMapper.GetIdentifier(d)

	if rv, err := v3cl.GetPackageKey(ctx, keyIdent.KeyId); err != nil {
		return diag.FromErr(err)
	} else {
		doLogJson(fmt.Sprintf("Recived the following JSON for package key %s", d.Id()), rv)

		if rv.Apikey != nil && rv.Secret != nil && len(*rv.Apikey) > 0 {
			mashschema.PackageKeyMapper.PersistTyped(ctx, rv, d)
			return nil
		} else {
			d.SetId("")
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       "Effectively missing object returned",
				Detail:        "V3 API returned object which is effectively deleted.",
				AttributePath: nil,
			}}
		}
	}
}

func packageKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pid := mashschema.PackageKeyMapper.GetPlanIdentifier(d)
	appIdent := mashschema.PackageKeyMapper.GetApplicationIdentifier(d)

	if !mashschema.IsIdentified(pid) || !mashschema.IsIdentified(appIdent) {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Ambiguous parent",
			Detail:   "Package plan and applications are not recognized",
		}}
	}

	upsertable, _ := mashschema.PackageKeyMapper.UpsertableTyped(d)

	v3cl := m.(v3client.Client)

	if rv, err := v3cl.CreatePackageKey(ctx, appIdent.AppId, upsertable); err != nil {
		return diag.FromErr(err)
	} else {
		mashschema.PackageKeyMapper.SetIdentifier(appIdent, rv.Id, d)
		mashschema.PackageKeyMapper.PersistTyped(ctx, rv, d)

		return nil
	}
}

func packageKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	upsertable, _ := mashschema.PackageKeyMapper.UpsertableTyped(d)
	v3cl := m.(v3client.Client)

	if rv, err := v3cl.UpdatePackageKey(ctx, upsertable); err != nil {
		return diag.FromErr(err)
	} else {
		mashschema.PackageKeyMapper.PersistTyped(ctx, rv, d)
		return nil
	}
}

func packageKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	upsertable, _ := mashschema.PackageKeyMapper.UpsertableTyped(d)
	if upsertable.Apikey == nil && upsertable.Secret == nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Attempt to delete effectively nil object",
			Detail:   fmt.Sprintf("Package key %s is effectively nil in state file. It is presumed as already deleted.", d.Id()),
		}}
	}

	pkId := mashschema.PackageKeyMapper.GetIdentifier(d)
	if err := v3cl.DeletePackageKey(ctx, pkId.KeyId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
