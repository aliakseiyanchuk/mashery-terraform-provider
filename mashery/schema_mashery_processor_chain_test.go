package mashery_test

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashery"
	"testing"
)

func TestProcessChainOperation(t *testing.T) {
	adapterA := v3client.Processor{
		PreProcessEnabled:  true,
		PostProcessEnabled: false,
		PostInputs:         nil,
		PreInputs:          []string{"a:1", "b:2", "c:3"},
		Adapter:            "adapter_a",
	}

	adapterB := v3client.Processor{
		PreProcessEnabled:  false,
		PostProcessEnabled: true,
		PostInputs:         []string{"d:4", "e:5", "f:6"},
		PreInputs:          nil,
		Adapter:            "adapter_b",
	}

	tfInput := []interface{}{
		mashery.V3ProcessorConfigurationToTerraform(adapterA),
		mashery.V3ProcessorConfigurationToTerraform(adapterB),
	}

	res := NewResourceData(&mashery.EndpointProcessorChainSchema)
	err := res.Set(mashery.MashEndpointProcessors, tfInput)
	assert.Nil(t, err)

	diags := mashery.ProcessorChainCreateUpdate(context.TODO(), res, nil)
	assert.Equal(t, 0, len(diags))

	assert.True(t, res.Get(mashery.MashEndpointProcessorPreProcessEnabled).(bool))
	assert.True(t, res.Get(mashery.MashEndpointProcessorPostProcessEnabled).(bool))

	var emptyArr []string
	cfg := mashery.ExtractStringArray(res, mashery.MashEndpointProcessorPreConfig, &emptyArr)

	assert.True(t, mashery.FindInArray("adapter_a.a:1", &cfg) >= 0)
	assert.True(t, mashery.FindInArray("adapter_a.b:2", &cfg) >= 0)
	assert.True(t, mashery.FindInArray("adapter_a.c:3", &cfg) >= 0)
	assert.True(t, mashery.FindInArray("processors:adapter_a", &cfg) >= 0)

	cfg = mashery.ExtractStringArray(res, mashery.MashEndpointProcessorPostConfig, &emptyArr)

	assert.True(t, mashery.FindInArray("adapter_b.d:4", &cfg) >= 0)
	assert.True(t, mashery.FindInArray("adapter_b.e:5", &cfg) >= 0)
	assert.True(t, mashery.FindInArray("adapter_b.f:6", &cfg) >= 0)
	assert.True(t, mashery.FindInArray("processors:adapter_b", &cfg) >= 0)
}
