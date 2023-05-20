package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailTemplateSetBuilderWillProduceSchema(t *testing.T) {
	schema := EmailTemplateSetResourceSchemaBuilder.ResourceSchema()
	assert.True(t, len(schema) > 0)
}

func TestEmailTemplateSetBuilderMapperIdentity(t *testing.T) {
	autoTestIdentity(t, EmailTemplateSetResourceSchemaBuilder, "2355")
}

func TestEmailTemplateSetBuilderMapper(t *testing.T) {
	autoTestMappings(t, EmailTemplateSetResourceSchemaBuilder, func() masherytypes.EmailTemplateSet {
		return masherytypes.EmailTemplateSet{}
	})
}
