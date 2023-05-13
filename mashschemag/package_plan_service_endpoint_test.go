package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestPackagePlanServiceEndpointBuilderWillProduceSchema(t *testing.T) {
	schema := PackagePlanServiceEndpointResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestPackagePlanServiceEndpointMapperIdentity(t *testing.T) {
	autoTestIdentity(t, PackagePlanServiceEndpointResourceSchemaBuilder, masherytypes.PackagePlanServiceEndpointIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "123",
			},
			PlanId: "456",
		},
		ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
			EndpointId: "abc",
			ServiceIdentifier: masherytypes.ServiceIdentifier{
				ServiceId: "defg",
			},
		},
	})
}

func TestPackagePlanServiceEndpointMapperServiceEndpointIdRef(t *testing.T) {
	mapper := PackagePlanServiceEndpointResourceSchemaBuilder.Mapper()
	testState := PackagePlanServiceEndpointResourceSchemaBuilder.TestResourceData()

	identIn := masherytypes.ServiceEndpointIdentifier{
		EndpointId: "abc",
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "def",
		},
	}
	err := mapper.TestAssign(mashschema.MashEndpointId, testState, identIn)
	assert.Nil(t, err)

	rv := PackagePlanServiceEndpointParam{}
	mapper.SchemaToRemote(testState, &rv)

	assert.True(t, reflect.DeepEqual(identIn, rv.ServiceEndpointIdentifier))
}
