package mashery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviderWillLoad(t *testing.T) {
	provider := Provider()
	assert.NotNil(t, provider)

}
