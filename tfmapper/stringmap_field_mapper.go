package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type StringMap map[string]string

type StringMapFieldMapper[MType any] struct {
	FieldMapperBase

	Locator LocatorFunc[MType, StringMap]
}

func (sfm *StringMapFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)

	// TOOO: repeating code that can be moved to the common method
	// deferred for the code optimization later on.
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

func (sfm *StringMapFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	// TODO: A candidate for the functional composition
	if sfm.Schema.Computed && !sfm.Schema.Optional {
		return
	}

	val := mashschema.ExtractStringMap(state, sfm.Key)
	*sfm.Locator(remote) = val
}
