package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"time"
)

type DurationFloat64FieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator LocatorFunc[MType, float64]
	Unit    time.Duration
}

func (sfm *DurationFloat64FieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
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

func (sfm *DurationFloat64FieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)
	thisVal := float64(sfm.durationAsValue(state))

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

func (sfm *DurationFloat64FieldMapper[MType]) conversionUnit() time.Duration {
	unit := sfm.Unit
	if unit == 0 {
		unit = time.Second
	}

	return unit
}

func (sfm *DurationFloat64FieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	val := sfm.durationAsValue(state)

	*sfm.Locator(remote) = float64(val)
}

func (sfm *DurationFloat64FieldMapper[MType]) durationAsValue(state *schema.ResourceData) int64 {
	dur := mashschema.ExtractString(state, sfm.Key, "")
	val := int64(0)

	if len(dur) > 0 {
		if theDur, err := time.ParseDuration(dur); err == nil {
			val = int64(theDur / sfm.conversionUnit())
		}
	}
	return val
}
