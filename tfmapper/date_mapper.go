package tfmapper

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DateMapper mapper that expects
type DateMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, masherytypes.MasheryJSONTime]
}

func (sfm *DateMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	return SetKeyWithDiag(state, sfm.Key, "")
}

func (sfm *DateMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	passingVal := ""
	if remoteVal := sfm.Locator(remote); remoteVal != nil {
		passingVal = remoteVal.ToString()
	}

	return SetKeyWithDiag(state, sfm.Key, passingVal)
}

func (sfm *DateMapper[MType]) SchemaToRemote(_ *schema.ResourceData, _ *MType) {
	// No action
}
