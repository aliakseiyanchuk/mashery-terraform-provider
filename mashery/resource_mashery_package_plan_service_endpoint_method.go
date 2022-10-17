package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

func resourceMasheryPlanMethod() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreatePlanMethod,
		DeleteContext: DeletePlanMethod,
		ReadContext:   ReadPlanMethod,
		UpdateContext: UpdatePlanMethod,
		Schema:        mashschema.PackagePlanServiceEndpointMethodMapper.TerraformSchema(),
	}
}

func CreatePlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mapper := mashschema.PackagePlanServiceEndpointMethodMapper
	v3cl := m.(v3client.Client)

	methIdent, _, dg := mapper.UpsertableTyped(d)
	if len(dg) > 0 {
		return dg
	}

	if _, err := v3cl.CreatePackagePlanMethod(ctx, methIdent); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(mashschema.CompoundId(&methIdent))

		if mapper.HasFilter(d) {
			filterIdent, _ := mapper.GetFilterIdentity(d) // TODO: handle the diagnostics earlier

			if _, err := v3cl.CreatePackagePlanMethodFilter(ctx, filterIdent); err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	}
}

func UpdatePlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mapper := mashschema.PackagePlanServiceEndpointMethodMapper

	v3Id, dgs := mapper.V3Identity(d)
	if len(dgs) > 0 {
		return dgs
	}

	v3cl := m.(v3client.Client)

	if mapper.HasFilterChange(d) {
		if len(dgs) > 0 {
			return dgs
		}
		_, after := mapper.GetFilterChange(d)

		if !mashschema.IsIdentified(after) {
			if err := v3cl.DeletePackagePlanMethodFilter(ctx, v3Id); err != nil {
				return diag.FromErr(err)
			}
		} else {

			filterIdent, _ := mapper.GetFilterIdentity(d)

			if _, err := v3cl.CreatePackagePlanMethodFilter(ctx, filterIdent); err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	} else {
		return nil
	}
}

func ReadPlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mapper := mashschema.PackagePlanServiceEndpointMethodMapper
	v3Id, dgs := mapper.V3Identity(d)
	if len(dgs) > 0 {
		return dgs
	}

	v3cl := m.(v3client.Client)
	if srv, err := v3cl.GetPackagePlanMethod(ctx, v3Id); err != nil {
		return diag.FromErr(err)
	} else {
		// Delete gone resources
		if srv == nil {
			d.SetId("")
		}

		if filter, err := v3cl.GetPackagePlanMethodFilter(ctx, v3Id); err != nil {
			return diag.FromErr(err)
		} else {
			if filter == nil {
				mashschema.PackagePlanServiceEndpointMethodMapper.ClearFilter(d)
			} else {
				mashschema.PackagePlanServiceEndpointMethodMapper.SetServiceFilterIdent(filter, d)
			}
		}

		return nil
	}
}

func DeletePlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mapper := mashschema.PackagePlanServiceEndpointMethodMapper
	v3Id, dgs := mapper.V3Identity(d)
	if len(dgs) > 0 {
		return dgs
	}

	v3cl := m.(v3client.Client)
	err := v3cl.DeletePackagePlanMethod(ctx, v3Id)
	if err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
