package mashery_test

import (
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashery"
	"testing"
)

func TestWillCreateUniquePath(t *testing.T) {
	v := mashery.CreateUniquePath("path", 1610977328)

	assert.Equal(t, "/path_sXEbVb", v)
}
