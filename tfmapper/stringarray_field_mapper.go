package tfmapper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type StringArrayFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, []string]
}

func (sfm *StringArrayFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	emptyArray := []string{}
	return SetKeyWithDiag(state, sfm.Key, emptyArray)
}

func (sfm *StringArrayFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	var settingVal []string

	if remoteVal := sfm.Locator(remote); remoteVal != nil {
		// The change to the state will be accepted if the remote value contains multiple elements
		// or if it contains a single, non-empty string. Other situations are normalized as an
		// empty array.
		if len(*remoteVal) > 1 || (len(*remoteVal) == 1 && len((*remoteVal)[0]) > 0) {
			settingVal = *remoteVal
		}
	}

	return SetKeyWithDiag(state, sfm.Key, settingVal)
}

func (sfm *StringArrayFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	val := mashschema.ExtractStringArray(state, sfm.Key, &[]string{})
	*sfm.Locator(remote) = val
}
