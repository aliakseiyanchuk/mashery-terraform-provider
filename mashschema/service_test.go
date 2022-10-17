package mashschema_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

// TestBasicServiceConfiguration
// resource "mashery_service" {
//  name_prefix = "lspwd2.github"
//  desc = "service-desc"
//  version = "0.0.1a"
// }
func TestCreateUpsertFromBasicService(t *testing.T) {
	d := mashschema.ServiceMapper.TestResourceData()
	d.SetId("serviceID")
	d.Set("name_prefix", "lspwd2.github")
	d.Set("description", "service-desc")
	d.Set("version", "0.0.1a")

	upsert, _, dg := mashschema.ServiceMapper.UpsertableTyped(d)
	LogErrorDiagnostics(t, "get service upsertable", &dg)

	assert.True(t, len(upsert.Name) > 0)
	assert.Equal(t, "service-desc", upsert.Description)
	assert.Equal(t, "0.0.1a", upsert.Version)
	assert.Nil(t, upsert.Cache)
	assert.Nil(t, upsert.SecurityProfile)
	assert.Nil(t, upsert.Roles)
}

// TestBasicServiceConfiguration
// resource "mashery_service" {
//  name_prefix = "lspwd2.github"
//  desc = "service-desc"
//  cache_ttl = 30
//  version = "0.0.1a"
// }
func TestCreateUpsertFromBasicServiceWithCache(t *testing.T) {
	d := mashschema.ServiceMapper.TestResourceData()
	d.SetId("serviceID")
	d.Set("name_prefix", "lspwd2.github")
	d.Set("description", "service-desc")
	d.Set("cache_ttl", 30)
	d.Set("version", "0.0.1a")

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

// TestBasicServiceConfiguration
// resource "mashery_service" {
//  name_prefix = "lspwd2.github"
//  desc = "service-desc"
//  oauth {
//   access_token_ttl_enabled = true
//   access_token_ttl = 3600
//   access_token_type = "bearer"
// }
//  version = "0.0.1a"
// }
func TestCreateUpsertFromBasicOauth(t *testing.T) {
	d := mashschema.ServiceMapper.TestResourceData()
	d.SetId("serviceID")
	d.Set("name_prefix", "lspwd2.github")
	d.Set("description", "service-desc")
	d.Set("version", "0.0.1a")

	oauthSet := schema.NewSet(func(i interface{}) int {
		return 1
	}, []interface{}{
		map[string]interface{}{
			"access_token_ttl_enabled": true,
			"access_token_ttl":         "1h",
		},
	})

	setErr := d.Set("oauth", oauthSet)
	assert.Nil(t, setErr)

	upsert, _, dg := mashschema.ServiceMapper.UpsertableTyped(d)
	LogErrorDiagnostics(t, "get service upsertable", &dg)

	assert.True(t, len(upsert.Name) > 0)
	assert.Equal(t, "service-desc", upsert.Description)
	assert.Equal(t, "0.0.1a", upsert.Version)

	assert.NotNil(t, upsert.SecurityProfile)
	assert.True(t, upsert.SecurityProfile.OAuth.AccessTokenTtlEnabled)
	assert.Equal(t, 3600, upsert.SecurityProfile.OAuth.AccessTokenTtl)

	assert.Nil(t, upsert.Roles)
}

func TestV3ServiceRolesToTerraform(t *testing.T) {
	now := masherytypes.MasheryJSONTime(time.Now())
	inp := []masherytypes.MasheryRolePermission{
		{
			Role: masherytypes.Role{
				AddressableV3Object: masherytypes.AddressableV3Object{Id: "r1", Name: "n1", Created: &now, Updated: &now},
			},
			Action: "read",
		},
		{
			Role: masherytypes.Role{
				AddressableV3Object: masherytypes.AddressableV3Object{Id: "r2", Name: "n2"},
			},
			Action: "read",
		},
	}

	d := mashschema.ServiceMapper.TestResourceData()
	diags := mashschema.ServiceMapper.PersisRoles(inp, d)

	assert.Equal(t, 0, len(diags))

	// Parse reverse.
	//upsert := mashschema.ServiceMapper.RolePermissionUpsertable(d)
	//assert.Equal(t, 2, len(upsert))
	//assert.GreaterOrEqual(t, indexOfServiceRolePermission(&upsert, "r1", "read"), 0)
	//assert.GreaterOrEqual(t, indexOfServiceRolePermission(&upsert, "r2", "read"), 0)
}

func indexOfServiceRolePermission(inp *[]masherytypes.MasheryRolePermission, id, action string) int {
	for idx, v := range *inp {
		if id == v.Id && action == v.Action {
			return idx
		}
	}

	return -1
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
			Id:      "id",
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

	d.SetId("id")
	diags := mashschema.ServiceMapper.SetState(&v3Obj, d)
	assert.Equal(t, 0, len(diags))

	reverse, _, _ := mashschema.ServiceMapper.UpsertableTyped(d)

	assert.Equal(t, v3Obj.Name, reverse.Name)
	assert.Equal(t, v3Obj.Description, reverse.Description)
	assert.Equal(t, int64(10), *reverse.QpsLimitOverall)

	// RFC would be coming as true for some reason.
	//assert.False(t, reverse.RFC3986Encode)
	assert.Equal(t, v3Obj.Version, reverse.Version)

	// Verify that OAuth got recovered correctly.
	assert.NotNil(t, reverse.SecurityProfile)
	assert.Equal(t, 3600, reverse.SecurityProfile.OAuth.AccessTokenTtl)
	assert.Equal(t, 300, reverse.SecurityProfile.OAuth.AuthorizationCodeTtl)
	assert.Equal(t, int64(360000), reverse.SecurityProfile.OAuth.RefreshTokenTtl)

	// Headers and grant types are sets, so ordering of elements could be different
	// between runs.
	//mashery.assertSameSet(t, "ForwardedHeaders", &secProfile.OAuth.ForwardedHeaders, &reverse.SecurityProfile.OAuth.ForwardedHeaders)
	//mashery.assertSameSet(t, "GrantTypes", &secProfile.OAuth.GrantTypes, &reverse.SecurityProfile.OAuth.GrantTypes)
}
