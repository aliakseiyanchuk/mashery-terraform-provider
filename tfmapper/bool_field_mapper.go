package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type BoolFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, bool]
}

func (sfm *BoolFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)

	if err := state.Set(sfm.Key, *remoteVal); err != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied value for field %s was not accepted: %s", sfm.Key, err.Error()),
			AttributePath: cty.GetAttrPath(sfm.Key),
		}
	} else {
		return nil
	}
}

func (sfm *BoolFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	impliedValue := false
	if v, ok := sfm.Schema.Default.(bool); ok {
		impliedValue = v
	}

	val := mashschema.ExtractBool(state, sfm.Key, impliedValue)
	*sfm.Locator(remote) = val
}
