package mashres

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourceTemplate_IsUpdateSuperfluous(t *testing.T) {
	assert.True(t, PackagePlanServiceResource.IsUpdateSuperfluous())
	assert.Nil(t, PackagePlanServiceResource.ResourceSchema().UpdateContext)

	assert.True(t, PackagePlanServiceEndpointResource.IsUpdateSuperfluous())
	assert.Nil(t, PackagePlanServiceEndpointResource.ResourceSchema().UpdateContext)

	assert.False(t, PackageResource.IsUpdateSuperfluous())
	assert.NotNil(t, PackageResource.ResourceSchema().UpdateContext)

	assert.False(t, ServiceResource.IsUpdateSuperfluous())
	assert.NotNil(t, ServiceResource.ResourceSchema().UpdateContext)
}
