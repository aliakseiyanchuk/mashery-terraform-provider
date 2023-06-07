package tfmapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type StringMap map[string]string

type StringMapFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, *StringMap]
}

func (sfm *StringMapFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	emptyMap := map[string]string{}
	return SetKeyWithDiag(state, sfm.Key, emptyMap)
}

func (sfm *StringMapFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	if remoteVal := sfm.Locator(remote); remoteVal == nil {
		emptyMap := map[string]string{}
		return SetKeyWithDiag(state, sfm.Key, emptyMap)
	} else {
		return SetKeyWithDiag(state, sfm.Key, *remoteVal)
	}
}

func (sfm *StringMapFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	val := mashschema.ExtractStringMap(state, sfm.Key)
	*sfm.Locator(remote) = (*StringMap)(&val)
}
