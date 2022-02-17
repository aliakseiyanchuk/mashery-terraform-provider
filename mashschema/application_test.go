package mashschema_test

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestApplicationCompoundIdentityGenerator(t *testing.T) {
	idt := mashschema.ApplicationMapper.CreateIdentifierTyped()
	idt.MemberId = "a-b-c-d-e"
	idt.Username = "lspwd.github"
	idt.AppId = "a-b-c-d-e"

	str := mashschema.CompoundId(&idt)
	fmt.Println(str)
}

func TestApplicationIdParsing(t *testing.T) {
	ident := mashschema.ApplicationIdentifier{}
	mashschema.CompoundIdFrom(&ident, "{\"mid\":\"a-b-c-d-e\",\"un\":\"lspwd.github\",\"appId\":\"a-b-c-d-e\"}")

	memberId := "a-b-c-d-e"
	userName := "lspwd.github"
	appId := "a-b-c-d-e"

	assert.Equal(t, ident.MemberId, memberId, "MemberId")
	assert.Equal(t, ident.Username, userName, "Username")
	assert.Equal(t, ident.AppId, appId, "AppId")
}

func TestMashAppUpsertable(t *testing.T) {
	d := mashschema.ApplicationMapper.NewResourceData()
	_ = d.Set(mashschema.MashAppOwner, "{\"mid\":\"a-b-c-d-e\",\"un\":\"m_username\"}")

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
	chk := mashschema.ApplicationMapper.SetState(context.TODO(), &orig, d)
	if len(chk) > 0 {
		t.Errorf("Setting fields encountered %d errors", len(chk))
	}
	//mashery.assertResourceHasStringKey(t, d, mashschema.MashAppOwnerUsername, "m_username")
	//mashery.assertResourceHasStringKey(t, d, mashschema.MashAppId, "appId")

	d.SetId("{\"mid\":\"mid\",\"un\":\"m_username\",\"appId\":\"appId\"}")

	// Reading back
	reverse, rvsDiags := mashschema.ApplicationMapper.UpsertableTyped(context.TODO(), d)
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
