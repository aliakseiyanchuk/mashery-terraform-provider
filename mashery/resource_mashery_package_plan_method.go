package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
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
	peIdent, svcMethIdent, filterIdent, dgs := mashschema.PackagePlanServiceEndpointMethodMapper.GetIdentifiers(d)

	if len(dgs) > 0 {
		return dgs
	}

	v3cl := m.(v3client.Client)

	mashMeth := masherytypes.MasheryMethod{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id: svcMethIdent.MethodId,
		},
	}

	if meth, err := v3cl.CreatePackagePlanMethod(ctx, *peIdent, mashMeth); err != nil {
		return diag.FromErr(err)
	} else {
		mashschema.PackagePlanServiceEndpointMethodMapper.SetIdentifier(peIdent, meth, d)

		if filterIdent != nil {
			v3Meth := masherytypes.MasheryPlanServiceEndpointMethod{
				MasheryPlanServiceEndpoint: masherytypes.MasheryPlanServiceEndpoint{
					MasheryPlanService: masherytypes.MasheryPlanService{
						PackageId: peIdent.PackageId,
						PlanId:    peIdent.PlanId,
						ServiceId: peIdent.ServiceId,
					},
					EndpointId: peIdent.EndpointId,
				},
				MethodId: meth.Id,
			}

			filterRef := masherytypes.MasheryServiceMethodFilter{
				FilterId: filterIdent.FilterId,
			}

			if _, err := v3cl.CreatePackagePlanMethodFilter(ctx, v3Meth, filterRef); err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	}
}

func UpdatePlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	if mashschema.PackagePlanServiceEndpointMethodMapper.HasFilterChange(d) {
		v3Id, dgs := mashschema.PackagePlanServiceEndpointMethodMapper.GetIdentifier(d)
		if len(dgs) > 0 {
			return dgs
		}

		_, after := mashschema.PackagePlanServiceEndpointMethodMapper.GetFilterChange(d)

		if !mashschema.IsIdentified(after) {
			if err := v3cl.DeletePackagePlanMethodFilter(ctx, v3Id); err != nil {
				return diag.FromErr(err)
			}
		} else {

			filter := masherytypes.MasheryServiceMethodFilter{
				MasheryServiceMethod: masherytypes.MasheryServiceMethod{
					MasheryServiceEndpoint: masherytypes.MasheryServiceEndpoint{
						ServiceId:  after.ServiceId,
						EndpointId: after.EndpointId,
					},
					MethodId: after.MethodId,
				},
				FilterId: after.FilterId,
			}

			if _, err := v3cl.CreatePackagePlanMethodFilter(ctx, v3Id, filter); err != nil {
				return diag.FromErr(err)
			}
		}

		return nil
	} else {
		return nil
	}
}

func ReadPlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Id, dgs := mashschema.PackagePlanServiceEndpointMethodMapper.GetIdentifier(d)
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
				mashschema.PackagePlanServiceEndpointMethodMapper.SetServiceFilterIdent(&v3Id, filter, d)
			}
		}

		return nil
	}
}

func DeletePlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Id, dgs := mashschema.PackagePlanServiceEndpointMethodMapper.GetIdentifier(d)
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
