package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackagePlanBuilderWillProduceSchema(t *testing.T) {
	schema := PackagePlanResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestPackagePlanMapperIdentity(t *testing.T) {
	autoTestIdentity(t, PackagePlanResourceSchemaBuilder, masherytypes.PackagePlanIdentifier{
		PackageIdentifier: masherytypes.PackageIdentifier{
			PackageId: "123",
		},
		PlanId: "456",
	})
}

func TestPackagePlanMapperParentIdentity(t *testing.T) {
	autoTestParentIdentity(t, PackagePlanResourceSchemaBuilder, masherytypes.PackageIdentifier{
		PackageId: "abc",
	})
}

func TestPackagePlanMapper(t *testing.T) {
	autoTestMappings(t, PackagePlanResourceSchemaBuilder, func() masherytypes.Plan {
		return masherytypes.Plan{}
	}, "Status")
}
