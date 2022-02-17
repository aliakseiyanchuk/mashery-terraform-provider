package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mashschema "terraform-provider-mashery/mashschema"
)

func resourceMasheryPackagePlanEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: createPackagePlanEndpoint,
		DeleteContext: deletePackagePlanEndpoint,
		ReadContext:   schema.NoopContext,
		// Update is not required
		//UpdateContext: noopResourceOperation,
		Schema: mashschema.PlanServiceEndpointMapper.TerraformSchema(),
	}
}

func createPackagePlanEndpoint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)

	// We can include endpoint only from the same plan as current.
	// Otherwise it's collision. The developer should fix it.
	upsert, collisionDiag := mashschema.PlanServiceEndpointMapper.UpsertableTyped(d)
	if len(collisionDiag) > 0 {
		return collisionDiag
	}

	doLogJson("WIll create package plan using the following upsert", &upsert)
	if _, err := v3cl.CreatePlanEndpoint(ctx, upsert); err != nil {
		return diag.FromErr(err)
	} else {
		// If the call have succeeded, then the endpoint was added to the plan. There is no useful information
		// that could be extracted from the returned data, since it repeates data already defined in other objects.

		mashschema.PlanServiceEndpointMapper.SetIdentifier(d)
		return nil
	}
}

func deletePackagePlanEndpoint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	v3cl := m.(v3client.Client)
	upsert, _ := mashschema.PlanServiceEndpointMapper.UpsertableTyped(d)

	err := v3cl.DeletePlanEndpoint(ctx, upsert)
	if err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}
