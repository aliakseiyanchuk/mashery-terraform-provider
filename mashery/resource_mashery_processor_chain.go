package mashery

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

func resourceMasheryProcessorChain() *schema.Resource {
	return &schema.Resource{
		ReadContext:   schema.NoopContext,
		DeleteContext: schema.NoopContext,
		UpdateContext: ProcessorChainCreateUpdate,
		CreateContext: ProcessorChainCreateUpdate,
		Schema:        mashschema.ProcessorChainMapper.TerraformSchema(),
	}
}

func ProcessorChainCreateUpdate(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	if len(d.Id()) == 0 {
		d.SetId(resource.PrefixedUniqueId("processorChain"))
		doLogf("Assigned ID to the resource")
	}

	mashschema.ProcessorChainMapper.PersistTyped(ctx, d)
	return nil
}
