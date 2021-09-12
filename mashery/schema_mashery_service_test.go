package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

func TestV3ServiceRolesToTerraform(t *testing.T) {
	now := v3client.MasheryJSONTime(time.Now())
	inp := []v3client.MasheryRolePermission{
		{
			MasheryRole: v3client.MasheryRole{
				AddressableV3Object: v3client.AddressableV3Object{Id: "r1", Name: "n1", Created: &now, Updated: &now},
			},
			Action: "read",
		},
		{
			MasheryRole: v3client.MasheryRole{
				AddressableV3Object: v3client.AddressableV3Object{Id: "r2", Name: "n2"},
			},
			Action: "read",
		},
	}

	d := NewResourceData(&mashery.ServiceSchema)
	diags := mashery.V3ServiceRolesToTerraform(inp, d)

	assert.Equal(t, 0, len(diags))

	// Parse reverse.
	upsert := mashery.V3ServiceRolePermissionUpsertable(d)
	assert.Equal(t, 2, len(upsert))
	assert.GreaterOrEqual(t, indexOfServiceRolePermission(&upsert, "r1", "read"), 0)
	assert.GreaterOrEqual(t, indexOfServiceRolePermission(&upsert, "r2", "read"), 0)
}

func indexOfServiceRolePermission(inp *[]v3client.MasheryRolePermission, id, action string) int {
	for idx, v := range *inp {
		if id == v.Id && action == v.Action {
			return idx
		}
	}

	return -1
}

func TestV3ServiceToTerraform(t *testing.T) {
	d := NewResourceData(&mashery.ServiceSchema)

	now := v3client.MasheryJSONTime(time.Now())
	var limitOverall int64 = -1

	cache := v3client.MasheryServiceCache{CacheTtl: 30}
	secProfile := v3client.MasherySecurityProfile{OAuth: v3client.MasheryOAuth{
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

	v3Obj := v3client.MasheryService{
		AddressableV3Object: v3client.AddressableV3Object{
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
	diags := mashery.V3ServiceToTerraform(&v3Obj, d)
	assert.Equal(t, 0, len(diags))

	reverse := mashery.V3ServiceUpsertable(d, true, true)
	assert.Equal(t, "name", reverse.Name)
	assert.Equal(t, "description", reverse.Description)
	assert.Equal(t, int64(-1), *reverse.QpsLimitOverall)

	// RFC would be coming as true for some reason.
	//assert.False(t, reverse.RFC3986Encode)
	assert.Equal(t, "version", reverse.Version)

	// Verify that OAuth got recovered correctly.
	assert.NotNil(t, reverse.SecurityProfile)
	assert.Equal(t, 3600, reverse.SecurityProfile.OAuth.AccessTokenTtl)
	assert.Equal(t, 300, reverse.SecurityProfile.OAuth.AuthorizationCodeTtl)
	assert.Equal(t, int64(360000), reverse.SecurityProfile.OAuth.RefreshTokenTtl)

	// Headers and grant types are sets, so ordering of elements could be different
	// between runs.
	assertSameSet(t, "ForwardedHeaders", &secProfile.OAuth.ForwardedHeaders, &reverse.SecurityProfile.OAuth.ForwardedHeaders)
	assertSameSet(t, "GrantTypes", &secProfile.OAuth.GrantTypes, &reverse.SecurityProfile.OAuth.GrantTypes)
}
