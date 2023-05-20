package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type StringArrayFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, []string]
}

func (sfm *StringArrayFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	emptyArray := []string{}
	if err := state.Set(sfm.Key, emptyArray); err != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied null-value for field %s was not accepted: %s", sfm.Key, err.Error()),
			AttributePath: cty.GetAttrPath(sfm.Key),
		}
	} else {
		return nil
	}
}

func (sfm *StringArrayFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)

	var setErr error

	if remoteVal == nil {
		var emptySet []string
		setErr = state.Set(sfm.Key, emptySet)
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

func (sfm *StringArrayFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	val := mashschema.ExtractStringArray(state, sfm.Key, &[]string{})
	*sfm.Locator(remote) = val
}
