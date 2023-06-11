package mashschema

// ---------------------------------------------------------
// Schemas for Mashery service resource and data sources.

// Unused args in mashschema -> to be moved to the provider resource
var supportedForwardedHeaders = []string{
	"access-token", "client-id", "scope", "user-context",
}

var supportedMasheryGrantTypes = []string{"authorization_code", "implicit", "password", "client_credentials"}
var SupportedMasheryMacAlgorithms = []string{"hmac-sha-1", "hmac-sha-256"}

const (
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
	MashSvcInteractiveDocsRoles = "iodocs_accessed_by"
	MashSvcOrganization         = "organization"

	// MashSvcId Mashery Service element keys
	MashSvcId                   = "service_id"
	MashSvcRef                  = "service_ref"
	MashSvcName                 = "name"
	MashSvcEditorHandle         = "editor_handle"
	MashSvcRevisionNumber       = "revision_number"
	MashSvcRobotsPolicy         = "robots_policy"
	MashSvcCrossdomainPolicy    = "crossdomain_policy"
	MashSvcDescription          = "description"
	MashSvcQpsLimitOverall      = "qps_limit_overall"
	MashSvcServiceRFC3986Encode = "rfc3986_encode"
	MashSvcVersion              = "version"

	MashSvcCacheTtl = "cache_ttl"
)
