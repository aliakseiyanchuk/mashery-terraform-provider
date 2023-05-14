package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type IntFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, int]
}

func (sfm *IntFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)

	if remoteVal != nil {
		if err := state.Set(sfm.Key, *remoteVal); err != nil {
			return &diag.Diagnostic{
				Severity:      diag.Error,
				Detail:        fmt.Sprintf("supplied value for field %s was not accepted: %s", sfm.Key, err.Error()),
				AttributePath: cty.GetAttrPath(sfm.Key),
			}
		}
	} else {
		_ = state.Set(sfm.Key, 0)
	}

	return nil
}

func (sfm *IntFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	implied := 0
	if schemaDefault, ok := sfm.Schema.Default.(int); ok {
		implied = schemaDefault
	}

	val := mashschema.ExtractInt(state, sfm.Key, implied)
	*sfm.Locator(remote) = val
}
