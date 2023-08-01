package tfmapper

import (
	"fmt"
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

func SuppressSameDuration(key, _, _ string, s *schema.ResourceData) bool {
	ovs, nvs := s.GetChange(key)

	ovt, _ := time.ParseDuration(ovs.(string))
	nvt, _ := time.ParseDuration(nvs.(string))

	return ovt == nvt
}

func (sfm *DurationFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	return SetKeyWithDiag(state, sfm.Key, "")
}

func (sfm *DurationFieldMapper[MType]) RemoteToSchema(remote *MType, state *schema.ResourceData) *diag.Diagnostic {
	remoteVal := sfm.Locator(remote)
	thisVal := sfm.durationAsValue(state)

	if thisVal != *remoteVal {
		// Convert to value only if the numbers disagree.
		settingVal := fmt.Sprintf("%s", time.Duration(*remoteVal)*sfm.conversionUnit())
		return SetKeyWithDiag(state, sfm.Key, settingVal)
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
