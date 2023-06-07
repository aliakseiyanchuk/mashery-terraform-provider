package tfmapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type Int64PointerFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, *int64]
}

func (sfm *Int64PointerFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	return SetKeyWithDiag(state, sfm.Key, 0)
}

func (sfm *Int64PointerFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	setVal := int64(0)

	if remoteVal := sfm.Locator(remote); remoteVal != nil {
		if ptr := *remoteVal; ptr != nil {
			setVal = *ptr
		}
	}

	return SetKeyWithDiag(state, sfm.Key, setVal)
}

func (sfm *Int64PointerFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	val := mashschema.ExtractInt64Pointer(state, sfm.Key, 0)
	*sfm.Locator(remote) = val
}
