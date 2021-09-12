package mashery

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryPlanService() *schema.Resource {
	return &schema.Resource{
		CreateContext: createPlanService,
		ReadContext:   noopResourceOperation,
		DeleteContext: deletePlanService,
		Schema:        PlanService,
	}
}

func createPlanService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	upsert := V3MasheryPlanServiceUpsertable(d)

	doLogJson("Creating plan service using this upsert", upsert)

	if rv, err := v3cl.CreatePlanService(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		if rv.Id != "" {
			psi := PlanServiceIdentifier{
				PlanIdentifier: PlanIdentifier{
					PackageId: upsert.PackageId,
					PlanId:    upsert.PlanId,
				},
				ServiceId: upsert.ServiceId,
			}
			d.SetId(psi.Id())

			return diag.Diagnostics{}
		} else {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "missing identifier",
				Detail:        "value returned from V3 did not include minimum required identifier to confirm that object was successfully created",
				AttributePath: cty.GetAttrPath("id"),
			}}
		}
	}
}

func deletePlanService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	upsert := V3MasheryPlanServiceUpsertable(d)

	// Guard against deleting plan services that still have Mashery endpoints entered.
	if cnt, err := v3cl.CountPlanEndpoints(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else if cnt > 0 {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Offending objects found",
			Detail:        fmt.Sprintf("There are still %d endpoints attached to this service", cnt),
			AttributePath: nil,
		}}
	}

	if err := v3cl.DeletePlanService(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId("")
		return diag.Diagnostics{}
	}
}
