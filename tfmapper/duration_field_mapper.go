package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"time"
)

type DurationFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, int64]
	Unit    time.Duration
}

func (sfm *DurationFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)
	thisVal := sfm.durationAsValue(state)

	if thisVal != *remoteVal {
		// Convert to value only if the numbers disagree.
		duration := time.Duration(*remoteVal) * sfm.conversionUnit()
		setErr := state.Set(sfm.Key, fmt.Sprintf("%s", duration))

		if setErr != nil {
			return &diag.Diagnostic{
				Severity:      diag.Error,
				Detail:        fmt.Sprintf("supplied value for field %s was not accepted: %s", sfm.Key, setErr.Error()),
				AttributePath: cty.GetAttrPath(sfm.Key),
			}
		}
	}

	return nil
}

func (sfm *DurationFieldMapper[MType]) conversionUnit() time.Duration {
	unit := sfm.Unit
	if unit == 0 {
		unit = time.Second
	}

	return unit
}

func (sfm *DurationFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	val := sfm.durationAsValue(state)

	*sfm.Locator(remote) = val
}

func (sfm *DurationFieldMapper[MType]) durationAsValue(state *schema.ResourceData) int64 {
	dur := mashschema.ExtractString(state, sfm.Key, "")
	val := int64(0)

	if len(dur) > 0 {
		if theDur, err := time.ParseDuration(dur); err == nil {
			val = int64(theDur / sfm.conversionUnit())
		}
	}
	return val
}
