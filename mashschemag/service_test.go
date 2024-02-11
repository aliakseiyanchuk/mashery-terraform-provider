package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceBuilderWillProduceSchema(t *testing.T) {
	schema := ServiceResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestServiceBuilderMappings(t *testing.T) {
	autoTestMappings(t, ServiceResourceSchemaBuilder, func() masherytypes.Service {
		return masherytypes.Service{}
	}, "EditorHandle", "RobotsPolicy", "CrossdomainPolicy", "RevisionNumber")
}

func TestServiceBuilderRoleMapping(t *testing.T) {
	perm := masherytypes.RolePermission{}
	perm.Id = "role1"
	perm.Action = "read"

	perm2 := masherytypes.RolePermission{}
	perm2.Id = "role2"
	perm2.Action = "read"

	roles := []masherytypes.RolePermission{perm, perm2}

	service := masherytypes.Service{
		Roles: &roles,
	}

	mapper, state := ServiceResourceSchemaBuilder.MapperAndTestData()

	dg := mapper.RemoteToSchema(&service, state)
	assert.Equal(t, 0, len(dg))

	readBack := masherytypes.Service{}
	mapper.SchemaToRemote(state, &readBack)

	assert.NotNil(t, readBack.Roles)
	assert.True(t, setEquals(*service.Roles, *readBack.Roles))
}
