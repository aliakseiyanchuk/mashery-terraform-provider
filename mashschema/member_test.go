package mashschema_test

import (
	"testing"
)

func TestWillGeneratePrefixedUsername(t *testing.T) {
	//d := NewResourceData(&mashschema.MemberSchema)
	//prefix := "lspwd_prefix"
	//mashery.assertOk(t, d.Set(mashschema.MashMemberUserNamePrefix, prefix))
	//
	//upsert := mashschema.MashMemberUpsertable(d)
	//if !strings.HasPrefix(upsert.Username, prefix) {
	//	t.Errorf("Username (%s) was not prefixed correctly", upsert.Username)
	//}
}

func TestV3MemberToResourceData(t *testing.T) {
	//var tm masherytypes.MasheryJSONTime = masherytypes.MasheryJSONTime(time.Now())
	//
	//source := masherytypes.MasheryMember{
	//	AddressableV3Object: masherytypes.AddressableV3Object{
	//		Id:      "id",
	//		Created: &tm,
	//		Updated: &tm,
	//	},
	//	Username:    "username",
	//	Email:       "a@b.com",
	//	DisplayName: "dn",
	//	Uri:         "uri",
	//	Blog:        "blog",
	//	Im:          "im",
	//	Imsvc:       "imsvc",
	//	Phone:       "phone",
	//	Company:     "company",
	//	Address1:    "addr1",
	//	Address2:    "addr2",
	//	Locality:    "loc",
	//	Region:      "reg",
	//	PostalCode:  "postal",
	//	CountryCode: "cc",
	//	FirstName:   "first",
	//	LastName:    "last",
	//	AreaStatus:  "active",
	//	ExternalId:  "extId",
	//}
	//
	//res := schema.Resource{
	//	Schema: mashschema.MemberSchema,
	//}
	//
	//d := res.TestResourceData()
	//refId := "id::username"
	//d.SetId(refId)
	//
	//diags := mashschema.V3MemberToResourceData(&source, d)
	//
	//// These two keys are required and must not be set to prevent indefinite loop
	//// with force-new.
	//mashery.assertResourceDoesNotHaveKey(t, d, mashschema.MashMemberEmail)
	//mashery.assertResourceDoesNotHaveKey(t, d, mashschema.MashMemberDisplayName)
	//
	//mashery.assertOk(t, d.Set(mashschema.MashMemberEmail, "a@b.com"))
	//mashery.assertOk(t, d.Set(mashschema.MashMemberDisplayName, "dn"))
	//
	//if len(diags) > 0 {
	//	t.Errorf("full conversion has encountered %d errors where none were expected", len(diags))
	//}
	//
	//reverse := mashschema.MashMemberUpsertable(d)
	//
	//mashery.assertSameString(t, "id", &source.Id, &reverse.Id)
	//
	//mashery.assertSameString(t, "username", &source.Username, &reverse.Username)
	//mashery.assertSameString(t, "email", &source.Email, &reverse.Email)
	//mashery.assertSameString(t, "display name", &source.DisplayName, &reverse.DisplayName)
	//mashery.assertSameString(t, "uri", &source.Uri, &reverse.Uri)
	//mashery.assertSameString(t, "blog", &source.Blog, &reverse.Blog)
	//mashery.assertSameString(t, "im", &source.Im, &reverse.Im)
	//mashery.assertSameString(t, "imsvc", &source.Imsvc, &reverse.Imsvc)
	//mashery.assertSameString(t, "phone", &source.Phone, &reverse.Phone)
	//mashery.assertSameString(t, "company", &source.Company, &reverse.Company)
	//mashery.assertSameString(t, "address1", &source.Address1, &reverse.Address1)
	//mashery.assertSameString(t, "address2", &source.Address2, &reverse.Address2)
	//mashery.assertSameString(t, "locality", &source.Locality, &reverse.Locality)
	//mashery.assertSameString(t, "region", &source.Region, &reverse.Region)
	//mashery.assertSameString(t, "postal code", &source.PostalCode, &reverse.PostalCode)
	//mashery.assertSameString(t, "country code", &source.CountryCode, &reverse.CountryCode)
	//mashery.assertSameString(t, "first name", &source.FirstName, &reverse.FirstName)
	//mashery.assertSameString(t, "last name", &source.LastName, &reverse.LastName)
	//mashery.assertSameString(t, "area status", &source.AreaStatus, &reverse.AreaStatus)
	//mashery.assertSameString(t, "external id", &source.ExternalId, &reverse.ExternalId)
}
