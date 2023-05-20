package mashschema

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

const MashEndpointProcessors = "processor_chain"
const MashEndpointCompiledProcessorChain = "compiled_chain"

var ProcessorChainMapper *ProcessorChaiMapperImpl

type ProcessorChaiMapperImpl struct {
	ResourceMapperImpl
}

func (pcm *ProcessorChaiMapperImpl) adapterPrefixed(adapter string, cfg []string) []string {
	rv := make([]string, len(cfg))
	for idx, v := range cfg {
		rv[idx] = adapter + "." + v
	}

	return rv
}

func (pcm *ProcessorChaiMapperImpl) PersistTyped(ctx context.Context, d *schema.ResourceData) {
	pcfg := pcm.ComputeChain(d)
	ServiceEndpointMapper.persistProcessor(&pcfg)
}

func (pcm *ProcessorChaiMapperImpl) ComputeChain(d *schema.ResourceData) masherytypes.Processor {
	mergedCfg := masherytypes.Processor{
		Adapter: "Mashery_Proxy_Processor_Chain",
	}

	var preProcessors []string
	var postProcessors []string
	var preCfg []string
	var postCfg []string

	if listRaw, ok := d.GetOk(MashEndpointProcessors); ok {
		list := listRaw.([]interface{})
		for _, tfProcCfg := range list {
			v3Cfg := ServiceEndpointMapper.processorUpsertableFromMap(tfProcCfg.(map[string]interface{}))

			if v3Cfg.PreProcessEnabled {
				preProcessors = append(preProcessors, v3Cfg.Adapter)
				mergedCfg.PreProcessEnabled = true
			}
			if v3Cfg.PostProcessEnabled {
				mergedCfg.PostProcessEnabled = true
				postProcessors = append(postProcessors, v3Cfg.Adapter)
			}

			//if len(v3Cfg.PreInputs) > 0 {
			//	preCfg = append(preCfg, pcm.adapterPrefixed(v3Cfg.Adapter, v3Cfg.PreInputs)...)
			//}
			//if len(v3Cfg.PostInputs) > 0 {
			//	postCfg = append(postCfg, pcm.adapterPrefixed(v3Cfg.Adapter, v3Cfg.PostInputs)...)
			//}
		}
	}

	if len(preProcessors) > 0 {
		preCfg = append(preCfg, fmt.Sprintf("processors:%s", strings.Join(preProcessors, ",")))
	}
	if len(postProcessors) > 0 {
		postCfg = append(postCfg, fmt.Sprintf("processors:%s", strings.Join(postProcessors, ",")))
	}

	//mergedCfg.PreInputs = preCfg
	//mergedCfg.PostInputs = postCfg

	return mergedCfg
}

func init() {
	ProcessorChainMapper = &ProcessorChaiMapperImpl{
		ResourceMapperImpl: ResourceMapperImpl{
			schema: map[string]*schema.Schema{
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
			},
		},
	}

	computedOutput := CloneAsComputed(EndpointProcessorSchema)
	inheritAll(&computedOutput, &ProcessorChainMapper.schema)
}
