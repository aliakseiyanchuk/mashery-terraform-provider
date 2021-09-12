package mashery

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryApplication() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		ReadContext:   applicationRead,
		CreateContext: applicationCreate,
		UpdateContext: applicationUpdate,
		DeleteContext: applicationDelete,
		// Schema
		Schema: AppSchema,
		// Importer by ID
		//Importer: &schema.ResourceImporter{
		//	StateContext: schema.ImportStatePassthroughContext,
		//},
	}
}

func applicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	appIdent := ApplicationIdentifier{}
	appIdent.From(d.Id())

	doLogf("-> Trying to read application %s", appIdent.AppId)

	if rv, err := v3cl.GetApplication(ctx, appIdent.AppId); err != nil {
		doLogf("<- Failed to read application %s: %s", appIdent.AppId, err.Error())
		return diag.FromErr(err)
	} else {
		doLogf("<- Application %s read", d.Id())
		V3AppToResourceData(rv, d)
		return diag.Diagnostics{}
	}
}

func applicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	memberIdent := MemberIdentifier{}
	memberIdent.From(extractString(d, MashAppOwner, ""))

	if len(memberIdent.MemberId) == 0 {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid parent object reference",
			Detail:        "Member reference is not valid",
			AttributePath: cty.GetAttrPath(MashAppOwner),
		}}
	}

	v3cl := m.(v3client.Client)
	upsert := MashAppUpsertable(d)

	if rv, err := v3cl.CreateApplication(ctx, memberIdent.MemberId, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		V3AppToResourceData(rv, d)

		appIdent := ApplicationIdentifier{
			MemberId: memberIdent.MemberId,
			Username: rv.Username,
			AppId:    rv.Id,
		}

		d.SetId(appIdent.Id())
		return diag.Diagnostics{}
	}
}

func applicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	upsert := MashAppUpsertable(d)
	if rv, err := v3cl.UpdateApplication(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		V3AppToResourceData(rv, d)
		return diag.Diagnostics{}
	}
}

func applicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	appIdent := ApplicationIdentifier{}
	appIdent.From(d.Id())

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
