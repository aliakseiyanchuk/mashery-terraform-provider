package tfmapper

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DateMapper mapper that expects
type DateMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, masherytypes.MasheryJSONTime]
}

func (sfm *DateMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	if err := state.Set(sfm.Key, ""); err != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied null-value for field %s was not accepted: %s", sfm.Key, err.Error()),
			AttributePath: cty.GetAttrPath(sfm.Key),
		}
	} else {
		return nil
	}
}

func (sfm *DateMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	passingVal := ""
	if remoteVal := sfm.Locator(remote); remoteVal != nil {
		passingVal = remoteVal.ToString()
	}

	if err := state.Set(sfm.Key, passingVal); err != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied value for field %s was not accepted: %s", sfm.Key, err.Error()),
			AttributePath: cty.GetAttrPath(sfm.Key),
		}
	} else {
		return nil
	}
}

func (sfm *DateMapper[MType]) SchemaToRemote(_ *schema.ResourceData, _ *MType) {
	// No action
}
