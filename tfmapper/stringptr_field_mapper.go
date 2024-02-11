package tfmapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type StringPtrFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, *string]
}

func (sfm *StringPtrFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	return SetKeyWithDiag(state, sfm.Key, "")
}

func (sfm *StringPtrFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	val := ""

	remoteVal := sfm.Locator(remote)
	if remoteVal != nil && *remoteVal != nil {
		if ptr := *remoteVal; ptr != nil {
			val = *ptr
		}
	}

	return SetKeyWithDiag(state, sfm.Key, val)
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
