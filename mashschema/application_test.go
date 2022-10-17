package mashschema_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestMashAppUpsertable(t *testing.T) {
	d := mashschema.ApplicationMapper.TestResourceData()
	_ = d.Set(mashschema.MashAppOwner, "{\"mid\":\"a-b-c-d-e\",\"un\":\"m_username\"}")

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
		Username:          "m_username",
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
	chk := mashschema.ApplicationMapper.SetState(&orig, d)
	if len(chk) > 0 {
		t.Errorf("Setting fields encountered %d errors", len(chk))
	}
	//mashery.assertResourceHasStringKey(t, d, mashschema.MashAppOwnerUsername, "m_username")
	//mashery.assertResourceHasStringKey(t, d, mashschema.MashAppId, "appId")

	d.SetId("{\"mid\":\"mid\",\"un\":\"m_username\",\"appId\":\"appId\"}")

	// Reading back
	reverse, _, rvsDiags := mashschema.ApplicationMapper.UpsertableTyped(d)
	LogErrorDiagnostics(t, "app", &rvsDiags)

	assert.Equal(t, orig.Id, reverse.Id)
	assert.Equal(t, orig.Name, reverse.Name)
	assert.Equal(t, orig.Username, reverse.Username)
	assert.Equal(t, orig.Description, reverse.Description)
	assert.Equal(t, orig.Type, reverse.Type)
	assert.Equal(t, orig.Commercial, reverse.Commercial)
	assert.Equal(t, orig.Ads, reverse.Ads)
	assert.Equal(t, orig.AdsSystem, reverse.AdsSystem)
	assert.Equal(t, orig.UsageModel, reverse.UsageModel)
	assert.Equal(t, orig.Tags, reverse.Tags)
	assert.Equal(t, orig.Notes, reverse.Notes)
	assert.Equal(t, orig.HowDidYouHear, reverse.HowDidYouHear)
	assert.Equal(t, orig.PreferredProtocol, reverse.PreferredProtocol)
	assert.Equal(t, orig.PreferredOutput, reverse.PreferredOutput)
	assert.Equal(t, orig.Uri, reverse.Uri)
	assert.Equal(t, orig.OAuthRedirectUri, reverse.OAuthRedirectUri)
	assert.Equal(t, orig.Eav, reverse.Eav)
}
