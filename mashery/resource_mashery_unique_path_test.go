package mashery_test

import (
	"terraform-provider-mashery/mashery"
	"testing"
)

func TestWillCreateUniquePath(t *testing.T) {
	v := mashery.CreateUniquePath("path", 1610977328)

	if v != "/path_sXEbVb" {
		t.Errorf("Unexpected value %s", v)
	}
}
