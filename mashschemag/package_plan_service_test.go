package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestPackagePlanServiceBuilderWillProduceSchema(t *testing.T) {
	schema := PackagePlanServiceResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestPackagePlanServiceMapperIdentity(t *testing.T) {
	autoTestIdentity(t, PackagePlanServiceResourceSchemaBuilder, masherytypes.PackagePlanServiceIdentifier{
		PackagePlanIdentifier: masherytypes.PackagePlanIdentifier{
			PackageIdentifier: masherytypes.PackageIdentifier{
				PackageId: "123",
			},
			PlanId: "456",
		},
		ServiceIdentifier: masherytypes.ServiceIdentifier{
			ServiceId: "abcd",
		},
	})
}

func TestPackagePlanServiceMapperServiceIdRef(t *testing.T) {
	mapper := PackagePlanServiceResourceSchemaBuilder.Mapper()
	testState := PackagePlanServiceResourceSchemaBuilder.TestResourceData()

	err := mapper.TestAssign(mashschema.MashSvcRef, testState, masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	})
	assert.Nil(t, err)

	rv := PackagePlanServiceParam{}
	mapper.SchemaToRemote(testState, &rv)

	assert.Equal(t, "abc", rv.ServiceIdentifier.ServiceId)
}
