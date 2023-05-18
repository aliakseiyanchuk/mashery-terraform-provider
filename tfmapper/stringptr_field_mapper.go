package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type StringPtrFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, *string]
}

func (sfm *StringPtrFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)

	val := ""

	if remoteVal != nil && *remoteVal != nil {
		val = **remoteVal
	}

	if err := state.Set(sfm.Key, val); err != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied value for field %s was not accepted: %s", sfm.Key, err.Error()),
			AttributePath: cty.GetAttrPath(sfm.Key),
		}
	} else {
		return nil
	}
}

func (sfm *StringPtrFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	impliedValue := ""
	if v, ok := sfm.Schema.Default.(string); ok {
		impliedValue = v
	}

	val := mashschema.ExtractString(state, sfm.Key, impliedValue)
	if len(val) > 0 {
		*sfm.Locator(remote) = &val
	}

}
