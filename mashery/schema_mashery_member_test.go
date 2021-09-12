package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestWillGeneratePrefixedUsername(t *testing.T) {
	d := NewResourceData(&mashery.MemberSchema)
	prefix := "lspwd_prefix"
	assertOk(t, d.Set(mashery.MashMemberUserNamePrefix, prefix))

	upsert := mashery.MashMemberUpsertable(d)
	if !strings.HasPrefix(upsert.Username, prefix) {
		t.Errorf("Username (%s) was not prefixed correctly", upsert.Username)
	}
}

func TestV3MemberToResourceData(t *testing.T) {
	var tm v3client.MasheryJSONTime = v3client.MasheryJSONTime(time.Now())

	source := v3client.MasheryMember{
		AddressableV3Object: v3client.AddressableV3Object{
			Id:      "id",
			Created: &tm,
			Updated: &tm,
		},
		Username:    "username",
		Email:       "a@b.com",
		DisplayName: "dn",
		Uri:         "uri",
		Blog:        "blog",
		Im:          "im",
		Imsvc:       "imsvc",
		Phone:       "phone",
		Company:     "company",
		Address1:    "addr1",
		Address2:    "addr2",
		Locality:    "loc",
		Region:      "reg",
		PostalCode:  "postal",
		CountryCode: "cc",
		FirstName:   "first",
		LastName:    "last",
		AreaStatus:  "active",
		ExternalId:  "extId",
	}

	res := schema.Resource{
		Schema: mashery.MemberSchema,
	}

	d := res.TestResourceData()
	refId := "id::username"
	d.SetId(refId)

	diags := mashery.V3MemberToResourceData(&source, d)

	// These two keys are required and must not be set to prevent indefinite loop
	// with force-new.
	assertResourceDoesNotHaveKey(t, d, mashery.MashMemberEmail)
	assertResourceDoesNotHaveKey(t, d, mashery.MashMemberDisplayName)

	assertOk(t, d.Set(mashery.MashMemberEmail, "a@b.com"))
	assertOk(t, d.Set(mashery.MashMemberDisplayName, "dn"))

	if len(diags) > 0 {
		t.Errorf("full conversion has encountered %d errors where none were expected", len(diags))
	}

	reverse := mashery.MashMemberUpsertable(d)

	assertSameString(t, "id", &source.Id, &reverse.Id)

	assertSameString(t, "username", &source.Username, &reverse.Username)
	assertSameString(t, "email", &source.Email, &reverse.Email)
	assertSameString(t, "display name", &source.DisplayName, &reverse.DisplayName)
	assertSameString(t, "uri", &source.Uri, &reverse.Uri)
	assertSameString(t, "blog", &source.Blog, &reverse.Blog)
	assertSameString(t, "im", &source.Im, &reverse.Im)
	assertSameString(t, "imsvc", &source.Imsvc, &reverse.Imsvc)
	assertSameString(t, "phone", &source.Phone, &reverse.Phone)
	assertSameString(t, "company", &source.Company, &reverse.Company)
	assertSameString(t, "address1", &source.Address1, &reverse.Address1)
	assertSameString(t, "address2", &source.Address2, &reverse.Address2)
	assertSameString(t, "locality", &source.Locality, &reverse.Locality)
	assertSameString(t, "region", &source.Region, &reverse.Region)
	assertSameString(t, "postal code", &source.PostalCode, &reverse.PostalCode)
	assertSameString(t, "country code", &source.CountryCode, &reverse.CountryCode)
	assertSameString(t, "first name", &source.FirstName, &reverse.FirstName)
	assertSameString(t, "last name", &source.LastName, &reverse.LastName)
	assertSameString(t, "area status", &source.AreaStatus, &reverse.AreaStatus)
	assertSameString(t, "external id", &source.ExternalId, &reverse.ExternalId)
}
