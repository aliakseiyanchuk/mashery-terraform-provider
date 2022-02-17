package mashery

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

func resourceMasheryApplication() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		ReadContext:   applicationRead,
		CreateContext: applicationCreate,
		UpdateContext: applicationUpdate,
		DeleteContext: applicationDelete,
		// Schema
		Schema: mashschema.ApplicationMapper.TerraformSchema(),
		// Importer by ID
		//Importer: &mashschema.ResourceImporter{
		//	StateContext: mashschema.ImportStatePassthroughContext,
		//},
	}
}

func applicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	appIdent := mashschema.ApplicationMapper.GetIdentifier(d)

	doLogf("-> Trying to read application %s", appIdent.AppId)

	if rv, err := v3cl.GetApplication(ctx, appIdent.AppId); err != nil {
		doLogf("<- Failed to read application %s: %s", appIdent.AppId, err.Error())
		return diag.FromErr(err)
	} else {
		if rv != nil {
			doLogf("<- Application %s read", d.Id())
			mashschema.ApplicationMapper.SetState(ctx, rv, d)
		} else {
			doLogf("<- Application is not found anymore", d.Id())
			d.SetId("")
		}

		return nil
	}
}

func applicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	memberIdent := mashschema.ApplicationMapper.GetOwnerIdentifier(d)

	if len(memberIdent.MemberId) == 0 {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid parent object reference",
			Detail:        "Member reference is not valid",
			AttributePath: cty.GetAttrPath(mashschema.MashAppOwner),
		}}
	}

	v3cl := m.(v3client.Client)
	upsert, diags := mashschema.ApplicationMapper.UpsertableTyped(ctx, d)
	if len(diags) > 0 {
		return diags
	}

	if rv, err := v3cl.CreateApplication(ctx, memberIdent.MemberId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		appIdent := mashschema.ApplicationMapper.CreateIdentifierTyped()
		appIdent.MemberId = memberIdent.MemberId
		appIdent.Username = memberIdent.Username
		appIdent.AppId = rv.Id

		d.SetId(mashschema.CompoundId(appIdent))
		return diag.Diagnostics{}
	}
}

func applicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	upsert, rd := mashschema.ApplicationMapper.UpsertableTyped(ctx, d)
	if len(rd) > 0 {
		return rd
	}

	if rv, err := v3cl.UpdateApplication(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		return mashschema.ApplicationMapper.PersistTyped(ctx, rv, d)
	}
}

func applicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	appIdent := mashschema.ApplicationMapper.CreateIdentifierTyped()
	mashschema.CompoundIdFrom(appIdent, d.Id())

	if appKeys, err := v3cl.CountApplicationPackageKeys(ctx, appIdent.AppId); err != nil {
		return diag.FromErr(err)
	} else if appKeys > 0 {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Offending objects",
			Detail:   fmt.Sprintf("there are still %d package keys associated with this application", appKeys),
		}}
	}

	if err := v3cl.DeleteApplication(ctx, appIdent.AppId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
