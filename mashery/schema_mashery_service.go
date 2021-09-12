package mashery

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------
// Schemas for Mashery service resource and data sources.

// Unused args in schema -> to be moved to the provider resource
var supportedForwardedHeaders = []string{
	"access-token", "client-id", "scope", "user-context",
}

var supportedMasheryGrantTypes = []string{"authorization_code", "implicit", "password", "client_credentials"}
var supportedMasheryMacAlgorithms = []string{"hmac-sha-1", "hmac-sha-256"}

const (
	// MashSvcOAuthAccessTokenTtlEnabled Mashery Service OAuth Security Profile keys
	MashSvcOAuthAccessTokenTtlEnabled       = "access_token_ttl_enabled"
	MashSvcOAuthAccessTokenTtl              = "access_token_ttl"
	MashSvcOAuthAccessTokenType             = "access_token_type"
	MashSvcOAuthAllowMultipleToken          = "allow_multiple_token"
	MashSvcOAuthAuthorizationCodeTtl        = "authorization_code_ttl"
	MashSvcOAuthForwardedHeaders            = "forwarded_headers"
	MashSvcOAuthMasheryTokenApiEnabled      = "mashery_token_api_enabled"
	MashSvcOAuthRefreshTokenEnabled         = "refresh_token_enabled"
	MashSvcOAuthEnableRefreshTokenTtl       = "enable_refresh_token_ttl"
	MashSvcOAuthTokenBasedRateLimitsEnabled = "token_based_rate_limits_enabled"
	MashSvcOAuthForceOAuthRedirectUrl       = "force_oauth_redirect_url"
	MashSvcOAuthForceSSLRedirectUrlEnabled  = "force_ssl_redirect_url_enabled"
	MashSvcOAuthGrantTypes                  = "grant_types"
	MashSvcOAuthMacAlgorithm                = "mac_algorithm"
	MashSvcOAuthQpsLimitCeiling             = "qps_limit_ceiling"
	MashSvcOAuthRateLimitCeiling            = "rate_limit_ceiling"
	MashSvcOAuthRefreshTokenTtl             = "refresh_token_ttl"
	MashSvcOAuthSecureTokensEnabled         = "secure_tokens_enabled"

	// MashSvcOAuth Mashery OAuth Service Element
	MashSvcOAuth                = "oauth"
	MashSvcInteractiveDocsRoles = "iodocs_accessed_by"

	// MashSvcId Mashery Service element keys
	MashSvcId                   = "service_id"
	MashSvcMultiRef             = "service_ids"
	MashSvcExplained            = "service_explained"
	MashSvcName                 = "name"
	MashSvcNamePrefix           = "name_prefix"
	MashSvcCreated              = "created"
	MashSvcUpdated              = "updated"
	MashSvcEditorHandle         = "editor_handle"
	MashSvcRevisionNumber       = "revision_number"
	MashSvcRobotsPolicy         = "robots_policy"
	MashSvcCrossdomainPolicy    = "crossdomain_policy"
	MashSvcDescription          = "description"
	MashSvcQpsLimitOverall      = "qps_limit_overall"
	MashSvcServiceRFC3986Encode = "rfc3986_encode"
	MashSvcVersion              = "version"
)

func validateDuration(i interface{}, path cty.Path) diag.Diagnostics {
	if _, err := time.ParseDuration(i.(string)); err != nil {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "invalid duration",
			Detail:        fmt.Sprintf("expression %s is not a valid duration expression", i),
			AttributePath: path,
		}}
	} else {
		return diag.Diagnostics{}
	}
}

func validateZeroOrGreater(i interface{}, path cty.Path) diag.Diagnostics {
	if v, ok := i.(int); ok {
		if v < 0 {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Field must be zero or positive",
				Detail:        fmt.Sprintf("Value %d is negative", v),
				AttributePath: path,
			}}
		} else {
			return diag.Diagnostics{}
		}
	} else if v, ok := i.(int64); ok {
		if v < 0 {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Field must be zero or positive",
				Detail:        fmt.Sprintf("Value %d is negative", v),
				AttributePath: path,
			}}
		} else {
			return diag.Diagnostics{}
		}
	}

	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       "int or in64 required at this path",
		Detail:        fmt.Sprintf("unsupported type is %s", reflect.TypeOf(i).Name()),
		AttributePath: path,
	}}

}

