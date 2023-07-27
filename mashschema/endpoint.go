package mashschema

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var OutboundTransportProtocolEnum = []string{"use-inbound", "http", "https"}
var RequestProtocolEnum = []string{"rest", "soap", "xml-rpc", "json-rpc", "other"}
var ConnectionTimeoutForSystemDomainRequestEnum = []int{2, 5, 10, 20, 30, 45, 60} // Mashery has recently changed this
var ConnectionTimeoutForSystemDomainResponseEnum = []int{2, 5, 10, 20, 30, 45, 60, 120, 300, 600, 900, 1200}
var RequestAuthenticationTypeEnum = []string{"apiKey", "apiKeyAndSecret_MD5", "apiKeyAndSecret_SHA256", "secureHash_SHA256", "oauth", "custom"}

var RequestBodyVal = "request-body"
var RequestParametersVal = "request-parameters"

var DeveloperAPIKeyLocationsEnum = []string{"request-header", RequestBodyVal, RequestParametersVal, "request-path", "custom"}
var MethodLocationsEnum = []string{"request-header", RequestBodyVal, RequestParametersVal, "request-path", "custom"}
var UserControlledErrorFormatEnum = []string{"request-header", "request-body", "request-parameters", "resource"}

var ForwardedHeadersEnum = []string{"mashery-host", "mashery-message-id", "mashery-service-id"}
var ReturnedHeadersEnum = []string{"mashery-message-id", "mashery-responder"}
var HttoMethodsEnum = []string{"get", "post", "put", "delete", "head", "options", "patch"}

const (
	MashServiceEndpointId                   = "service_endpoint_id"
	MashServiceEndpointRef                  = "service_endpoint_ref"
	MashEndpointAllowMissingApiKey          = "allow_missing_api_key"
	MashEndpointApiKeyValueLocationKey      = "developer_api_key_field_name"
	MashEndpointApiKeyValueLocations        = "developer_api_key_locations"
	MashEndpointApiMethodDetectionKey       = "api_method_detection_key"
	MashEndpointApiMethodDetectionLocations = "api_method_detection_locations"

	MashEndpointCache                               = "cache"
	MashEndpointCacheClientSurrogateControlEnabled  = "client_surrogate_control_enabled"
	MashEndpointCacheContentCacheKeyHeaders         = "content_cache_key_headers"
	MashEndpointCacheTTLOverride                    = "cache_ttl_override"
	MashEndpointCacheIncludeApiKeyInContentCacheKey = "include_api_key_in_content_cache_key"
	MashEndpointCacheRespondFromStaleCacheEnabled   = "respond_from_stale_cache_enabled"
	MashEndpointCacheResponseCacheControlEnabled    = "response_cache_control_enabled"
	MashEndpointCacheVaryHeaderEnabled              = "vary_header_enabled"

	MashEndpointConnectionTimeoutForSystemDomainRequest  = "connection_timeout_for_system_domain_request"
	MashEndpointConnectionTimeoutForSystemDomainResponse = "connection_timeout_for_system_domain_response"
	MashEndpointCookiesDuringHttpRedirectsEnabled        = "cookies_during_http_redirects_enabled"

	MashEndpointCors                         = "cors"
	MashEndpointCorsAllDomainsEnabled        = "all_domains_enabled"
	MashEndpointCorsMaxAge                   = "max_age"
	MashEndpointCorsCookiesAllowed           = "cookies_allowed"
	MashEndpointCorsAllowedDomains           = "allowed_domains"
	MashEndpointCorsAllowedHeaders           = "allowed_headers"
	MashEndpointCorsExposedHeaders           = "exposed_headers"
	MashEndpointCorsSubDomainMatchingAllowed = "sub_domain_matching_allowed"

	MashEndpointCustomRequestAuthenticationAdapter         = "custom_request_authentication_adapter"
	MashEndpointDropApiKeyFromIncomingCall                 = "drop_api_key_from_incoming_call"
	MashEndpointForceGzipOfBackendCall                     = "force_gzip_of_backend_call"
	MashEndpointGzipPassthroughSupportEnabled              = "gzip_passthrough_support_enabled"
	MashEndpointHeadersToExcludeFromIncomingCall           = "headers_to_exclude_from_incoming_call"
	MashEndpointHighSecurity                               = "high_security"
	MashEndpointHostPassthroughIncludedInBackendCallHeader = "host_passthrough_included_in_backend_call_header"
	MashEndpointInboundSslRequired                         = "inbound_ssl_required"
	MashEndpointInboundMutualSslRequired                   = "inbound_mutual_ssl_required"
	MashEndpointJsonpCallbackParameter                     = "jsonp_callback_parameter"
	MashEndpointJsonpCallbackParameterValue                = "jsonp_callback_parameter_value"
	MashEndpointForwardedHeaders                           = "forwarded_headers"
	MashEndpointReturnedHeaders                            = "returned_headers"
	MashEndpointName                                       = "name"
	MashEndpointNumberOfHttpRedirectsToFollow              = "number_of_http_redirects_to_follow"
	MashEndpointOutboundRequestTargetPath                  = "outbound_request_target_path"
	MashEndpointOutboundRequestTargetQueryParameters       = "outbound_request_target_query_parameters"
	MashEndpointOutboundTransportProtocol                  = "outbound_transport_protocol"

	MashEndpointProcessor                   = "processor"
	MashEndpointProcessorAdapter            = "adapter"
	MashEndpointProcessorPreProcessEnabled  = "pre_process_enabled"
	MashEndpointProcessorPostProcessEnabled = "post_process_enabled"
	MashEndpointProcessorPreConfig          = "pre_config"
	MashEndpointProcessorPostConfig         = "post_config"

	MashEndpointPublicDomains             = "public_domains"
	MashEndpointRequestAuthenticationType = "request_authentication_type"
	MashEndpointRequestPathAlias          = "request_path_alias"
	MashEndpointRequestProtocol           = "request_protocol"
	MashEndpointOauthGrantTypes           = "oauth_grant_types"
	MashEndpointStringsToTrimFromApiKey   = "strings_to_trim_from_api_key"
	MashEndpointSupportedHttpMethods      = "supported_http_methods"

	MashEndpointSystemDomainAuthentication            = "system_domain_authentication"
	MashEndpointSystemDomainAuthenticationType        = "type"
	MashEndpointSystemDomainAuthenticationUsername    = "username"
	MashEndpointSystemDomainAuthenticationCertificate = "certificate"
	MashEndpointSystemDomainAuthenticationPassword    = "password"

	MashEndpointSystemDomains                = "system_domains"
	MashEndpointTrafficManagerDomain         = "traffic_manager_domain"
	MashEndpointUseSystemDomainCredentials   = "use_system_domain_credentials"
	MashEndpointSystemDomainCredentialKey    = "system_domain_credential_key"
	MashEndpointSystemDomainCredentialSecret = "system_domain_credential_secret"

	MashEndpointErrorSetRef                    = "error_set_ref"
	MashEndpointUserControlledErrorLocation    = "user_controlled_error_format_location"
	MashEndpointUserControlledErrorLocationKey = "user_controlled_error_format_location_key"
)

