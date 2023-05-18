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
