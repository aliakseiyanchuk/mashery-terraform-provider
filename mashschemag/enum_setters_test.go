package mashschemag

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/tfmapper"
	"testing"
)

// Automatically test that the mapper will accept all possible values and reject where these are malformed.
func autoTestFieldEnumValuesValidation[ParentIdent, Ident, MType any](t *testing.T, builder *tfmapper.SchemaBuilder[ParentIdent, Ident, MType], key string, values []string) {
	for _, val := range values {

		mapper, state := builder.MapperAndTestData()
		err := state.Set(key, []string{val})
		assert.Nil(t, err)

		dg := mapper.IsStateValid(state)
		assert.Equal(t, 0, len(dg))

		err = state.Set(key, []string{fmt.Sprintf("%s.invalid", val)})
		assert.Nil(t, err)
		dg = mapper.IsStateValid(state)
		assert.Equal(t, 1, len(dg))
	}
}
