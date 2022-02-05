package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestApplicationIdParsing(t *testing.T) {
	ident := mashery.ApplicationIdentifier{}
	ident.From("fe733c9a-2c20-4efe-ad46-5070ebfe28bb::x002924_terraform20210126222356303100000003::a8d65611-2d9f-4647-8099-8861e47126c1")

	memberId := "fe733c9a-2c20-4efe-ad46-5070ebfe28bb"
	userName := "x002924_terraform20210126222356303100000003"
	appId := "a8d65611-2d9f-4647-8099-8861e47126c1"

	assertSameString(t, "MemberId", &ident.MemberId, &memberId)
	assertSameString(t, "Username", &ident.Username, &userName)
	assertSameString(t, "AppId", &ident.AppId, &appId)
}

func TestMashAppUpsertable(t *testing.T) {
	res := schema.Resource{
		Schema: mashery.AppSchema,
	}

	d := res.TestResourceData()
	_ = d.Set(mashery.MashAppOwner, "mid::m_username")

	now := masherytypes.MasheryJSONTime(time.Now())
	eav := masherytypes.EAV(map[string]string{
		"A": "B",
	})

	orig := masherytypes.MasheryApplication{
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
	chk := mashery.V3AppToResourceData(&orig, d)
	if len(chk) > 0 {
		t.Errorf("Setting fields encountered %d errors", len(chk))
	}
	assertResourceHasStringKey(t, d, mashery.MashAppOwnerUsername, "m_username")
	assertResourceHasStringKey(t, d, mashery.MashAppId, "appId")

	d.SetId("memberId::username::appId")

	// Reading back
	reverse := mashery.MashAppUpsertable(d)
	assertSameString(t, "id", &orig.Id, &reverse.Id)
	assertSameString(t, "name", &orig.Name, &reverse.Name)
	assertSameString(t, "username", &orig.Username, &reverse.Username)
	assertSameString(t, "description", &orig.Description, &reverse.Description)
	assertSameString(t, "type", &orig.Type, &reverse.Type)
	assertSameString(t, "type", &orig.Type, &reverse.Type)
	assertSameBool(t, "commercial", &orig.Commercial, &reverse.Commercial)
	assertSameBool(t, "ads", &orig.Ads, &reverse.Ads)
	assertSameString(t, "adsSystem", &orig.AdsSystem, &reverse.AdsSystem)
	assertSameString(t, "usageModel", &orig.UsageModel, &reverse.UsageModel)
	assertSameString(t, "tags", &orig.Tags, &reverse.Tags)
	assertSameString(t, "notes", &orig.Notes, &reverse.Notes)
	assertSameString(t, "howDidYouHear", &orig.HowDidYouHear, &reverse.HowDidYouHear)
	assertSameString(t, "preferredProtocol", &orig.PreferredProtocol, &reverse.PreferredProtocol)
	assertSameString(t, "preferredOutput", &orig.PreferredOutput, &reverse.PreferredOutput)
	assertSameString(t, "externalId", &orig.ExternalId, &reverse.ExternalId)
	assertSameString(t, "uri", &orig.Uri, &reverse.Uri)
	assertSameString(t, "OAuthRedirectUri", &orig.OAuthRedirectUri, &reverse.OAuthRedirectUri)
	assertSameEAV(t, "eav", orig.Eav, reverse.Eav)
}