// OAuthSecurityProfileSchema Mother schema for OAuth security profile.
var OAuthSecurityProfileSchema = map[string]*schema.Schema{
	MashSvcOAuthAccessTokenTtlEnabled: {
		Type:        schema.TypeBool,
		Description: "If enabled, the Access Token will expire after the specified time has passed",
		Optional:    true,
		Default:     true,
	},
	MashSvcOAuthAccessTokenTtl: {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Access token expires after the specified time has passed. TTL time is specified in seconds",
		Default:          "1h",
		ValidateDiagFunc: validateDuration,
	},
	MashSvcOAuthAccessTokenType: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Access token type, bearer of mac",
		Default:     "bearer",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			if str, ok := i.(string); ok {
				switch str {
				case "bearer":
				case "mac":
					return diag.Diagnostics{}
				default:
					return diag.Errorf("invalid value %s for access token type at path %s", str, path)
				}
			}
			return diag.Errorf("value %s is not a string at path %s", i, path)
		},
	},
	MashSvcOAuthAllowMultipleToken: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "If enabled, a unique access token will be issued for each access token request regardless of user context",
	},
	MashSvcOAuthAuthorizationCodeTtl: {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          "5m",
		Description:      "Authorization Code will expire after the specified time has passed. TTL time is specified in seconds.",
		ValidateDiagFunc: validateDuration,
	},
	MashSvcOAuthForwardedHeaders: {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "Mashery-generated headers that should be forwarded to the back-end",
	},
	MashSvcOAuthMasheryTokenApiEnabled: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If enabled, Access Token requests will be handled directly by Mashery via a dedicated Endpoint Request endpoint",
	},
	MashSvcOAuthRefreshTokenEnabled: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Allow developers to refresh tokens. The token can be refreshed when the partner is trusted but the risk lifespan is short",
	},
	MashSvcOAuthEnableRefreshTokenTtl: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "If enabled, the Refresh Token will expire after the specified TTL. TTL time is specified in seconds",
	},
	MashSvcOAuthTokenBasedRateLimitsEnabled: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Limit API calls per access token separately from API key rate limits",
	},
	MashSvcOAuthForceOAuthRedirectUrl: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "At the time of creating an access token, Mashery will validate that the client application provided a redirect URI field that matches with the callback URL specified during application registration",
	},
	MashSvcOAuthForceSSLRedirectUrlEnabled: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Allows Mashery to reject the request for authorization codes or access tokens that consist of a redirection URL other than HTTPS",
	},
	MashSvcOAuthGrantTypes: {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "Grant types selected for this service",
	},
	MashSvcOAuthMacAlgorithm: {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "MAC token algorithm",
		ValidateDiagFunc: func(inp interface{}, pth cty.Path) diag.Diagnostics {
			return validateStringValueInSet(inp, pth, &supportedMasheryMacAlgorithms)
		},
	},
	MashSvcOAuthQpsLimitCeiling: {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          0,
		ValidateDiagFunc: validateZeroOrGreater,
		Description:      "The throttle limit, i.e. calls per second, is applied to all access tokens granted for the API",
	},
	MashSvcOAuthRateLimitCeiling: {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          0,
		ValidateDiagFunc: validateZeroOrGreater,
		Description:      "The quota limit is applied to all access tokens granted for the API.",
	},
	MashSvcOAuthRefreshTokenTtl: {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          "768h",
		Description:      "The refresh token gets expired after it crosses the TTL value",
		ValidateDiagFunc: validateDuration,
	},
	MashSvcOAuthSecureTokensEnabled: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If enabled, Mashery stores tokens using a one-way SHA-256 hashed value",
	},
}

var OAuthSecurityProfileReadOnlySchema = cloneAsComputed(OAuthSecurityProfileSchema)

