package tfmapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type IntFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, int]
}

func (sfm *IntFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	return SetKeyWithDiag(state, sfm.Key, 0)
}

func (sfm *IntFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	setVal := 0
	if remoteVal := sfm.Locator(remote); remoteVal != nil {
		setVal = *remoteVal
	}

	return SetKeyWithDiag(state, sfm.Key, setVal)
}

func (sfm *IntFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	implied := 0
	if schemaDefault, ok := sfm.Schema.Default.(int); ok {
		implied = schemaDefault
	}

	val := mashschema.ExtractInt(state, sfm.Key, implied)
	*sfm.Locator(remote) = val
}
