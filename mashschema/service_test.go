package mashschema_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"strings"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

func TestServiceMapperSetup(t *testing.T) {
	assert.True(t, len(mashschema.ServiceMapper.V3ObjectName()) > 0)
}

func TestServiceFailsOnEmptyV3Mapper(t *testing.T) {
	mapper := mashschema.ServiceMapper
	d := mashschema.ServiceMapper.TestResourceData()
	_, dg := mapper.V3Identity(d)

	assert.Equal(t, 1, len(dg))

	assert.Equal(t, "lacking identification", dg[0].Summary)
	assert.Equal(t, "field(s) id must be set to identify V3 service object and must match object schema", dg[0].Detail)
}

// TestBasicServiceConfiguration
//
//	resource "mashery_service" {
//	 name_prefix = "lspwd2.github"
//	 desc = "service-desc"
//	 version = "0.0.1a"
//	}
func TestServiceCreateUpsertFromBasicConfig(t *testing.T) {
	d, dg := mashschema.ServiceMapper.TestResourceDataWith(case1MinimalConfiguration())
	assert.Equal(t, 0, len(dg))

	upsert, _, dg := mashschema.ServiceMapper.UpsertableTyped(d)

	assert.True(t, strings.HasPrefix(upsert.Name, "lspwd2.github"))
	assert.Equal(t, "service-desc", upsert.Description)
	assert.Equal(t, "0.0.1a", upsert.Version)

	assert.Nil(t, upsert.Cache)
	assert.Nil(t, upsert.SecurityProfile)
	assert.Nil(t, upsert.Roles)
}

func case1MinimalConfiguration() map[string]interface{} {
	data := map[string]interface{}{
		mashschema.MashSvcNamePrefix:  "lspwd2.github",
		mashschema.MashSvcDescription: "service-desc",
		mashschema.MashSvcVersion:     "0.0.1a",
	}
	return data
}

func case2CacheEnabled() map[string]interface{} {
	data := case1MinimalConfiguration()
	data[mashschema.MashSvcCacheTtl] = 30

	return data
}

func case3WithOAuth() map[string]interface{} {
	data := case1MinimalConfiguration()

	oauthSet := schema.NewSet(func(i interface{}) int {
		return 1
	}, []interface{}{
		map[string]interface{}{
			"grant_types": []string{"authorization_code"},
			"forwarded_headers": []string{"client-id",
				"scope", "user-context"},
			"access_token_ttl_enabled": true,
			"access_token_ttl":         "1h",
		},
	})
	data[mashschema.MashSvcOAuth] = oauthSet

	return data
}

func case4WithRoles(t *testing.T) map[string]interface{} {

	role := masherytypes.Role{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   "role-uuid",
			Name: "testRole",
		},
	}

	mapper := mashschema.RoleMapper
	d := mapper.TestResourceData()

	dg := mapper.PersistTyped(role, d)
	assert.Equal(t, 0, len(dg))

	rvData := case1MinimalConfiguration()

	roleRef := d.Get(mashschema.MashReadRolePermission).(map[string]interface{})

	rvData[mashschema.MashSvcInteractiveDocsRoles] = schema.NewSet(func(i interface{}) int {
		return 1
	}, []interface{}{
		roleRef,
	})

	return rvData
}

func TestServiceIdentityMapping(t *testing.T) {
	srv := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id: "svcId",
		},
	}

	mapper := mashschema.ServiceMapper
	d := mashschema.ServiceMapper.TestResourceData()

	dg := mapper.PersistTyped(srv, d)
	assert.Equal(t, 0, len(dg))

	ident, dg := mapper.V3IdentityTyped(d)
	assert.Equal(t, 0, len(dg))
	assert.Equal(t, "svcId", ident.ServiceId)
}

// TestServiceCreateUpsertFromBasicServiceWithCache
// Testing extraction of an upsert where the cache has been defined.
func TestServiceCreateUpsertFromBasicServiceWithCache(t *testing.T) {
	d, dg := mashschema.ServiceMapper.TestResourceDataWith(case2CacheEnabled())
	assert.Equal(t, 0, len(dg))

	upsert, _, dg := mashschema.ServiceMapper.UpsertableTyped(d)
	LogErrorDiagnostics(t, "get service upsertable", &dg)

	assert.True(t, len(upsert.Name) > 0)
	assert.Equal(t, "service-desc", upsert.Description)
	assert.Equal(t, "0.0.1a", upsert.Version)
	assert.NotNil(t, upsert.Cache)
	assert.Equal(t, 30, upsert.Cache.CacheTtl)

	assert.Nil(t, upsert.SecurityProfile)
	assert.Nil(t, upsert.Roles)
}

// TestServiceCreateUpsertFromBasicOauth
// Testing upsert extraction where OAuth has been defined
func TestServiceCreateUpsertFromBasicOauth(t *testing.T) {
	d, dg := mashschema.ServiceMapper.TestResourceDataWith(case3WithOAuth())
	assert.Equal(t, 0, len(dg))

	upsert, _, dg := mashschema.ServiceMapper.UpsertableTyped(d)

	assert.True(t, len(upsert.Name) > 0)
	assert.Equal(t, "service-desc", upsert.Description)
	assert.Equal(t, "0.0.1a", upsert.Version)

	assert.NotNil(t, upsert.SecurityProfile)
	assert.True(t, upsert.SecurityProfile.OAuth.AccessTokenTtlEnabled)
	assert.Equal(t, 3600, upsert.SecurityProfile.OAuth.AccessTokenTtl)

	assert.Nil(t, upsert.Cache)
	assert.Nil(t, upsert.Roles)
}

