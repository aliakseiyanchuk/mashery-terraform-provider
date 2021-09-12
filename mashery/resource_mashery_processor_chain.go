package mashery

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMasheryProcessorChain() *schema.Resource {
	return &schema.Resource{
		ReadContext:   noopResourceOperation,
		DeleteContext: noopResourceOperation,
		UpdateContext: ProcessorChainCreateUpdate,
		CreateContext: ProcessorChainCreateUpdate,
		Schema:        EndpointProcessorChainSchema,
	}
}

func ProcessorChainCreateUpdate(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	proc := ComputeChain(d)
	doLogJson("Computed processor chain", proc)

	if len(d.Id()) == 0 {
		d.SetId(resource.PrefixedUniqueId("processorChain"))
		doLogf("Assigned ID to the resource")
	}

	data := V3ProcessorConfigurationToTerraform(proc)

	data[MashEndpointCompiledProcessorChain] = []interface{}{
		//V3ProcessorConfigurationToTerraform(proc),
	}

	doLogJson("Processor chain output", data)
	return SetResourceFields(data, d)
}
