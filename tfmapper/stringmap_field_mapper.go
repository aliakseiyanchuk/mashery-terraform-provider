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

	Locator LocatorFunc[MType, *StringMap]
}

func (sfm *StringMapFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)

	var setErr error

	if remoteVal == nil {
		emptyMap := map[string]string{}
		setErr = state.Set(sfm.Key, emptyMap)
	} else {
		setErr = state.Set(sfm.Key, *remoteVal)
	}

	// TOOO: repeating code that can be moved to the common method
	// deferred for the code optimization later on.
	if setErr != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied value for field %s was not accepted: %s", sfm.Key, setErr.Error()),
			AttributePath: cty.GetAttrPath(sfm.Key),
		}
	} else {
		return nil
	}
}

func (sfm *StringMapFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	val := mashschema.ExtractStringMap(state, sfm.Key)
	*sfm.Locator(remote) = (*StringMap)(&val)
}