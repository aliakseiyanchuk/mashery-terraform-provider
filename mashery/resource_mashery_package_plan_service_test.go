package mashery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackagePlanServiceResource(t *testing.T) {
	schema := PackagePlanServiceResource.TFDataSourceSchema()
	assert.NotNil(t, schema.Read)
	assert.Nil(t, schema.Update)
	assert.Nil(t, schema.UpdateContext)
}
