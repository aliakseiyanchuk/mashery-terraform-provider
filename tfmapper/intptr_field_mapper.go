package tfmapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type IntPointerFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, *int]
}

func (sfm *IntPointerFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	return SetKeyWithDiag(state, sfm.Key, 0)
}

func (sfm *IntPointerFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	setVal := 0

	if remoteVal := sfm.Locator(remote); remoteVal != nil {
		if ptr := *remoteVal; ptr != nil {
			setVal = *ptr
		}
	}

	return SetKeyWithDiag(state, sfm.Key, setVal)
}

func (sfm *IntPointerFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	val := mashschema.ExtractIntPointer(state, sfm.Key)
	*sfm.Locator(remote) = val
}
