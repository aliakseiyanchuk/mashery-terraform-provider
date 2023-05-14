package tfmapper

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
)

type EAVFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, *masherytypes.EAV]
}

func (sfm *EAVFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)

	var setErr error

	if *remoteVal == nil {
		emptyMap := map[string]string{}
		setErr = state.Set(sfm.Key, emptyMap)
	} else {
		passMap := map[string]string{}
		for k, v := range **remoteVal {
			passMap[k] = v
		}
		setErr = state.Set(sfm.Key, passMap)
	}

	// TOOO: repeating code that can be moved to the common method
	// deferred for the code optimization later on.
	if setErr != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied value for field %s was not accepted: %s", sfm.Key, setErr.Error()),
			AttributePath: cty.GetAttrPath(sfm.Key),
		}
	} else {
		return nil
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