// TestServiceCreateUpsertWithRoles
// Testing service being upserted with defined roles
func TestServiceCreateUpsertWithRoles(t *testing.T) {
	d, dg := mashschema.ServiceMapper.TestResourceDataWith(case4WithRoles(t))
	assert.Equal(t, 0, len(dg))

	upsert, _, dg := mashschema.ServiceMapper.UpsertableTyped(d)

	assert.True(t, len(upsert.Name) > 0)
	assert.Equal(t, "service-desc", upsert.Description)
	assert.Equal(t, "0.0.1a", upsert.Version)

	assert.Nil(t, upsert.SecurityProfile)
	assert.Nil(t, upsert.Cache)
	assert.Nil(t, upsert.Roles)

	roles := mashschema.ServiceMapper.UpsertableServiceRoles(d)
	assert.NotNil(t, roles)

	assert.Equal(t, 1, len(*roles))
	assert.Equal(t, "role-uuid", (*roles)[0].Id)
	assert.Equal(t, "read", (*roles)[0].Action)
}

func TestServiceReceivedNilCacheWillDeletePreviousState(t *testing.T) {
	mapper := mashschema.ServiceMapper

	d, dg := mapper.TestResourceDataWith(case2CacheEnabled())
	assert.Equal(t, 0, len(dg))

	inboundService := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id: "svcId",
		},
	}

	dg = mapper.PersistTyped(inboundService, d)
	assert.Equal(t, 0, len(dg))

	upsert, _, _ := mapper.UpsertableTyped(d)
	assert.Nil(t, upsert.Cache)
}

func TestServiceReceivedNilOAuthWillDeletePreviousState(t *testing.T) {
	mapper := mashschema.ServiceMapper

	d, dg := mapper.TestResourceDataWith(case3WithOAuth())
	assert.Equal(t, 0, len(dg))

	inboundService := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id: "svcId",
		},
	}

	dg = mapper.PersistTyped(inboundService, d)
	assert.Equal(t, 0, len(dg))

	upsert := mapper.UpsertableSecurityProfile(d)
	assert.Nil(t, upsert)
}

func TestServiceReceivedNilRolesWillDeletePreviousState(t *testing.T) {
	mapper := mashschema.ServiceMapper

	d, dg := mapper.TestResourceDataWith(case4WithRoles(t))
	assert.Equal(t, 0, len(dg))

	inboundService := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id: "svcId",
		},
	}

	dg = mapper.PersistTyped(inboundService, d)
	assert.Equal(t, 0, len(dg))

	upsert := mapper.UpsertableServiceRoles(d)
	assert.Nil(t, upsert)
}

func TestV3ServiceToTerraform(t *testing.T) {
	d := mashschema.ServiceMapper.TestResourceData()

	now := masherytypes.MasheryJSONTime(time.Now())
	var limitOverall int64 = 10

	cache := masherytypes.ServiceCache{CacheTtl: 30}
	secProfile := masherytypes.MasherySecurityProfile{OAuth: &masherytypes.MasheryOAuth{
		AccessTokenTtlEnabled:       true,
		AccessTokenTtl:              3600,
		AccessTokenType:             "bearer",
		AllowMultipleToken:          true,
		AuthorizationCodeTtl:        300,
		ForwardedHeaders:            []string{"a", "b", "c"},
		MasheryTokenApiEnabled:      true,
		RefreshTokenEnabled:         true,
		EnableRefreshTokenTtl:       true,
		TokenBasedRateLimitsEnabled: true,
		ForceOauthRedirectUrl:       true,
		ForceSslRedirectUrlEnabled:  true,
		GrantTypes:                  []string{"t", "u", "v"},
		MACAlgorithm:                "hmac",
		QPSLimitCeiling:             -1,
		RateLimitCeiling:            -1,
		RefreshTokenTtl:             int64(360000),
		SecureTokensEnabled:         true,
	}}

	v3Obj := masherytypes.Service{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:      "svcId-num",
			Name:    "name",
			Created: &now,
			Updated: &now,
		},
		Cache:             &cache,
		Endpoints:         nil,
		EditorHandle:      "editor",
		RevisionNumber:    10,
		RobotsPolicy:      "robots",
		CrossdomainPolicy: "x-domand",
		Description:       "description",
		ErrorSets:         nil,
		QpsLimitOverall:   &limitOverall,
		RFC3986Encode:     false,
		SecurityProfile:   &secProfile,
		Version:           "version",
		Roles:             nil,
	}

	diags := mashschema.ServiceMapper.SetState(&v3Obj, d)
	assert.Equal(t, 0, len(diags))

	reverse, _, _ := mashschema.ServiceMapper.UpsertableTyped(d)
	assert.NotNil(t, reverse)

	// TODO: Figure out a way to compare what came back.
}
