package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryPackagePlanEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: createPackagePlanEndpoint,
		DeleteContext: deletePackagePlanEndpoint,
		ReadContext:   noopResourceOperation,
		// Update is not required
		//UpdateContext: noopResourceOperation,
		Schema: PlanEndpointSchema,
	}
}

func createPackagePlanEndpoint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	// We can include endpoint only from the same plan as current.
	// Otherwise it's collision. The developer should fix it.
	upsert, collisionDiag := V3MasheryPlanEndpointUpsertable(d)
	if len(collisionDiag) > 0 {
		return collisionDiag
	}

	doLogJson("WIll create package plan using the following upsert", &upsert)
	if _, err := v3cl.CreatePlanEndpoint(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		// If the call have succeeded, then the endpoint was added to the plan. There is no useful information
		// that could be extracted from the returned data, since it repeates data already defined in other objects.

		ident := PlanEndpointIdentifier{
			PlanServiceIdentifier: PlanServiceIdentifier{
				PlanIdentifier: PlanIdentifier{
					PackageId: upsert.PackageId,
					PlanId:    upsert.PlanId,
				},
				ServiceId: upsert.ServiceId,
			},
			EndpointId: upsert.EndpointId,
		}

		d.SetId(ident.Id())
		return diag.Diagnostics{}
	}
}

func deletePackagePlanEndpoint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	upsert, _ := V3MasheryPlanEndpointUpsertable(d)

	err := v3cl.DeletePlanEndpoint(ctx, upsert)
	if err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
