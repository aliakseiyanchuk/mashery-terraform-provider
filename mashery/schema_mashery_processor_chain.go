package mashery

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

const MashEndpointProcessors = "processor_chain"
const MashEndpointCompiledProcessorChain = "compiled_chain"

var EndpointProcessorChainSchema = map[string]*schema.Schema{
	MashEndpointProcessors: {
		Type:        schema.TypeList,
		MinItems:    2,
		Required:    true,
		Description: "List of processors to join in the processor chain",
		Elem: &schema.Resource{
			Schema: EndpointProcessorSchema,
		},
	},

	MashEndpointCompiledProcessorChain: {
		Type:        schema.TypeSet,
		Computed:    true,
		Description: "Computed processor chain",
		Elem: &schema.Resource{
			Schema: EndpointProcessorSchema,
		},
	},
}

func adapterPrefixed(adapter string, cfg []string) []string {
	rv := make([]string, len(cfg))
	for idx, v := range cfg {
		rv[idx] = adapter + "." + v
	}

	return rv
}

func ComputeChain(d *schema.ResourceData) v3client.Processor {
	mergedCfg := v3client.Processor{
		Adapter: "Mashery_Proxy_Processor_Chain",
	}

	var preProcessors []string
	var postProcessors []string
	var preCfg []string
	var postCfg []string

	if listRaw, ok := d.GetOk(MashEndpointProcessors); ok {
		list := listRaw.([]interface{})
		for _, tfProcCfg := range list {
			v3Cfg := V3ProcessorConfigurationFrom(tfProcCfg.(map[string]interface{}))

			if v3Cfg.PreProcessEnabled {
				preProcessors = append(preProcessors, v3Cfg.Adapter)
				mergedCfg.PreProcessEnabled = true
			}
			if v3Cfg.PostProcessEnabled {
				mergedCfg.PostProcessEnabled = true
				postProcessors = append(postProcessors, v3Cfg.Adapter)
			}

			if len(v3Cfg.PreInputs) > 0 {
				preCfg = append(preCfg, adapterPrefixed(v3Cfg.Adapter, v3Cfg.PreInputs)...)
			}
			if len(v3Cfg.PostInputs) > 0 {
				postCfg = append(postCfg, adapterPrefixed(v3Cfg.Adapter, v3Cfg.PostInputs)...)
			}
		}
	}

	if len(preProcessors) > 0 {
		preCfg = append(preCfg, fmt.Sprintf("processors:%s", strings.Join(preProcessors, ",")))
	}
	if len(postProcessors) > 0 {
		postCfg = append(postCfg, fmt.Sprintf("processors:%s", strings.Join(postProcessors, ",")))
	}

	mergedCfg.PreInputs = preCfg
	mergedCfg.PostInputs = postCfg

	return mergedCfg
}

func init() {
	computedOutput := cloneAsComputed(EndpointProcessorSchema)
	inheritAll(&computedOutput, &EndpointProcessorChainSchema)
}
