package mashschema_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/stretchr/testify/assert"
	"strings"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestMemberMapperHasObjectName(t *testing.T) {
	assert.True(t, len(mashschema.MemberMapper.V3ObjectName()) > 0)
}

func TestMemberUpsertableFromEmptyConfiguration(t *testing.T) {
	d := mashschema.MemberMapper.TestResourceData()
	m, ctx, dg := mashschema.MemberMapper.UpsertableTyped(d)

	assert.Nil(t, ctx)
	assert.Equal(t, 0, len(dg))
	assert.Equal(t, "waiting", m.AreaStatus)
	assert.True(t, strings.HasPrefix(m.Username, "terraform-"))
}

func TestMemberUpsertableFromUserNamePrefix(t *testing.T) {
	cfg := map[string]interface{}{
		mashschema.MashMemberUserNamePrefix: "dtt",
	}

	d, rvd := mashschema.MemberMapper.TestResourceDataWith(cfg)
	assert.Equal(t, 0, len(rvd), "initial data set must be correct")

	m, ctx, dg := mashschema.MemberMapper.UpsertableTyped(d)

	assert.Nil(t, ctx)
	assert.Equal(t, 0, len(dg))
	assert.Equal(t, "waiting", m.AreaStatus)
	assert.True(t, strings.HasPrefix(m.Username, "dtt"))
}

func createCompleteMember() masherytypes.Member {
	tm := masherytypes.MasheryJSONTime(time.Now())

	source := masherytypes.Member{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:      "member-id",
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

	return source
}

func TestMemberSchemaMapping(t *testing.T) {
	mapper := mashschema.MemberMapper

	d := mapper.TestResourceData()

	source := createCompleteMember()
	dg := mapper.SetState(&source, d)

	assert.Equal(t, 0, len(dg))

	reverseIdentRaw, dg := mapper.V3Identity(d)
	assert.Equal(t, 0, len(dg))

	reverseIdent := reverseIdentRaw.(masherytypes.MemberIdentifier)
	assert.Equal(t, source.Id, reverseIdent.MemberId)
	assert.Equal(t, source.Username, reverseIdent.Username)

	reverseUpsert, _, dg := mapper.Upsertable(d)
	assert.Equal(t, 0, len(dg))

	source.Created = nil
	source.Updated = nil

	assert.True(t, assert.ObjectsAreEqualValues(source, reverseUpsert))
}
