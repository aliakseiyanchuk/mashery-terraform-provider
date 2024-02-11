package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestPackagePlanServiceEndpointMethodBuilderWillProduceSchema(t *testing.T) {
	schema := PackagePlanServiceEndpointMethodResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestPackagePlanServiceEndpointMethodBuilderParentIdentity(t *testing.T) {
	autoTestParentIdentity(t, PackagePlanServiceEndpointMethodResourceSchemaBuilder, masherytypes.PackagePlanServiceEndpointIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "abc",
			},
			PlanId: "bcd",
		},
		ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
			ServiceIdentifier: masherytypes.ServiceIdentifier{ServiceId: "cde"},
			EndpointId:        "def",
		},
	})
}

func TestPackagePlanServiceEndpointMethodMapperIdentity(t *testing.T) {
	autoTestIdentity(t, PackagePlanServiceEndpointMethodResourceSchemaBuilder, masherytypes.PackagePlanServiceEndpointMethodIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "123",
			},
			PlanId: "456",
		},
		ServiceEndpointMethodIdentifier: masherytypes.ServiceEndpointMethodIdentifier{
			MethodId: "methodId",
			ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
				EndpointId: "abc",
				ServiceIdentifier: masherytypes.ServiceIdentifier{
					ServiceId: "defg",
				},
			},
		},
	})
}

func TestPackagePlanServiceEndpointMethodMapperMethodIdent(t *testing.T) {
	mapper, testData := PackagePlanServiceEndpointMethodResourceSchemaBuilder.MapperAndTestData()

	filterIdent := masherytypes.ServiceEndpointMethodIdentifier{
		MethodId: "234",
		ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
			EndpointId: "345",
			ServiceIdentifier: masherytypes.ServiceIdentifier{
				ServiceId: "456",
			},
		},
	}

	err := mapper.TestAssign(mashschema.ServiceEndpointMethodRef, testData, filterIdent)
	assert.Nil(t, err)

	rv := PackagePlanServiceEndpointMethodParam{}
	mapper.SchemaToRemote(testData, &rv)

	assert.True(t, reflect.DeepEqual(filterIdent, rv.ServiceEndpointMethod))
}

func TestPackagePlanServiceEndpointMethodMapperFilterIdent(t *testing.T) {
	mapper, testData := PackagePlanServiceEndpointMethodResourceSchemaBuilder.MapperAndTestData()

	filterIdent := masherytypes.ServiceEndpointMethodFilterIdentifier{
		FilterId: "1234",
		ServiceEndpointMethodIdentifier: masherytypes.ServiceEndpointMethodIdentifier{
			MethodId: "234",
			ServiceEndpointIdentifier: masherytypes.ServiceEndpointIdentifier{
				EndpointId: "345",
				ServiceIdentifier: masherytypes.ServiceIdentifier{
					ServiceId: "456",
				},
			},
		},
	}

	err := mapper.TestAssign(mashschema.MashSvcEndpointMethodFilterRef, testData, filterIdent)
	assert.Nil(t, err)

	rv := PackagePlanServiceEndpointMethodParam{}
	mapper.SchemaToRemote(testData, &rv)

	assert.True(t, reflect.DeepEqual(filterIdent, rv.ServiceEndpointMethodFilterDesired))
}
