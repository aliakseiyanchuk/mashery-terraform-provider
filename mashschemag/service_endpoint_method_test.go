package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceEndpointMethodBuilderWillProduceSchema(t *testing.T) {
	schema := ServiceEndpointMethodResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestServiceEndpointMethodIdentityMapping(t *testing.T) {
	autoTestIdentity(t, ServiceEndpointMethodResourceSchemaBuilder, masherytypes.ServiceEndpointMethodIdentifier{
		MethodId: "mid",
		ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
			EndpointId: "defg",
			ServiceIdentifier: masherytypes.ServiceIdentifier{
				ServiceId: "abc",
			},
		},
	})

	autoTestParentIdentity(t, ServiceEndpointMethodResourceSchemaBuilder, masherytypes.ServiceEndpointIdentifier{
		EndpointId: "endpId",
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "defg",
		},
	})
}

func TestServiceEndpointMethodBuilderMappings(t *testing.T) {
	autoTestMappings(t, ServiceEndpointMethodResourceSchemaBuilder, func() masherytypes.ServiceEndpointMethod {
		return masherytypes.ServiceEndpointMethod{}
	})
}
