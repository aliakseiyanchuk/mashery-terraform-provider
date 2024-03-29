package tfmapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type BoolFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, bool]
}

func (sfm *BoolFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	return SetKeyWithDiag(state, sfm.Key, false)
}

func (sfm *BoolFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)
	return SetKeyWithDiag(state, sfm.Key, *remoteVal)
}

func (sfm *BoolFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	impliedValue := false
	if v, ok := sfm.Schema.Default.(bool); ok {
		impliedValue = v
	}

	val := mashschema.ExtractBool(state, sfm.Key, impliedValue)
	*sfm.Locator(remote) = val
}
