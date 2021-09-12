package mashery

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryPackageKey() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		ReadContext:   packageKeyRead,
		CreateContext: packageKeyCreate,
		UpdateContext: packageKeyUpdate,
		DeleteContext: packageKeyDelete,
		// Schema
		Schema: PackageKeySchema,
		// Importer by ID
		//Importer: &schema.ResourceImporter{
		//	StateContext: schema.ImportStatePassthroughContext,
		//},
	}
}

func packageKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	if rv, err := v3cl.GetPackageKey(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	} else {
		doLogJson(fmt.Sprintf("Recived the following JSON for package key %s", d.Id()), rv)

		if rv.Apikey != nil && rv.Secret != nil && len(*rv.Apikey) > 0 {
			V3PackageKeyToResourceData(rv, d)
			return diag.Diagnostics{}
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
	pid := PlanIdentifier{}
	pid.From(extractString(d, MashPlanId, ""))

	appIdent := ApplicationIdentifier{}
	appIdent.From(extractString(d, MashAppId, ""))

	if !pid.IsIdentified() || !appIdent.IsIdentified() {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Ambiguous parent",
			Detail:   "Package plan and applications are not recognized",
		}}
	}

	upsertable := MashPackageKeyUpsertable(d)

	v3cl := m.(v3client.Client)

	if rv, err := v3cl.CreatePackageKey(ctx, appIdent.AppId, upsertable); err != nil {
		return diag.FromErr(err)
	} else {
		V3PackageKeyToResourceData(rv, d)
		d.SetId(rv.Id)

		return diag.Diagnostics{}
	}
}

func packageKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	upsertable := MashPackageKeyUpsertable(d)
	v3cl := m.(v3client.Client)

	if rv, err := v3cl.UpdatePackageKey(ctx, upsertable); err != nil {
		return diag.FromErr(err)
	} else {
		V3PackageKeyToResourceData(rv, d)
		return diag.Diagnostics{}
	}
}

func packageKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	upsertable := MashPackageKeyUpsertable(d)
	if upsertable.Apikey == nil && upsertable.Secret == nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Attempt to delete effectively nil object",
			Detail:   fmt.Sprintf("Package key %s is effectively nil in state file. It is presumed as already deleted.", d.Id()),
		}}
	}

	if err := v3cl.DeletePackageKey(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
