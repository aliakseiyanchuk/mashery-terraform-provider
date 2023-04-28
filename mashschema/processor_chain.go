package mashschema

import (
	"context"
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

func (pcm *ProcessorChaiMapperImpl) adapterPrefixed(adapter string, cfg map[string]string) map[string]string {
	rv := map[string]string{}

	for key, value := range cfg {
		rv[adapter+"."+key] = value
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
	preCfg := map[string]string{}
	postCfg := map[string]string{}

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

			if len(v3Cfg.PreInputs) > 0 {
				v := pcm.adapterPrefixed(v3Cfg.Adapter, v3Cfg.PreInputs)
				mapMerge(&v, &preCfg)
			}
			if len(v3Cfg.PostInputs) > 0 {
				v := pcm.adapterPrefixed(v3Cfg.Adapter, v3Cfg.PostInputs)
				mapMerge(&v, &postCfg)
			}
		}
	}

	if len(preProcessors) > 0 {
		preCfg["processors"] = strings.Join(preProcessors, ",")
	}
	if len(postProcessors) > 0 {
		postCfg["processors"] = strings.Join(postProcessors, ",")
	}

	mergedCfg.PreInputs = preCfg
	mergedCfg.PostInputs = postCfg

	return mergedCfg
}

func mapMerge(source *map[string]string, into *map[string]string) {
	for k, v := range *source {
		(*into)[k] = v
	}
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

	computedOutput := cloneAsComputed(EndpointProcessorSchema)
	inheritAll(&computedOutput, &ProcessorChainMapper.schema)
}
