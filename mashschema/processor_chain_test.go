package mashschema_test

import (
	"testing"
)

func TestProcessChainOperation(t *testing.T) {
	//adapterA := masherytypes.Processor{
	//	PreProcessEnabled:  true,
	//	PostProcessEnabled: false,
	//	PostInputs:         nil,
	//	PreInputs:          []string{"a:1", "b:2", "c:3"},
	//	Adapter:            "adapter_a",
	//}
	//
	//adapterB := masherytypes.Processor{
	//	PreProcessEnabled:  false,
	//	PostProcessEnabled: true,
	//	PostInputs:         []string{"d:4", "e:5", "f:6"},
	//	PreInputs:          nil,
	//	Adapter:            "adapter_b",
	//}
	//
	//tfInput := []interface{}{
	//	mashschema.V3ProcessorConfigurationToTerraform(adapterA),
	//	mashschema.V3ProcessorConfigurationToTerraform(adapterB),
	//}
	//
	//res := TestResourceData(&mashschema.EndpointProcessorChainSchema)
	//err := res.Set(mashschema.MashEndpointProcessors, tfInput)
	//assert.Nil(t, err)
	//
	////diags := mashery.ProcessorChainCreateUpdate(context.TODO(), res, nil)
	////assert.Equal(t, 0, len(diags))
	//
	//assert.True(t, res.Get(mashschema.MashEndpointProcessorPreProcessEnabled).(bool))
	//assert.True(t, res.Get(mashschema.MashEndpointProcessorPostProcessEnabled).(bool))
	//
	////var emptyArr []string
	////cfg := ExtractStringArray(res, mashschema.MashEndpointProcessorPreConfig, &emptyArr)
	//
	////assert.True(t, FindInArray("adapter_a.a:1", &cfg) >= 0)
	////assert.True(t, FindInArray("adapter_a.b:2", &cfg) >= 0)
	////assert.True(t, FindInArray("adapter_a.c:3", &cfg) >= 0)
	////assert.True(t, FindInArray("processors:adapter_a", &cfg) >= 0)
	////
	////cfg = ExtractStringArray(res, mashschema.MashEndpointProcessorPostConfig, &emptyArr)
	////
	////assert.True(t, FindInArray("adapter_b.d:4", &cfg) >= 0)
	////assert.True(t, FindInArray("adapter_b.e:5", &cfg) >= 0)
	////assert.True(t, FindInArray("adapter_b.f:6", &cfg) >= 0)
	////assert.True(t, FindInArray("processors:adapter_b", &cfg) >= 0)
}
