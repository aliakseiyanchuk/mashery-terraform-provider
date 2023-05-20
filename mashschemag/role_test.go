package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoleBuilderWillProduceSchema(t *testing.T) {
	schema := RoleResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestRoleBuilderMapperIdentity(t *testing.T) {
	autoTestIdentity(t, RoleResourceSchemaBuilder, "2355")
}

func TestRoleBuilderMapper(t *testing.T) {
	autoTestMappings(t, RoleResourceSchemaBuilder, func() masherytypes.Role {
		return masherytypes.Role{}
	}, "Predefined")
}
