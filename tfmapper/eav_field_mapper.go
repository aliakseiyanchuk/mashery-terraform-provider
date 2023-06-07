package tfmapper

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type EAVFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, *masherytypes.EAV]
}

func (sfm *EAVFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	emptyMap := map[string]interface{}{}

	return SetKeyWithDiag(state, sfm.Key, emptyMap)
}

func (sfm *EAVFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)

	if *remoteVal == nil {
		return SetKeyWithDiag(state, sfm.Key, map[string]string{})
	} else {
		passMap := map[string]string{}
		for k, v := range **remoteVal {
			passMap[k] = v
		}
		return SetKeyWithDiag(state, sfm.Key, passMap)
	}
}

func (sfm *EAVFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	val := mashschema.ExtractStringMap(state, sfm.Key)
	if len(val) == 0 {
		*sfm.Locator(remote) = nil
	} else {
		eav := masherytypes.EAV{}
		for k, v := range val {
			eav[k] = v
		}
		*sfm.Locator(remote) = &eav
	}
}
