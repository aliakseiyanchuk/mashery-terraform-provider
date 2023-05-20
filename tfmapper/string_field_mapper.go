package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type StringFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, string]
}

func (sfm *StringFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	if err := state.Set(sfm.Key, ""); err != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied null-value for field %s was not accepted: %s", sfm.Key, err.Error()),
			AttributePath: cty.GetAttrPath(sfm.Key),
		}
	} else {
		return nil
	}
}

func (sfm *StringFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
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

func (sfm *StringFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	impliedValue := ""
	if v, ok := sfm.Schema.Default.(string); ok {
		impliedValue = v
	}

	val := mashschema.ExtractString(state, sfm.Key, impliedValue)
	*sfm.Locator(remote) = val
}
