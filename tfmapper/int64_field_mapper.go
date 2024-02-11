package tfmapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type Int64FieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, int64]
}

func (sfm *Int64FieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	return SetKeyWithDiag(state, sfm.Key, 0)
}

func (sfm *Int64FieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	setVal := int64(0)

	if remoteVal := sfm.Locator(remote); remoteVal != nil {
		setVal = *remoteVal
	}

	return SetKeyWithDiag(state, sfm.Key, setVal)
}

func (sfm *Int64FieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	implied := int64(0)
	if schemaDefault, ok := sfm.Schema.Default.(int64); ok {
		implied = schemaDefault
	}

	val := mashschema.ExtractInt64(state, sfm.Key, implied)
	*sfm.Locator(remote) = val
}
