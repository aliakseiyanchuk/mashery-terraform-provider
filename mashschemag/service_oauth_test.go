package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceOAuthBuilderWillProduceSchema(t *testing.T) {
	schema := ServiceOAuthResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestServiceOAuthIdentityMapping(t *testing.T) {
	autoTestIdentity(t, ServiceOAuthResourceSchemaBuilder, masherytypes.ServiceIdentifier{
		ServiceId: "abc",
	})

	autoTestParentIdentity(t, ServiceOAuthResourceSchemaBuilder, masherytypes.ServiceIdentifier{
		ServiceId: "defg",
	})
}

func TestServiceOAuthBuilderMappings(t *testing.T) {
	autoTestMappings(t, ServiceOAuthResourceSchemaBuilder, func() masherytypes.MasheryOAuth {
		return masherytypes.MasheryOAuth{}
	})
}
