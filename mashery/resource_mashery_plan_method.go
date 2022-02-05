package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryPlanMethod() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreatePlanMethod,
		DeleteContext: DeletePlanMethod,
		ReadContext:   ReadPlanMethod,
		UpdateContext: UpdatePlanMethod,
		Schema:        PlanMethodSchema,
	}
}

func CreatePlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	endpointIdent := MashPlanServiceEndpointUpsert(d)

	methIdent := ServiceEndpointMethodIdentifier{}
	methIdent.From(extractString(d, MashServiceEndpointMethodRef, ""))

	// Detect a collision
	if methIdent.EndpointId != endpointIdent.EndpointId {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Mismatching object",
			Detail:        "Method does not belong to the same endpoint",
			AttributePath: cty.GetAttrPath(MashServiceEndpointMethodRef),
		}}
	}

	filterIdentStr := extractString(d, MashServiceEndpointMethodFilterRef, "")
	var uFilter *ServiceEndpointMethodFilterIdentifier = nil

	if len(filterIdentStr) > 0 {
		filterIdent := ServiceEndpointMethodFilterIdentifier{}
		filterIdent.From(filterIdentStr)

		if methIdent.MethodId != filterIdent.MethodId {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Mismatching object",
				Detail:        "Filter does not belong to the containing method",
				AttributePath: cty.GetAttrPath(MashServiceEndpointMethodFilterRef),
			}}
		}

		uFilter = &filterIdent
	}

	v3cl := m.(v3client.Client)

	mashMeth := masherytypes.MasheryMethod{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id: methIdent.MethodId,
		},
	}

	if meth, err := v3cl.CreatePackagePlanMethod(ctx, endpointIdent, mashMeth); err != nil {
		return diag.FromErr(err)
	} else {
		data := map[string]interface{}{}

		rid := PlanEndpointMethodIdentifier{
			PlanEndpointIdentifier: PlanEndpointIdentifier{
				PlanServiceIdentifier: PlanServiceIdentifier{
					PlanIdentifier: PlanIdentifier{
						PackageId: endpointIdent.PackageId,
						PlanId:    endpointIdent.PlanId,
					},
					ServiceId: endpointIdent.ServiceId,
				},
				EndpointId: endpointIdent.EndpointId,
			},
			MethodId: meth.Id,
		}

		d.SetId(rid.Id())
		data[PlanEndpointMethodId] = meth.Id

		if uFilter != nil {
			v3Meth := masherytypes.MasheryPlanServiceEndpointMethod{
				MasheryPlanServiceEndpoint: masherytypes.MasheryPlanServiceEndpoint{
					MasheryPlanService: masherytypes.MasheryPlanService{
						PackageId: endpointIdent.PackageId,
						PlanId:    endpointIdent.PlanId,
						ServiceId: endpointIdent.ServiceId,
					},
					EndpointId: endpointIdent.EndpointId,
				},
				MethodId: meth.Id,
			}

			filterRef := masherytypes.MasheryServiceMethodFilter{
				FilterId: uFilter.FilterId,
			}

			if fltr, err := v3cl.CreatePackagePlanMethodFilter(ctx, v3Meth, filterRef); err != nil {
				return diag.FromErr(err)
			} else {
				data[PlanEndpointMethodFilterId] = fltr.Id
			}
		}

		return SetResourceFields(data, d)
	}
}

func UpdatePlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	if d.HasChange(MashServiceEndpointMethodFilterRef) {
		v3Id, dgs := MashPlanMethodIdentifier(d)
		if len(dgs) > 0 {
			return dgs
		}

		_, after := d.GetChange(MashServiceEndpointMethodFilterRef)

		afterStr := after.(string)

		data := map[string]interface{}{}

		if len(afterStr) == 0 {
			if err := v3cl.DeletePackagePlanMethodFilter(ctx, v3Id); err != nil {
				return diag.FromErr(err)
			}
			data[PlanEndpointMethodFilterId] = nil
		} else {
			filter, dgs := MashServiceMethodIdentifier(d)
			if len(dgs) > 0 {
				return dgs
			}

			if mf, err := v3cl.CreatePackagePlanMethodFilter(ctx, v3Id, filter); err != nil {
				return diag.FromErr(err)
			} else {
				data[PlanEndpointMethodFilterId] = mf.Id
			}
		}

		return SetResourceFields(data, d)
	} else {
		return diag.Diagnostics{}
	}
}

func ReadPlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Id, dgs := MashPlanMethodIdentifier(d)
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

		data := map[string]interface{}{}

		if filter, err := v3cl.GetPackagePlanMethodFilter(ctx, v3Id); err != nil {
			return diag.FromErr(err)
		} else {
			if filter == nil {
				data[PlanEndpointMethodFilterId] = nil
			} else {
				methId := ServiceEndpointMethodIdentifier{}
				methId.From(extractString(d, MashServiceEndpointMethodRef, ""))

				filterIdent := ServiceEndpointMethodFilterIdentifier{
					FilterId: filter.Id,
				}
				filterIdent.Inherit(methId)

				data[MashServiceEndpointMethodFilterRef] = filterIdent.Id()
				data[PlanEndpointMethodFilterId] = filter.Id
			}
		}

		return SetResourceFields(data, d)
	}
}

func DeletePlanMethod(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3Id, dgs := MashPlanMethodIdentifier(d)
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
