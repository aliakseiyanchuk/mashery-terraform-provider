package mashery

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var outboundTransportProtocolEnum = []string{"use-inbound", "http", "https"}
var requestProtocolEnum = []string{"rest", "soap", "xml-rpc", "json-rpc", "other"}
var connectionTimeoutForSystemDomainRequestEnum = []int{2, 5, 10, 20, 30, 45, 60}
var connectionTimeoutForSystemDomainResponseEnum = []int{2, 5, 10, 20, 30, 45, 60, 120, 300, 600, 900, 1200}

const (
	MashEndpointId                          = "endpoint_id"
	MashEndpointAllowMissingApiKey          = "allow_missing_api_key"
	MashEndpointApiKeyValueLocationKey      = "api_key_value_location_key"
	MashEndpointApiKeyValueLocations        = "api_key_value_locations"
	MashEndpointApiMethodDetectionKey       = "api_method_detection_key"
	MashEndpointApiMethodDetectionLocations = "api_method_detection_locations"

	MashEndpointCache                              = "cache"
	MashEndpointCacheClientSurrogateControlEnabled = "client_surrogate_control_enabled"
	MashEndpointCacheContentCacheKeyHeaders        = "content_cache_key_headers"

	MashEndpointConnectionTimeoutForSystemDomainRequest  = "connection_timeout_for_system_domain_request"
	MashEndpointConnectionTimeoutForSystemDomainResponse = "connection_timeout_for_system_domain_response"
	MashEndpointCookiesDuringHttpRedirectsEnabled        = "cookies_during_http_redirects_enabled"

	MashEndpointCors                  = "cors"
	MashEndpointCorsAllDomainsEnabled = "all_domains_enabled"
	MashEndpointCorsMaxAge            = "max_age"

	MashEndpointCreated                                    = "created"
	MashEndpointCustomRequestAuthenticationAdapter         = "custom_request_authentication_adapter"
	MashEndpointDropApiKeyFromIncomingCall                 = "drop_api_key_from_incoming_call"
	MashEndpointForceGzipOfBackendCall                     = "force_gzip_of_backend_call"
	MashEndpointGzipPassthroughSupportEnabled              = "gzip_passthrough_support_enabled"
	MashEndpointHeadersToExcludeFromIncomingCall           = "headers_to_exclude_from_incoming_call"
	MashEndpointHighSecurity                               = "high_security"
	MashEndpointHostPassthroughIncludedInBackendCallHeader = "host_passthrough_included_in_backend_call_header"
	MashEndpointInboundSslRequired                         = "inbound_ssl_required"
	MashEndpointInboundMutualSslRequired                   = "inbound_mutila_ssl_required"
	MashEndpointJsonpCallbackParameter                     = "jsonp_callback_parameter"
	MashEndpointJsonpCallbackParameterValue                = "jsonp_callback_parameter_value"
	MashEndpointScheduledMaintenanceEvent                  = "scheduled_maintenance_event"
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
	MashEndpointUpdated                      = "updated"
	MashEndpointUseSystemDomainCredentials   = "use_system_domain_credentials"
	MashEndpointSystemDomainCredentialKey    = "system_domain_credential_key"
	MashEndpointSystemDomainCredentialSecret = "system_domain_credential_secret"

	MashEndpointMultiRef                = "endpoint_ids"
	MashEndpointsExplained              = "endpoints_explained"
	DataSourceServiceEndpointPathRegexp = "filter_request_path_alias"
)

// Terraform specification for the Mashery endpoint.

type CompoundIdentifier struct {
}

func (ci *CompoundIdentifier) MalformedDiagnostic(path string) diag.Diagnostics {
	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       "Incomplete id",
		Detail:        "Identifier supplies incomplete data or is malformed",
		AttributePath: cty.GetAttrPath(path),
	}}
}

type ServiceEndpointIdentifier struct {
	CompoundIdentifier
	ServiceId  string
	EndpointId string
}

func (ei *ServiceEndpointIdentifier) IsIdentified() bool {
	return len(ei.ServiceId) > 0 && len(ei.EndpointId) > 0
}

