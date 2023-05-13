package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type Int64PointerFieldMapper[MType any] struct {
	FieldMapperBase

	Locator LocatorFunc[MType, *int64]
}

func (sfm *Int64PointerFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
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

func (sfm *Int64PointerFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	if sfm.Schema.Computed && !sfm.Schema.Optional {
		return
	}

	val := mashschema.ExtractInt64Pointer(state, sfm.Key, 0)
	*sfm.Locator(remote) = val
}