// ServiceSchema Mashery Service Definition schema.
var ServiceSchema = map[string]*schema.Schema{
	MashSvcId: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Service Id of this service",
	},
	MashSvcName: {
		Type:          schema.TypeString,
		Optional:      true,
		Computed:      true,
		ConflictsWith: []string{MashSvcNamePrefix},
		Description:   "Service name",
	},
	MashSvcNamePrefix: {
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{MashSvcName},
	},
	MashSvcCreated: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Created timestamp of this service",
	},
	MashSvcUpdated: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Timestamp of the latest service update",
	},
	MashSvcEditorHandle: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "User id which perform latest modification",
	},
	MashSvcRevisionNumber: {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Count of updates that were applied on this service after update",
	},
	MashSvcRobotsPolicy: {
		Type:     schema.TypeString,
		Computed: true,
	},
	MashSvcCrossdomainPolicy: {
		Type:     schema.TypeString,
		Computed: true,
	},
	MashSvcDescription: {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "Managed by Terraform",
		Description: "Description of this service",
	},
	MashSvcQpsLimitOverall: {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          0,
		ValidateDiagFunc: validateZeroOrGreater,
		Description:      "Maximum number of calls handled per second (QPS) across all developer keys for the API. Most customers do not set a value for this particular setting.",
	},
	MashSvcServiceRFC3986Encode: {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	MashSvcCacheTtl: {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	MashSvcOAuth: {
		Type:     schema.TypeSet,
		MinItems: 1,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: OAuthSecurityProfileSchema,
		},
	},
	MashSvcVersion: {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "0.0.1/TF",
		Description: "Deployed-defined version designator",
	},
	MashSvcInteractiveDocsRoles: {
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Description: "Set of role (or portal access groups) that can use IODocs of this service",
		// TODO: Figure out how diff suppression function
		// TODO: why can't this be assigned
		Elem: &schema.Resource{
			Schema: RolePermissionReferenceSchema,
		},
	},
}

// DataSourceMashSvcSchema Schema for service data source, which allows for the first service that matches
// the specified query
var DataSourceMashSvcSchema = DataSourceBaseSchema()

// -----------------------------------------
// V3 client -> Terraform schema conversion

// Formats the number of seconds as a duration for easier reading.
func toDurationFormat(seconds int64) string {
	if seconds <= 0 {
		return ""
	}

	d := time.Second * time.Duration(seconds)

	rv := strings.Builder{}

	// Reduce hours
	hours := d / time.Hour
	if hours > 0 {
		rv.Write([]byte(strconv.FormatInt(int64(hours), 10)))
		rv.Write([]byte("h"))

		d -= hours * time.Hour
	}

	// Reduce minutes
	minutes := d / time.Minute
	if minutes > 0 {
		rv.Write([]byte(strconv.FormatInt(int64(minutes), 10)))
		rv.Write([]byte("m"))

		d -= minutes * time.Minute
	}

	// Reduce seconds
	secsLeft := d / time.Second
	if secsLeft > 0 || rv.Len() == 0 {
		rv.Write([]byte(strconv.FormatInt(int64(secsLeft), 10)))
		rv.Write([]byte("s"))
	}

	return rv.String()
}

func durationToSeconds(dur string) int64 {
	if len(dur) == 0 {
		return 0
	}

	if theDur, err := time.ParseDuration(dur); err != nil {
		return 0
	} else {
		return int64(theDur.Seconds())
	}
}

func MapV3OAuthPolicyToTerraform(inp *v3client.MasheryOAuth) map[string]interface{} {
	return map[string]interface{}{
		MashSvcOAuthAccessTokenTtlEnabled:       inp.AccessTokenTtlEnabled,
		MashSvcOAuthAccessTokenTtl:              toDurationFormat(int64(inp.AccessTokenTtl)),
		MashSvcOAuthAccessTokenType:             inp.AccessTokenType,
		MashSvcOAuthAllowMultipleToken:          inp.AllowMultipleToken,
		MashSvcOAuthAuthorizationCodeTtl:        toDurationFormat(int64(inp.AuthorizationCodeTtl)),
		MashSvcOAuthForwardedHeaders:            inp.ForwardedHeaders,
		MashSvcOAuthMasheryTokenApiEnabled:      inp.MasheryTokenApiEnabled,
		MashSvcOAuthRefreshTokenEnabled:         inp.RefreshTokenEnabled,
		MashSvcOAuthEnableRefreshTokenTtl:       inp.EnableRefreshTokenTtl,
		MashSvcOAuthTokenBasedRateLimitsEnabled: inp.TokenBasedRateLimitsEnabled,
		MashSvcOAuthForceOAuthRedirectUrl:       inp.ForceOauthRedirectUrl,
		MashSvcOAuthForceSSLRedirectUrlEnabled:  inp.ForceSslRedirectUrlEnabled,
		MashSvcOAuthGrantTypes:                  inp.GrantTypes,
		MashSvcOAuthMacAlgorithm:                inp.MACAlgorithm,
		MashSvcOAuthQpsLimitCeiling:             inp.QPSLimitCeiling,
		MashSvcOAuthRateLimitCeiling:            inp.RateLimitCeiling,
		MashSvcOAuthRefreshTokenTtl:             toDurationFormat(inp.RefreshTokenTtl),
		MashSvcOAuthSecureTokensEnabled:         inp.SecureTokensEnabled,
	}
}

