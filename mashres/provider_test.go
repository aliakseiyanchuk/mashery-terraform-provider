package mashres

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviderWillLoad(t *testing.T) {
	assert.NotNil(t, Provider())
}
