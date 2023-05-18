package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceEndpointMethodFilterBuilderWillProduceSchema(t *testing.T) {
	schema := ServiceEndpointMethodFilterResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestServiceEndpointMethodFilterIdentityMapping(t *testing.T) {
	autoTestIdentity(t, ServiceEndpointMethodFilterResourceSchemaBuilder, masherytypes.ServiceEndpointMethodFilterIdentifier{
		FilterId: "filterId",
		ServiceEndpointMethodIdentifier: masherytypes.ServiceEndpointMethodIdentifier{
			MethodId: "mid",
			ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
				EndpointId: "defg",
				ServiceIdentifier: masherytypes.ServiceIdentifier{
					ServiceId: "abc",
				},
			},
		},
	})

	autoTestParentIdentity(t, ServiceEndpointMethodFilterResourceSchemaBuilder, masherytypes.ServiceEndpointMethodIdentifier{
		MethodId: "methId",
		ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
			EndpointId: "endpId",
			ServiceIdentifier: masherytypes.ServiceIdentifier{
				ServiceId: "defg",
			},
		},
	})
}

func TestServiceEndpointMethodFilterBuilderMappings(t *testing.T) {
	autoTestMappings(t, ServiceEndpointMethodFilterResourceSchemaBuilder, func() masherytypes.ServiceEndpointMethodFilter {
		return masherytypes.ServiceEndpointMethodFilter{}
	})
}