func MashSvcHasDirectUpsertableModifications(d *schema.ResourceData) bool {
	return d.HasChanges(MashSvcName, MashSvcDescription, MashSvcQpsLimitOverall, MashSvcServiceRFC3986Encode, MashSvcVersion)
}

func V3ServiceRolePermissionUpsertable(d *schema.ResourceData) []v3client.MasheryRolePermission {
	if setRaw, ok := d.GetOk(MashSvcInteractiveDocsRoles); ok {
		set, _ := setRaw.(*schema.Set)
		rv := make([]v3client.MasheryRolePermission, set.Len())

		for idx, vRaw := range set.List() {
			v, _ := vRaw.(map[string]interface{})
			rv[idx] = V3RolePermissionUpsertable(v)
		}

		return rv
	} else {
		return []v3client.MasheryRolePermission{}
	}
}

func V3ServiceOAuthProfileToTerraform(inp *v3client.MasheryOAuth, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashSvcOAuth: MapV3OAuthPolicyToTerraform(inp),
	}

	return SetResourceFields(data, d)
}

func V3ServiceToTerraform(inp *v3client.MasheryService, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashSvcId:                   inp.Id,
		MashSvcName:                 inp.Name,
		MashSvcCreated:              inp.Created.ToString(),
		MashSvcUpdated:              inp.Updated.ToString(),
		MashSvcEditorHandle:         inp.EditorHandle,
		MashSvcRevisionNumber:       inp.RevisionNumber,
		MashSvcRobotsPolicy:         inp.RobotsPolicy,
		MashSvcCrossdomainPolicy:    inp.CrossdomainPolicy,
		MashSvcDescription:          inp.Description,
		MashSvcQpsLimitOverall:      inp.QpsLimitOverall,
		MashSvcServiceRFC3986Encode: inp.RFC3986Encode,
		MashSvcVersion:              inp.Version,
	}

	if inp.SecurityProfile != nil {
		data[MashSvcOAuth] = []interface{}{
			MapV3OAuthPolicyToTerraform(&inp.SecurityProfile.OAuth),
		}
	} else {
		data[MashSvcOAuth] = nil
	}

	if inp.Cache != nil {
		data[MashSvcCacheTtl] = inp.Cache.CacheTtl
	} else {
		data[MashSvcCacheTtl] = 0
	}

	return SetResourceFields(data, d)
}

func V3ServiceRolesToTerraform(inp []v3client.MasheryRolePermission, d *schema.ResourceData) diag.Diagnostics {
	conv := V3RolesPermissionsToTerraform(inp)
	if err := d.Set(MashSvcInteractiveDocsRoles, conv); err != nil {
		return diag.FromErr(err)
	} else {
		return diag.Diagnostics{}
	}
}

// -----------------------------------------------------------------------
// Terraform -> V3 conversion routines

