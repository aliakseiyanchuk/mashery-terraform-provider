package mashschema_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"strings"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestApplicationEmptyConfiguration(t *testing.T) {
	mapper := mashschema.ApplicationMapper
	d := mashschema.ApplicationMapper.TestResourceData()

	_, _, dg := mapper.Upsertable(d)
	assert.Equal(t, 1, len(dg))

}

func getPersistedIdOf(mapper mashschema.ResourceMapper, val interface{}) string {
	d := mapper.TestResourceData()

	mapper.SetState(val, d)
	return d.Id()
}

func TestApplicationMinimalConfiguration(t *testing.T) {
	memberObj := createCompleteMember()

	mapper := mashschema.ApplicationMapper
	data := map[string]interface{}{
		mashschema.MashAppOwner: getPersistedIdOf(mashschema.MemberMapper, &memberObj),
	}

	d, dg := mapper.TestResourceDataWith(data)
	assert.Equal(t, 0, len(dg))

	upsert, ctx, rvd := mapper.UpsertableTyped(d)
	assert.Equal(t, 0, len(rvd))

	assert.True(t, assert.ObjectsAreEqualValues(memberObj.Identifier(), ctx))
	assert.True(t, strings.HasPrefix(upsert.Name, "terraform-"))
}

func TestMashAppUpsertable(t *testing.T) {
	mapper := mashschema.ApplicationMapper

	memberObj := createCompleteMember()
	data := map[string]interface{}{
		mashschema.MashAppOwner: getPersistedIdOf(mashschema.MemberMapper, &memberObj),
	}

	d, dg := mapper.TestResourceDataWith(data)
	assert.Equal(t, 0, len(dg))

	now := masherytypes.MasheryJSONTime(time.Now())
	eav := masherytypes.EAV(map[string]string{
		"A": "B",
	})

	orig := masherytypes.Application{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:      "appId",
			Name:    "appName",
			Created: &now,
			Updated: &now,
		},
		Username:          memberObj.Username,
		Description:       "desc",
		Type:              "Type",
		Commercial:        true,
		Ads:               true,
		AdsSystem:         "AdSys",
		UsageModel:        "usage",
		Tags:              "tags",
		Notes:             "notes",
		HowDidYouHear:     "how",
		PreferredProtocol: "proto",
		PreferredOutput:   "output",
		ExternalId:        "extId",
		Uri:               "uri",
		OAuthRedirectUri:  "oauth",
		Eav:               &eav,
	}

	// Setting forward.
	chk := mapper.SetState(&orig, d)
	assert.Equal(t, 0, len(chk), "Setting fields encountered %d errors", len(chk))

	// Reading back
	reverse, userCtx, rvsDiags := mashschema.ApplicationMapper.UpsertableTyped(d)
	assert.Equal(t, 0, len(rvsDiags))
	assert.True(t, assert.ObjectsAreEqualValues(memberObj.Identifier(), userCtx))

	// Remove fields that are not set
	orig.Created = nil
	orig.Updated = nil

	assert.True(t, assert.ObjectsAreEqualValues(orig, reverse))
}
