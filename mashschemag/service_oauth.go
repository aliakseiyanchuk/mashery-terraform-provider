package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
	"time"
)

var ServiceOAuthResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.MasheryOAuth]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.ServiceIdentifier]{
		IdentityFunc: func() masherytypes.ServiceIdentifier {
			return masherytypes.ServiceIdentifier{}
		},
	})

// Parent service identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.ServiceIdentifier]{
		Key: mashschema.MashSvcId,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Service Id, to which this OAuth security profile belongs",
		},
		IdentityFunc: func() masherytypes.ServiceIdentifier {
			return masherytypes.ServiceIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.ServiceIdentifier) bool {
			return len(inp.ServiceId) > 0
		},
	}

	ServiceOAuthResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

func init() {
	ServiceOAuthResourceSchemaBuilder.Add(&tfmapper.BoolFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *bool {
			return &in.AccessTokenTtlEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthAccessTokenTtlEnabled,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If enabled, the Access Token will expire after the specified time has passed",
				Optional:    true,
				Default:     true,
			},
		},
	}).Add(&tfmapper.DurationFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *int64 {
			return &in.AccessTokenTtl
		},
		Unit: time.Second,
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthAccessTokenTtl,
			Schema: &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Access token expires after the specified time has passed. TTL time is specified in seconds",
				Default:          "1h",
				ValidateDiagFunc: mashschema.ValidateDuration,
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *string {
			return &in.AccessTokenType
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthAccessTokenType,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access token type, bearer of mac",
				Default:     "bearer",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(i, path, &[]string{"bearer", "mac"})
				},
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *bool {
			return &in.AllowMultipleToken
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthAllowMultipleToken,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If enabled, a unique access token will be issued for each access token request regardless of user context",
			},
		},
	}).Add(&tfmapper.DurationFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *int64 {
			return &in.AuthorizationCodeTtl
		},
		Unit: time.Second,
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthAuthorizationCodeTtl,
			Schema: &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "5m",
				Description:      "Authorization Code will expire after the specified time has passed. TTL time is specified in seconds.",
				ValidateDiagFunc: mashschema.ValidateDuration,
			},
		},
	}).Add(&tfmapper.StringArrayFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *[]string {
			return &in.ForwardedHeaders
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthForwardedHeaders,
			Schema: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         mashschema.StringHashcode,
				Description: "Mashery-generated headers that should be forwarded to the back-end",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *bool {
			return &in.MasheryTokenApiEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthMasheryTokenApiEnabled,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, Access Token requests will be handled directly by Mashery via a dedicated Endpoint Request endpoint",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *bool {
			return &in.RefreshTokenEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthRefreshTokenEnabled,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow developers to refresh tokens. The token can be refreshed when the partner is trusted but the risk lifespan is short",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *bool {
			return &in.EnableRefreshTokenTtl
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthEnableRefreshTokenTtl,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, the Refresh Token will expire after the specified TTL. TTL time is specified in seconds",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *bool {
			return &in.TokenBasedRateLimitsEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthTokenBasedRateLimitsEnabled,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Limit API calls per access token separately from API key rate limits",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *bool {
			return &in.ForceOauthRedirectUrl
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthForceOAuthRedirectUrl,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "At the time of creating an access token, Mashery will validate that the client application provided a redirect URI field that matches with the callback URL specified during application registration",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *bool {
			return &in.ForceSslRedirectUrlEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthForceSSLRedirectUrlEnabled,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allows Mashery to reject the request for authorization codes or access tokens that consist of a redirection URL other than HTTPS",
			},
		},
	}).Add(&tfmapper.StringArrayFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *[]string {
			return &in.GrantTypes
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthGrantTypes,
			Schema: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         mashschema.StringHashcode,
				Description: "Grant types selected for this service",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *string {
			return &in.MACAlgorithm
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthMacAlgorithm,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "MAC token algorithm",
				ValidateDiagFunc: func(inp interface{}, pth cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(inp, pth, &mashschema.SupportedMasheryMacAlgorithms)
				},
			},
		},
	}).Add(&tfmapper.Int64FieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *int64 {
			return &in.QPSLimitCeiling
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthQpsLimitCeiling,
			Schema: &schema.Schema{
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          0,
				ValidateDiagFunc: mashschema.ValidateZeroOrGreater,
				Description:      "The throttle limit, i.e. calls per second, is applied to all access tokens granted for the API",
			},
		},
	}).Add(&tfmapper.Int64FieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *int64 {
			return &in.RateLimitCeiling
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthRateLimitCeiling,
			Schema: &schema.Schema{
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          0,
				ValidateDiagFunc: mashschema.ValidateZeroOrGreater,
				Description:      "The quota limit is applied to all access tokens granted for the API.",
			},
		},
	}).Add(&tfmapper.DurationFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *int64 {
			return &in.RefreshTokenTtl
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthRefreshTokenTtl,
			Schema: &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "768h",
				Description:      "The refresh token gets expired after it crosses the TTL value",
				ValidateDiagFunc: mashschema.ValidateDuration,
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.MasheryOAuth]{
		Locator: func(in *masherytypes.MasheryOAuth) *bool {
			return &in.SecureTokensEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.MasheryOAuth]{
			Key: mashschema.MashSvcOAuthSecureTokensEnabled,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, Mashery stores tokens using a one-way SHA-256 hashed value",
			},
		},
	})
}