func (ei *ServiceEndpointIdentifier) Id() string {
	return CreateCompoundId(ei.ServiceId, ei.EndpointId)
}

func (ei *ServiceEndpointIdentifier) IdOf(comment string) string {
	return fmt.Sprintf("%s # %s", ei.Id(), comment)
}

func (ei *ServiceEndpointIdentifier) From(id string) {
	ParseCompoundId(id, &ei.ServiceId, &ei.EndpointId)
}

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
		Type:     schema.TypeList,
		Optional: true,
		Elem:     stringElem(),
	},
	MashEndpointProcessorPostConfig: {
		Type:     schema.TypeList,
		Optional: true,
		Elem:     stringElem(),
	},
}

var EndpointSchema = map[string]*schema.Schema{
	MashSvcId: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Service Id to which this endpoint shall be attached",
	},
	MashEndpointId: {
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
		Description: "Method identifier, aka the string that uniquely identifies the incoming method to the Traffic Manager. I",
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
			return validateIntValueInSet(i, path, &connectionTimeoutForSystemDomainRequestEnum)
		},
	},
	MashEndpointConnectionTimeoutForSystemDomainResponse: {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     60,
		Description: "Timeout to receive response from the back-end",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return validateIntValueInSet(i, path, &connectionTimeoutForSystemDomainResponseEnum)
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
	MashEndpointCreated: {
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
			return validateStringValueInSet(i, path, &outboundTransportProtocolEnum)
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
	// Public domains use a bit different schema. For simplicity of a Terraform developer, the
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
			return validateStringValueInSet(i, path, &requestProtocolEnum)
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
	MashEndpointUpdated: {
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

var DataSourceSvcEndpointsSchema = DataSourceBaseSchema()

func initDataSourceSvcEndpointsSchema() {
	DataSourceSvcEndpointsSchema[MashEndpointMultiRef] = &schema.Schema{
		Computed:    true,
		Type:        schema.TypeSet,
		Description: "Matched endpoints",
		Elem:        stringElem(),
	}
	DataSourceSvcEndpointsSchema[DataSourceServiceEndpointPathRegexp] = &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Regular expressions to match",
		//ValidateDiagFunc: validateRegularExpressionSet,
		Elem: stringElem(),
	}
	DataSourceSvcEndpointsSchema[MashEndpointsExplained] = &schema.Schema{
		Type:        schema.TypeMap,
		Computed:    true,
		Description: "Mapping from endpoint ids to endpoint names",
		Elem:        stringElem(),
	}
	addRequiredString(&DataSourceSvcEndpointsSchema, MashSvcId, "Service to search")
}

// ----------------------------------------------------------------------------------------------
// Terraform - V3 conversion functions.

var impliedKeyValueLocations = []string{"request-header"}
var impliedMethodDetectionLocations = []string{"request-path"}

func V3ProcessorConfigurationToTerraform(inp masherytypes.Processor) map[string]interface{} {
	return map[string]interface{}{
		MashEndpointProcessorAdapter:            inp.Adapter,
		MashEndpointProcessorPreProcessEnabled:  inp.PreProcessEnabled,
		MashEndpointProcessorPostProcessEnabled: inp.PostProcessEnabled,
		MashEndpointProcessorPreConfig:          inp.PreInputs,
		MashEndpointProcessorPostConfig:         inp.PostInputs,
	}
}

func V3ProcessorConfigurationFrom(inp map[string]interface{}) masherytypes.Processor {
	return masherytypes.Processor{
		Adapter:            inp[MashEndpointProcessorAdapter].(string),
		PreProcessEnabled:  inp[MashEndpointProcessorPreProcessEnabled].(bool),
		PostProcessEnabled: inp[MashEndpointProcessorPostProcessEnabled].(bool),
		PreInputs:          ConvertInterfaceArrayToStringArray(inp[MashEndpointProcessorPreConfig].([]interface{})),
		PostInputs:         ConvertInterfaceArrayToStringArray(inp[MashEndpointProcessorPostConfig].([]interface{})),
	}
}

func mashEndpointCacheUpsertable(d *schema.ResourceData) *masherytypes.Cache {
	if cacheSet, exists := d.GetOk(MashEndpointCache); exists {
		tfCache := unwrapStructFromTerraformSet(cacheSet)

		// TODO: unsafe lookups from the map. Should be extracted to avoid panic
		rv := masherytypes.Cache{
			ClientSurrogateControlEnabled: tfCache[MashEndpointCacheClientSurrogateControlEnabled].(bool),
			ContentCacheKeyHeaders:        schemaSetToStringArray(tfCache[MashEndpointCacheContentCacheKeyHeaders]),
		}

		return &rv
	} else {
		return nil
	}
}

func mashEndpointCorsUpsertable(d *schema.ResourceData) *masherytypes.Cors {
	if corsSet, exists := d.GetOk(MashEndpointCors); exists {
		tfCors := unwrapStructFromTerraformSet(corsSet)

		rv := masherytypes.Cors{
			AllDomainsEnabled: tfCors[MashEndpointCorsAllDomainsEnabled].(bool),
			MaxAge:            tfCors[MashEndpointCorsMaxAge].(int),
		}

		return &rv
	} else {
		return nil
	}
}

func mashEndpointProcessorUpsertable(d *schema.ResourceData) *masherytypes.Processor {
	if procSet, exists := d.GetOk(MashEndpointProcessor); exists {
		tfProcMap := unwrapStructFromTerraformSet(procSet)

		// TODO: Unsafe lookup from map
		// Maybe this conversion won't event work.
		rv := masherytypes.Processor{
			PreProcessEnabled:  tfProcMap[MashEndpointProcessorPreProcessEnabled].(bool),
			PostProcessEnabled: tfProcMap[MashEndpointProcessorPostProcessEnabled].(bool),
			//PreInputs:          schemaMapToStringMap(tfProcMap[MashEndpointProcessorPreConfig]),
			//PostInputs:         schemaMapToStringMap(tfProcMap[MashEndpointProcessorPostConfig]),
			Adapter: tfProcMap[MashEndpointProcessorAdapter].(string),
		}

		if rv.IsEmpty() {
			doLogf("Processor for this endpoint id=%s is effectively empty", afterApplyKnown(d.Id()))
			return nil
		} else {
			return &rv
		}
	} else {
		return nil
	}
}

func mashEndpointDomainsArrayUpsertable(d *schema.ResourceData, key string) []masherytypes.Domain {
	defDomains := ExtractStringArray(d, key, &EmptyStringArray)

	rv := make([]masherytypes.Domain, len(defDomains))
	for idx, v := range defDomains {
		rv[idx] = masherytypes.Domain{Address: v}
	}

	return rv
}

func mashEndpointSystemDomainAuthenticationUpsertable(d *schema.ResourceData) *masherytypes.SystemDomainAuthentication {
	if sysDomainSet, exists := d.GetOk(MashEndpointSystemDomainAuthentication); exists {
		tfSysDomain := unwrapStructFromTerraformSet(sysDomainSet)

		// TODO: Unsafe lookups
		rv := masherytypes.SystemDomainAuthentication{
			Type:        tfSysDomain[MashEndpointSystemDomainAuthenticationType].(string),
			Username:    safeLookupStringPointer(&tfSysDomain, MashEndpointSystemDomainAuthenticationUsername),
			Certificate: safeLookupStringPointer(&tfSysDomain, MashEndpointSystemDomainAuthenticationCertificate),
			Password:    safeLookupStringPointer(&tfSysDomain, MashEndpointSystemDomainAuthenticationPassword),
		}

		return &rv
	} else {
		return nil
	}
}

// Create V3 Mashery Endpoint data structure from the resource data
func MashEndpointUpsertable(d *schema.ResourceData) (masherytypes.MasheryEndpoint, diag.Diagnostics) {

	enpdIdent := ServiceEndpointIdentifier{}
	enpdIdent.From(d.Id())

	rv := masherytypes.MasheryEndpoint{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   enpdIdent.EndpointId,
			Name: extractString(d, MashEndpointName, resource.UniqueId()),
		},
		AllowMissingApiKey:                         extractBool(d, MashEndpointAllowMissingApiKey, false),
		ApiKeyValueLocationKey:                     extractString(d, MashEndpointApiKeyValueLocationKey, "api_key"),
		ApiKeyValueLocations:                       ExtractStringArray(d, MashEndpointApiKeyValueLocations, &impliedKeyValueLocations),
		ApiMethodDetectionKey:                      extractString(d, MashEndpointApiMethodDetectionKey, ""),
		ApiMethodDetectionLocations:                ExtractStringArray(d, MashEndpointApiMethodDetectionLocations, &impliedMethodDetectionLocations),
		Cache:                                      mashEndpointCacheUpsertable(d),
		ConnectionTimeoutForSystemDomainRequest:    extractInt(d, MashEndpointConnectionTimeoutForSystemDomainRequest, 10),
		ConnectionTimeoutForSystemDomainResponse:   extractInt(d, MashEndpointConnectionTimeoutForSystemDomainResponse, 60),
		CookiesDuringHttpRedirectsEnabled:          extractBool(d, MashEndpointCookiesDuringHttpRedirectsEnabled, false),
		Cors:                                       mashEndpointCorsUpsertable(d),
		CustomRequestAuthenticationAdapter:         ExtractStringPointer(d, MashEndpointCustomRequestAuthenticationAdapter),
		DropApiKeyFromIncomingCall:                 extractBool(d, MashEndpointDropApiKeyFromIncomingCall, true),
		ForceGzipOfBackendCall:                     extractBool(d, MashEndpointForceGzipOfBackendCall, false),
		GzipPassthroughSupportEnabled:              extractBool(d, MashEndpointGzipPassthroughSupportEnabled, true),
		HeadersToExcludeFromIncomingCall:           ExtractStringArray(d, MashEndpointHeadersToExcludeFromIncomingCall, &EmptyStringArray),
		HighSecurity:                               extractBool(d, MashEndpointHighSecurity, false),
		HostPassthroughIncludedInBackendCallHeader: extractBool(d, MashEndpointHostPassthroughIncludedInBackendCallHeader, true),
		InboundSslRequired:                         extractBool(d, MashEndpointInboundSslRequired, true),
		InboundMutualSslRequired:                   extractBool(d, MashEndpointInboundMutualSslRequired, true),
		JsonpCallbackParameter:                     extractString(d, MashEndpointJsonpCallbackParameter, ""),
		JsonpCallbackParameterValue:                extractString(d, MashEndpointJsonpCallbackParameterValue, ""),
		ScheduledMaintenanceEvent:                  nil,
		ForwardedHeaders:                           ExtractStringArray(d, MashEndpointForwardedHeaders, &EmptyStringArray),
		ReturnedHeaders:                            ExtractStringArray(d, MashEndpointReturnedHeaders, &EmptyStringArray),
		Methods:                                    nil,
		NumberOfHttpRedirectsToFollow:              extractInt(d, MashEndpointNumberOfHttpRedirectsToFollow, 0),
		OutboundRequestTargetPath:                  extractString(d, MashEndpointOutboundRequestTargetPath, ""),
		OutboundRequestTargetQueryParameters:       extractString(d, MashEndpointOutboundRequestTargetQueryParameters, ""),
		OutboundTransportProtocol:                  extractString(d, MashEndpointOutboundTransportProtocol, "rest"),
		Processor:                                  mashEndpointProcessorUpsertable(d),
		PublicDomains:                              mashEndpointDomainsArrayUpsertable(d, MashEndpointPublicDomains),
		RequestAuthenticationType:                  extractString(d, MashEndpointRequestAuthenticationType, "public"),
		RequestPathAlias:                           extractString(d, MashEndpointRequestPathAlias, "/"),
		RequestProtocol:                            extractString(d, MashEndpointRequestProtocol, "rest"),
		OAuthGrantTypes:                            ExtractStringArray(d, MashEndpointOauthGrantTypes, &EmptyStringArray),
		StringsToTrimFromApiKey:                    extractString(d, MashEndpointStringsToTrimFromApiKey, ""),
		SupportedHttpMethods:                       ExtractStringArray(d, MashEndpointSupportedHttpMethods, &EmptyStringArray),
		SystemDomainAuthentication:                 mashEndpointSystemDomainAuthenticationUpsertable(d),
		SystemDomains:                              mashEndpointDomainsArrayUpsertable(d, MashEndpointSystemDomains),
		TrafficManagerDomain:                       extractString(d, MashEndpointTrafficManagerDomain, ""),
		UseSystemDomainCredentials:                 extractBool(d, MashEndpointUseSystemDomainCredentials, false),
		SystemDomainCredentialKey:                  ExtractStringPointer(d, MashEndpointSystemDomainCredentialKey),
		SystemDomainCredentialSecret:               ExtractStringPointer(d, MashEndpointSystemDomainCredentialSecret),
	}

	// TODO: Perform diagnostics for offending values.

	rvDiag := diag.Diagnostics{}
	return rv, rvDiag
}

func v3EndpointCacheToTerraform(cache *masherytypes.Cache) []interface{} {
	if !cache.IsEmpty() {
		return []interface{}{
			map[string]interface{}{
				MashEndpointCacheClientSurrogateControlEnabled: cache.ClientSurrogateControlEnabled,
				MashEndpointCacheContentCacheKeyHeaders:        cache.ContentCacheKeyHeaders,
			},
		}
	} else {
		return nil
	}
}

func v3EndpointCorsToTerraform(cors *masherytypes.Cors) []interface{} {
	if cors != nil {
		return []interface{}{
			map[string]interface{}{
				MashEndpointCorsAllDomainsEnabled: cors.AllDomainsEnabled,
				MashEndpointCorsMaxAge:            cors.MaxAge,
			},
		}
	} else {
		return nil
	}
}

func v3EndpointProcessorToTerraform(processor *masherytypes.Processor) []interface{} {
	// We will persist the processor in the state only if will contain useful information.
	// Mashery V3 API will reuturn processor data structure also when the call transformation is not enabled.
	if processor != nil && !processor.IsEmpty() {
		rv := []interface{}{
			map[string]interface{}{
				MashEndpointProcessorAdapter:            processor.Adapter,
				MashEndpointProcessorPreProcessEnabled:  processor.PreProcessEnabled,
				MashEndpointProcessorPostProcessEnabled: processor.PostProcessEnabled,
				MashEndpointProcessorPreConfig:          processor.PreInputs,
				MashEndpointProcessorPostConfig:         processor.PostInputs,
			},
		}

		doLogJson("Converted processor data structure for TF", rv)

		return rv
	} else {
		return nil
	}
}

func v3EndpointDomainsToTerraform(d []masherytypes.Domain) []string {
	if d != nil {
		rv := make([]string, len(d))

		for idx, v := range d {
			rv[idx] = v.Address
		}

		return rv
	} else {
		return []string{}
	}
}

func v3EndpointSystemDomainAuthenticationToTerraform(d *masherytypes.SystemDomainAuthentication) []interface{} {
	if d != nil {
		return []interface{}{
			map[string]interface{}{
				MashEndpointSystemDomainAuthenticationType:        d.Type,
				MashEndpointSystemDomainAuthenticationUsername:    d.Username,
				MashEndpointSystemDomainAuthenticationCertificate: d.Certificate,
				MashEndpointSystemDomainAuthenticationPassword:    d.Password,
			},
		}
	} else {
		return nil
	}
}

func V3EndpointToResourceData(endpoint *masherytypes.MasheryEndpoint, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashEndpointId:                                         endpoint.Id,
		MashEndpointAllowMissingApiKey:                         endpoint.AllowMissingApiKey,
		MashEndpointApiKeyValueLocationKey:                     endpoint.ApiKeyValueLocationKey,
		MashEndpointApiKeyValueLocations:                       endpoint.ApiKeyValueLocations,
		MashEndpointApiMethodDetectionKey:                      endpoint.ApiMethodDetectionKey,
		MashEndpointApiMethodDetectionLocations:                endpoint.ApiMethodDetectionLocations,
		MashEndpointCache:                                      v3EndpointCacheToTerraform(endpoint.Cache),
		MashEndpointConnectionTimeoutForSystemDomainRequest:    endpoint.ConnectionTimeoutForSystemDomainRequest,
		MashEndpointConnectionTimeoutForSystemDomainResponse:   endpoint.ConnectionTimeoutForSystemDomainResponse,
		MashEndpointCookiesDuringHttpRedirectsEnabled:          endpoint.CookiesDuringHttpRedirectsEnabled,
		MashEndpointCors:                                       v3EndpointCorsToTerraform(endpoint.Cors),
		MashEndpointCreated:                                    endpoint.Created.ToString(),
		MashEndpointCustomRequestAuthenticationAdapter:         endpoint.CustomRequestAuthenticationAdapter,
		MashEndpointDropApiKeyFromIncomingCall:                 endpoint.DropApiKeyFromIncomingCall,
		MashEndpointForceGzipOfBackendCall:                     endpoint.ForceGzipOfBackendCall,
		MashEndpointGzipPassthroughSupportEnabled:              endpoint.GzipPassthroughSupportEnabled,
		MashEndpointHeadersToExcludeFromIncomingCall:           endpoint.HeadersToExcludeFromIncomingCall,
		MashEndpointHighSecurity:                               endpoint.HighSecurity,
		MashEndpointHostPassthroughIncludedInBackendCallHeader: endpoint.HostPassthroughIncludedInBackendCallHeader,
		MashEndpointInboundSslRequired:                         endpoint.InboundSslRequired,
		MashEndpointJsonpCallbackParameter:                     endpoint.JsonpCallbackParameter,
		MashEndpointJsonpCallbackParameterValue:                endpoint.JsonpCallbackParameterValue,
		MashEndpointForwardedHeaders:                           endpoint.ForwardedHeaders,
		MashEndpointReturnedHeaders:                            endpoint.ReturnedHeaders,
		MashEndpointName:                                       endpoint.Name,
		MashEndpointNumberOfHttpRedirectsToFollow:              endpoint.NumberOfHttpRedirectsToFollow,
		MashEndpointOutboundRequestTargetPath:                  endpoint.OutboundRequestTargetPath,
		MashEndpointOutboundRequestTargetQueryParameters:       endpoint.OutboundRequestTargetQueryParameters,
		MashEndpointOutboundTransportProtocol:                  endpoint.OutboundTransportProtocol,
		MashEndpointProcessor:                                  v3EndpointProcessorToTerraform(endpoint.Processor),
		MashEndpointPublicDomains:                              v3EndpointDomainsToTerraform(endpoint.PublicDomains),
		MashEndpointRequestAuthenticationType:                  endpoint.RequestAuthenticationType,
		MashEndpointRequestPathAlias:                           endpoint.RequestPathAlias,
		MashEndpointRequestProtocol:                            endpoint.RequestProtocol,
		MashEndpointOauthGrantTypes:                            endpoint.OAuthGrantTypes,
		MashEndpointStringsToTrimFromApiKey:                    endpoint.StringsToTrimFromApiKey,
		MashEndpointSupportedHttpMethods:                       endpoint.SupportedHttpMethods,
		MashEndpointSystemDomainAuthentication:                 v3EndpointSystemDomainAuthenticationToTerraform(endpoint.SystemDomainAuthentication),
		MashEndpointSystemDomains:                              v3EndpointDomainsToTerraform(endpoint.SystemDomains),
		MashEndpointTrafficManagerDomain:                       endpoint.TrafficManagerDomain,
		MashEndpointUpdated:                                    endpoint.Updated.ToString(),
		MashEndpointUseSystemDomainCredentials:                 endpoint.UseSystemDomainCredentials,
		MashEndpointSystemDomainCredentialKey:                  endpoint.SystemDomainCredentialKey,
		MashEndpointSystemDomainCredentialSecret:               endpoint.SystemDomainCredentialSecret,
	}

	doLogJson("Converted endpoint to TF", &data)

	return SetResourceFields(data, d)
}

// init
func init() {
	initDataSourceSvcEndpointsSchema()
}