// Terraform specification for the Mashery endpoint.

var EndpointProcessorSchema = map[string]*schema.Schema{
	MashEndpointProcessorAdapter: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Adapter to be used",
	},
	MashEndpointProcessorPreProcessEnabled: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Pre-processor is enabled",
	},
	MashEndpointProcessorPostProcessEnabled: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Pre-processor is enabled",
	},
	MashEndpointProcessorPreConfig: {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     StringElem(),
	},
	MashEndpointProcessorPostConfig: {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     StringElem(),
	},
}

var EndpointSchema = map[string]*schema.Schema{
	MashSvcId: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Service Id to which this endpoint shall be attached",
	},
	MashServiceEndpointId: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Endpoint ID value generated by Mashery",
	},
	MashEndpointAllowMissingApiKey: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Set to true to allow calls without API keys",
	},
	MashEndpointApiKeyValueLocationKey: {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "api_key",
		Description: "Key in request parameters that identifies developer's api key",
	},
	MashEndpointApiKeyValueLocations: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Locations where the developer should place key",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashEndpointApiMethodDetectionKey: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Method v3Identity, aka the string that uniquely identifies the incoming method to the Traffic Manager. I",
	},
	MashEndpointApiMethodDetectionLocations: {
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true, // Would be assigned by Mashery to "request-path"
		Description: "Locations to derive method from",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashEndpointCache: {
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				MashEndpointCacheClientSurrogateControlEnabled: {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				MashEndpointCacheContentCacheKeyHeaders: {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	},
	MashEndpointConnectionTimeoutForSystemDomainRequest: {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     10,
		Description: "Timeout to connect to the back-end",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return ValidateIntValueInSet(i, path, &ConnectionTimeoutForSystemDomainRequestEnum)
		},
	},
	MashEndpointConnectionTimeoutForSystemDomainResponse: {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     60,
		Description: "Timeout to receive response from the back-end",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return ValidateIntValueInSet(i, path, &ConnectionTimeoutForSystemDomainResponseEnum)
		},
	},
	MashEndpointCookiesDuringHttpRedirectsEnabled: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Set to true to enable cookies during redirects",
	},
	MashEndpointCors: {
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				MashEndpointCorsAllDomainsEnabled: {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				MashEndpointCorsMaxAge: {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  300,
				},
			},
		},
	},
	MashObjCreated: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date/time the endpoint was created",
	},
	MashEndpointCustomRequestAuthenticationAdapter: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Custom adapter for authentication for an endpoint",
	},
	MashEndpointDropApiKeyFromIncomingCall: {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	MashEndpointForceGzipOfBackendCall: {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	MashEndpointGzipPassthroughSupportEnabled: {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	MashEndpointHeadersToExcludeFromIncomingCall: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "HTTP Headers to Drop from Incoming Call",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashEndpointHighSecurity: {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	MashEndpointHostPassthroughIncludedInBackendCallHeader: {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	MashEndpointInboundSslRequired: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Whether SSL is required when receiving traffic",
	},
	MashEndpointInboundMutualSslRequired: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Whether mutual SSL is required for this endpoint",
	},
	MashEndpointJsonpCallbackParameter: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The parameter is used by Traffic Manager while handling the JSON responses.",
	},
	MashEndpointJsonpCallbackParameterValue: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The default parameter value can be set to be used by Traffic Manager to effectively handle the JSON response.",
	},

	// ------------------------
	// Scheduled maintenance event is a separate resource. The presence of the scheduled event could be
	// queried from the data source.
	// --------------------------
	MashEndpointForwardedHeaders: {
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Description: "Specific Mashery headers to be forwarded to the back-end",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashEndpointReturnedHeaders: {
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true, // Would be assigned by Mashery to "mashery-responder"
		Description: "Specific Mashery headers to be forwarded to the back-end",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// ---------------------------
	// Methods are separate resources. These are managed by the resources. All loaded resource could be
	// queried
	// ------------------------------
	MashEndpointName: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Endpoint name",
	},
	MashEndpointNumberOfHttpRedirectsToFollow: {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Number of HTTP redirects to follow",
	},
	MashEndpointOutboundRequestTargetPath: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Root context of the outbound domain to serve this data from",
	},
	MashEndpointOutboundRequestTargetQueryParameters: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Query string to be added to the original request",
	},
	MashEndpointOutboundTransportProtocol: {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "https",
		Description: "Outbound request protocol, defaults to https",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return ValidateStringValueInSet(i, path, &OutboundTransportProtocolEnum)
		},
	},
	MashEndpointProcessor: {
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: EndpointProcessorSchema,
		},
	},
	// Public domains use a bit different mashschema. For simplicity of a Terraform developer, the
	// domains are specified as a list.
	MashEndpointPublicDomains: {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// Mashery V3 API has request authentication type as optional. From the concept of Terraform developent,
	// the request authentication type is extrementy important and must be supplied.
	MashEndpointRequestAuthenticationType: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Authentication type for the endpoint",
		// TODO: add constriant on the authentication options
	},
	MashEndpointRequestPathAlias: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Request path",
	},
	MashEndpointRequestProtocol: {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "rest",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return ValidateStringValueInSet(i, path, &RequestProtocolEnum)
		},
	},
	MashEndpointOauthGrantTypes: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "OAuth grant types supported at this endpoint",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashEndpointStringsToTrimFromApiKey: {
		Type:     schema.TypeString,
		Optional: true,
	},
	// Another deviation from V3: the developer must declare the supported methods.
	MashEndpointSupportedHttpMethods: {
		Type:        schema.TypeSet,
		Required:    true,
		Description: "Methods this endpoint will support",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashEndpointSystemDomainAuthentication: {
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				MashEndpointSystemDomainAuthenticationType: {
					Type:     schema.TypeString,
					Required: true,
				},
				MashEndpointSystemDomainAuthenticationUsername: {
					Type:     schema.TypeString,
					Optional: true,
					// TODO: How to use AtLeastOneOf???
					//AtLeastOneOf: []string{"username", "certificate"},
				},
				MashEndpointSystemDomainAuthenticationCertificate: {
					Type:     schema.TypeString,
					Optional: true,
					//AtLeastOneOf: []string{"username", "certificate"},
				},
				MashEndpointSystemDomainAuthenticationPassword: {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
	MashEndpointSystemDomains: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The domain name of the client API server",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashEndpointTrafficManagerDomain: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Traffic Manager internal hostname (domain name) to which the requested public hostname is CNAMED.",
	},
	MashObjUpdated: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date/time the object was last updated.",
	},
	MashEndpointUseSystemDomainCredentials: {
		Type:        schema.TypeBool,
		Description: "To suit to the client server's requirement, Mashery can swap the API credentials such as API keys and send the swapped Mashery credentials to the client API server.",
		Optional:    true,
		Default:     false,
	},
	MashEndpointSystemDomainCredentialKey: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Key to use when making call to the client API server",
	},
	MashEndpointSystemDomainCredentialSecret: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Secret to use when making call to the client API server.",
	},
}

// ----------------------------------------------------------------------------------------------
// Terraform - V3 conversion functions.

var impliedKeyValueLocations = []string{"request-header"}
var impliedMethodDetectionLocations = []string{"request-path"}