func V3SecurityProfileUpsertable(d *schema.ResourceData) v3client.MasherySecurityProfile {
	oauth := v3client.MasheryOAuth{}

	if inpRaw, ok := d.GetOk(MashSvcOAuth); ok {
		tfOauth := unwrapStructFromTerraformSet(inpRaw)

		oauth.AccessTokenTtlEnabled = tfOauth[MashSvcOAuthAccessTokenTtlEnabled].(bool)
		oauth.AccessTokenTtl = int(durationToSeconds(tfOauth[MashSvcOAuthAccessTokenTtl].(string)))
		oauth.AccessTokenType = tfOauth[MashSvcOAuthAccessTokenType].(string)
		oauth.AllowMultipleToken = tfOauth[MashSvcOAuthAllowMultipleToken].(bool)
		oauth.AuthorizationCodeTtl = int(durationToSeconds(tfOauth[MashSvcOAuthAuthorizationCodeTtl].(string)))
		oauth.ForwardedHeaders = convertSetToStringArray(tfOauth[MashSvcOAuthForwardedHeaders])
		oauth.MasheryTokenApiEnabled = tfOauth[MashSvcOAuthMasheryTokenApiEnabled].(bool)
		oauth.RefreshTokenEnabled = tfOauth[MashSvcOAuthRefreshTokenEnabled].(bool)
		oauth.EnableRefreshTokenTtl = tfOauth[MashSvcOAuthEnableRefreshTokenTtl].(bool)
		oauth.TokenBasedRateLimitsEnabled = tfOauth[MashSvcOAuthTokenBasedRateLimitsEnabled].(bool)
		oauth.ForceOauthRedirectUrl = tfOauth[MashSvcOAuthForceOAuthRedirectUrl].(bool)
		oauth.ForceSslRedirectUrlEnabled = tfOauth[MashSvcOAuthForceSSLRedirectUrlEnabled].(bool)
		oauth.GrantTypes = convertSetToStringArray(tfOauth[MashSvcOAuthGrantTypes])
		oauth.MACAlgorithm = tfOauth[MashSvcOAuthMacAlgorithm].(string)
		oauth.QPSLimitCeiling = int64(tfOauth[MashSvcOAuthQpsLimitCeiling].(int))
		oauth.RateLimitCeiling = int64(tfOauth[MashSvcOAuthRateLimitCeiling].(int))
		oauth.RefreshTokenTtl = durationToSeconds(tfOauth[MashSvcOAuthRefreshTokenTtl].(string))
		oauth.SecureTokensEnabled = tfOauth[MashSvcOAuthSecureTokensEnabled].(bool)
	}

	return v3client.MasherySecurityProfile{
		OAuth: oauth,
	}
}

func ClearServiceOAuthProfile(d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashSvcOAuth: nil,
	}
	return SetResourceFields(data, d)
}

func V3ServiceUpsertable(d *schema.ResourceData, includeCache, includeProfile bool) v3client.MasheryService {

	mashServ := v3client.MasheryService{
		AddressableV3Object: v3client.AddressableV3Object{
			Name: extractSetOrPrefixedString(d, MashSvcName, MashSvcNamePrefix),
		},
		Description:     extractString(d, MashSvcDescription, "Managed by Terraform"),
		QpsLimitOverall: extractInt64Pointer(d, MashSvcQpsLimitOverall, -10),
		RFC3986Encode:   extractBool(d, MashSvcServiceRFC3986Encode, true),
		Version:         extractString(d, MashSvcVersion, "0.0.1-TF"),
	}

	if includeCache {
		ttl := extractInt(d, MashSvcCacheTtl, 0)
		mashServ.Cache = &v3client.MasheryServiceCache{CacheTtl: ttl}
	}

	if includeProfile && d.Get(MashSvcOAuth) != nil {
		profile := V3SecurityProfileUpsertable(d)
		mashServ.SecurityProfile = &profile
	}

	return mashServ
}

// -----------------------------------
// Init section
// Perform data copy to various resources.

func inheritMasheryDataSourceSchema() {
	appendAsComputedInto(&ServiceSchema, &DataSourceMashSvcSchema)
	DataSourceMashSvcSchema[MashSvcMultiRef] = &schema.Schema{
		Type:        schema.TypeSet,
		Computed:    true,
		Description: "If multiple services matched, the ids of the matched services",
		Elem:        stringElem(),
	}

	DataSourceMashSvcSchema[MashSvcExplained] = &schema.Schema{
		Type:        schema.TypeMap,
		Computed:    true,
		Description: "Service ID to service name mapping",
		Elem:        stringElem(),
	}
}

func init() {
	inheritMasheryDataSourceSchema()
}
