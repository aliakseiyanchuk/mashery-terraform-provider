package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SerOrPrefixedFieldMapper mapper that expects
type SerOrPrefixedFieldMapper[MType any] struct {
	StringFieldMapper[MType]
	CompositeMapperBase

	PrefixKey string
}

func (sfm *SerOrPrefixedFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	for fld, _ := range sfm.CompositeSchema {
		if err := state.Set(fld, ""); err != nil {
			return &diag.Diagnostic{
				Severity:      diag.Error,
				Detail:        fmt.Sprintf("supplied null-value for field %s was not accepted: %s", sfm.Key, err.Error()),
				AttributePath: cty.GetAttrPath(sfm.Key),
			}
		}
	}

	return nil
}

func (sfm *SerOrPrefixedFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	value := ""

	if v, exists := state.GetOk(sfm.Key); exists {
		value = v.(string)
	} else {
		if prefix, exists := state.GetOk(sfm.PrefixKey); exists {
			value = resource.PrefixedUniqueId(prefix.(string))
		} else {
			value = resource.UniqueId()
		}

		_ = state.Set(sfm.Key, value)
	}

	dest := sfm.Locator(remote)
	*dest = value
}
