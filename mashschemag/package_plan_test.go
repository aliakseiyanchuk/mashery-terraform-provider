package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPackagePlanBuilderWillProduceSchema(t *testing.T) {
	schema := PackagePlanResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestPackagePlanMapperIdentity(t *testing.T) {
	mapper := PackagePlanResourceSchemaBuilder.Mapper()
	testState := PackagePlanResourceSchemaBuilder.TestResourceData()

	testPack := masherytypes.PackagePlanIdentifier{
		PackageIdentifier: masherytypes.PackageIdentifier{
			PackageId: "123",
		},
		PlanId: "456",
	}

	err := mapper.AssignIdentity(testPack, testState)
	assert.Nil(t, err)

	readBack, idErr := mapper.Identity(testState)
	assert.Nil(t, idErr)

	assert.True(t, reflect.DeepEqual(testPack, readBack))
}

func TestPackagePlanMapperParentIdentity(t *testing.T) {
	mapper := PackagePlanResourceSchemaBuilder.Mapper()
	testState := PackagePlanResourceSchemaBuilder.TestResourceData()

	parentIdent := masherytypes.PackageIdentifier{
		PackageId: "abc",
	}

	err := mapper.TestSetPrentIdentity(parentIdent, testState)
	assert.Nil(t, err)

	readBackParentIdent, err := mapper.ParentIdentity(testState)
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(parentIdent, readBackParentIdent))
}

func TestPackagePlanMapperReadWriteBooleanFields(t *testing.T) {
	autoTestBoolMappings(t, PackagePlanResourceSchemaBuilder, func() masherytypes.Plan {
		return masherytypes.Plan{}
	})
}

func TestPackagePlanMapperReadWriteStringFields(t *testing.T) {
	autoTestStringMappings(t, PackagePlanResourceSchemaBuilder, func() masherytypes.Plan {
		return masherytypes.Plan{}
	}, "Status")
}

func TestPackagePlanMapperReadWriteIntFields(t *testing.T) {
	autoTestIntMappings(t, PackagePlanResourceSchemaBuilder, func() masherytypes.Plan {
		return masherytypes.Plan{}
	})
}

func TestPackagePlanMapperReadWriteInt64PtrFields(t *testing.T) {
	autoTestInt64PtrMappings(t, PackagePlanResourceSchemaBuilder, func() masherytypes.Plan {
		return masherytypes.Plan{}
	})
}

func TestPackagePlanMapperReadWriteEAVFields(t *testing.T) {
	autoTestEAVMappings(t, PackagePlanResourceSchemaBuilder, func() masherytypes.Plan {
		return masherytypes.Plan{}
	})
}
