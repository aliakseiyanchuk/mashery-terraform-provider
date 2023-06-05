package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrganizationBuilderWillProduceSchema(t *testing.T) {
	schema := OrganizationResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestOrganizationBuilderMapperIdentity(t *testing.T) {
	autoTestIdentity(t, OrganizationResourceSchemaBuilder, "28840-dfgcc-34t32")
}

func TestOrganizationBuilderMapper(t *testing.T) {
	autoTestMappings(t, OrganizationResourceSchemaBuilder, func() masherytypes.Organization {
		return masherytypes.Organization{}
	})
}
