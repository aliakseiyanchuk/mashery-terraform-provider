package mashery

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

func resourceMasheryMember() *schema.Resource {
	return &schema.Resource{
		// CRUD operations
		ReadContext:   memberRead,
		CreateContext: memberCreate,
		UpdateContext: memberUpdate,
		DeleteContext: memberDelete,
		// Schema
		Schema: mashschema.MemberMapper.TerraformSchema(),
		// Importer by ID
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func memberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	memberIdent := mashschema.MemberMapper.GetIdentifier(d)

	if rv, err := v3cl.GetMember(ctx, memberIdent.MemberId); err != nil {
		return diag.FromErr(err)
	} else {
		mashschema.MemberMapper.PersistTyped(ctx, rv, d)
		return diag.Diagnostics{}
	}
}

func memberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	upsert, _ := mashschema.MemberMapper.UpsertableTyped(d)

	if rv, err := v3cl.CreateMember(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		mashschema.MemberMapper.PersistTyped(ctx, rv, d)

		memberIdent := mashschema.MemberIdentifier{
			MemberId: rv.Id,
			Username: rv.Username,
		}

		d.SetId(mashschema.CompoundId(memberIdent))
		return diag.Diagnostics{}
	}
}

func memberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	upsert, _ := mashschema.MemberMapper.UpsertableTyped(d)

	if rv, err := v3cl.UpdateMember(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		mashschema.MemberMapper.PersistTyped(ctx, rv, d)
		return diag.Diagnostics{}
	}
}

func memberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	memberIdent := mashschema.MemberMapper.GetIdentifier(d)

	// Guard against unintended deletion.
	if apps, err := v3cl.CountApplicationsOfMember(ctx, memberIdent.MemberId); err != nil {
		return diag.FromErr(err)
	} else if apps > 0 {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Offending objects",
			Detail:   fmt.Sprintf("Member %s still contactin %d applications", d.Id(), apps),
		}}
	}

	// There's only member left.
	if err := v3cl.DeleteMember(ctx, memberIdent.MemberId); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
