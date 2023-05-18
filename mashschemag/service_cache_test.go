package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestServiceCacheBuilderWillProduceSchema(t *testing.T) {
	schema := ServiceCacheResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestServiceCacheIdentityMapping(t *testing.T) {
	autoTestIdentity(t, ServiceCacheResourceSchemaBuilder, masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	})

	autoTestParentIdentity(t, ServiceCacheResourceSchemaBuilder, masherytypes.ServiceIdentifier{
		ServiceId: "defg",
	})
}

func TestServiceCacheBuilderMappings(t *testing.T) {
	autoTestMappings(t, ServiceCacheResourceSchemaBuilder, func() masherytypes.ServiceCache {
		return masherytypes.ServiceCache{}
	})
}

// TestServiceCacheMappingDurationExpressionLeniency verifies that time expressions are accepted as long as
// they equate to the same number of minutes.
func TestServiceCacheMappingDurationExpressionLeniency(t *testing.T) {
	mapper, state := ServiceCacheResourceSchemaBuilder.MapperAndTestData()

	err := state.Set(mashschema.MashSvcCacheTtl, "60m")
	assert.Nil(t, err)

	inbound := masherytypes.ServiceCache{
		CacheTtl: 60,
	}

	dg := mapper.RemoteToSchema(&inbound, state)
	assert.Equal(t, 0, len(dg))

	err = state.Set(mashschema.MashSvcCacheTtl, "1h")
	assert.Nil(t, err)

	assert.Equal(t, "1h", state.Get(mashschema.MashSvcCacheTtl))
}
