package tfmapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type StringFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, string]
}

func (sfm *StringFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	return SetKeyWithDiag(state, sfm.Key, "")
}

func (sfm *StringFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	setVal := ""
	if remoteVal := sfm.Locator(remote); remoteVal != nil {
		setVal = *remoteVal
	}

	return SetKeyWithDiag(state, sfm.Key, setVal)
}

func (sfm *StringFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	impliedValue := ""
	if v, ok := sfm.Schema.Default.(string); ok {
		impliedValue = v
	}

	val := mashschema.ExtractString(state, sfm.Key, impliedValue)
	*sfm.Locator(remote) = val
}
