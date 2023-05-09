package mashschema

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
	MashEndpointUseSystemDomainCredentials   = "use_system_domain_credentials"
	MashEndpointSystemDomainCredentialKey    = "system_domain_credential_key"
	MashEndpointSystemDomainCredentialSecret = "system_domain_credential_secret"

	MashEndpointMultiRef                = "endpoint_ids"
	MashEndpointsExplained              = "endpoints_explained"
	DataSourceServiceEndpointPathRegexp = "filter_request_path_alias"
)

// Terraform specification for the Mashery endpoint.

var ServiceEndpointMapper *ServiceEndpointMapperImpl

type ServiceEndpointMapperImpl struct {
	ResourceMapperImpl
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
			return ValidateStringValueInSet(i, path, &outboundTransportProtocolEnum)
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
			return ValidateStringValueInSet(i, path, &requestProtocolEnum)
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

func (sem *ServiceEndpointMapperImpl) GetServiceId(d *schema.ResourceData) string {
	return d.Get(MashSvcId).(string)
}

func (sem *ServiceEndpointMapperImpl) cacheUpsertable(d *schema.ResourceData) *masherytypes.Cache {
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

func (sem *ServiceEndpointMapperImpl) corsUpsertable(d *schema.ResourceData) *masherytypes.Cors {
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

func (sem *ServiceEndpointMapperImpl) processorUpsertable(d *schema.ResourceData) *masherytypes.Processor {
	if procSet, exists := d.GetOk(MashEndpointProcessor); exists {
		tfProcMap := unwrapStructFromTerraformSet(procSet)

		rv := sem.processorUpsertableFromMap(tfProcMap)

		if rv.IsEmpty() {
			return nil
		} else {
			return &rv
		}
	} else {
		return nil
	}
}

func (sem *ServiceEndpointMapperImpl) schemaArrayToString(in []interface{}) []string {
	rv := make([]string, len(in))
	for idx, vRaw := range in {
		if v, ok := vRaw.(string); ok {
			rv[idx] = v
		} else {
			rv[idx] = fmt.Sprintf("%s", vRaw)
		}
	}

	return rv
}

func (sem *ServiceEndpointMapperImpl) processorUpsertableFromMap(tfProcMap map[string]interface{}) masherytypes.Processor {
	// TODO: Unsafe lookup from map
	// Maybe this conversion won't event work.
	rv := masherytypes.Processor{
		PreProcessEnabled:  tfProcMap[MashEndpointProcessorPreProcessEnabled].(bool),
		PostProcessEnabled: tfProcMap[MashEndpointProcessorPostProcessEnabled].(bool),
		//PreInputs:          sem.schemaArrayToString(tfProcMap[MashEndpointProcessorPreConfig].([]interface{})),
		//PostInputs:         sem.schemaArrayToString(tfProcMap[MashEndpointProcessorPostConfig].([]interface{})),
		Adapter: tfProcMap[MashEndpointProcessorAdapter].(string),
	}
	return rv
}

func (sem *ServiceEndpointMapperImpl) domainsArrayUpsertable(d *schema.ResourceData, key string) []masherytypes.Domain {
	defDomains := ExtractStringArray(d, key, &EmptyStringArray)

	rv := make([]masherytypes.Domain, len(defDomains))
	for idx, v := range defDomains {
		rv[idx] = masherytypes.Domain{Address: v}
	}

	return rv
}

func (sem *ServiceEndpointMapperImpl) systemDomainAuthenticationUpsertable(d *schema.ResourceData) *masherytypes.SystemDomainAuthentication {
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

// UpsertableTyped Create V3 Mashery Endpoint data structure from the resource data
func (sem *ServiceEndpointMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.Endpoint, masherytypes.ServiceIdentifier, diag.Diagnostics) {

	rvd := diag.Diagnostics{}

	enpdIdent := masherytypes.ServiceEndpointIdentifier{}
	CompoundIdFrom(&enpdIdent, d.Id())

	serviceIdent := masherytypes.ServiceIdentifier{}
	if !CompoundIdFrom(&serviceIdent, ExtractString(d, MashSvcId, "")) {
		rvd = append(rvd, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "service context is incomplete",
		})
	}

	rv := masherytypes.Endpoint{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   enpdIdent.EndpointId,
			Name: ExtractString(d, MashEndpointName, resource.UniqueId()),
		},
		AllowMissingApiKey:                         ExtractBool(d, MashEndpointAllowMissingApiKey, false),
		ApiKeyValueLocationKey:                     ExtractString(d, MashEndpointApiKeyValueLocationKey, "api_key"),
		ApiKeyValueLocations:                       ExtractStringArray(d, MashEndpointApiKeyValueLocations, &impliedKeyValueLocations),
		ApiMethodDetectionKey:                      ExtractString(d, MashEndpointApiMethodDetectionKey, ""),
		ApiMethodDetectionLocations:                ExtractStringArray(d, MashEndpointApiMethodDetectionLocations, &impliedMethodDetectionLocations),
		Cache:                                      sem.cacheUpsertable(d),
		ConnectionTimeoutForSystemDomainRequest:    extractInt(d, MashEndpointConnectionTimeoutForSystemDomainRequest, 10),
		ConnectionTimeoutForSystemDomainResponse:   extractInt(d, MashEndpointConnectionTimeoutForSystemDomainResponse, 60),
		CookiesDuringHttpRedirectsEnabled:          ExtractBool(d, MashEndpointCookiesDuringHttpRedirectsEnabled, false),
		Cors:                                       sem.corsUpsertable(d),
		CustomRequestAuthenticationAdapter:         ExtractStringPointer(d, MashEndpointCustomRequestAuthenticationAdapter),
		DropApiKeyFromIncomingCall:                 ExtractBool(d, MashEndpointDropApiKeyFromIncomingCall, true),
		ForceGzipOfBackendCall:                     ExtractBool(d, MashEndpointForceGzipOfBackendCall, false),
		GzipPassthroughSupportEnabled:              ExtractBool(d, MashEndpointGzipPassthroughSupportEnabled, true),
		HeadersToExcludeFromIncomingCall:           ExtractStringArray(d, MashEndpointHeadersToExcludeFromIncomingCall, &EmptyStringArray),
		HighSecurity:                               ExtractBool(d, MashEndpointHighSecurity, false),
		HostPassthroughIncludedInBackendCallHeader: ExtractBool(d, MashEndpointHostPassthroughIncludedInBackendCallHeader, true),
		InboundSslRequired:                         ExtractBool(d, MashEndpointInboundSslRequired, true),
		InboundMutualSslRequired:                   ExtractBool(d, MashEndpointInboundMutualSslRequired, true),
		JsonpCallbackParameter:                     ExtractString(d, MashEndpointJsonpCallbackParameter, ""),
		JsonpCallbackParameterValue:                ExtractString(d, MashEndpointJsonpCallbackParameterValue, ""),
		ScheduledMaintenanceEvent:                  nil,
		ForwardedHeaders:                           ExtractStringArray(d, MashEndpointForwardedHeaders, &EmptyStringArray),
		ReturnedHeaders:                            ExtractStringArray(d, MashEndpointReturnedHeaders, &EmptyStringArray),
		Methods:                                    nil,
		NumberOfHttpRedirectsToFollow:              extractInt(d, MashEndpointNumberOfHttpRedirectsToFollow, 0),
		OutboundRequestTargetPath:                  ExtractString(d, MashEndpointOutboundRequestTargetPath, ""),
		OutboundRequestTargetQueryParameters:       ExtractString(d, MashEndpointOutboundRequestTargetQueryParameters, ""),
		OutboundTransportProtocol:                  ExtractString(d, MashEndpointOutboundTransportProtocol, "rest"),
		Processor:                                  sem.processorUpsertable(d),
		PublicDomains:                              sem.domainsArrayUpsertable(d, MashEndpointPublicDomains),
		RequestAuthenticationType:                  ExtractString(d, MashEndpointRequestAuthenticationType, "public"),
		RequestPathAlias:                           ExtractString(d, MashEndpointRequestPathAlias, "/"),
		RequestProtocol:                            ExtractString(d, MashEndpointRequestProtocol, "rest"),
		OAuthGrantTypes:                            ExtractStringArray(d, MashEndpointOauthGrantTypes, &EmptyStringArray),
		StringsToTrimFromApiKey:                    ExtractString(d, MashEndpointStringsToTrimFromApiKey, ""),
		SupportedHttpMethods:                       ExtractStringArray(d, MashEndpointSupportedHttpMethods, &EmptyStringArray),
		SystemDomainAuthentication:                 sem.systemDomainAuthenticationUpsertable(d),
		SystemDomains:                              sem.domainsArrayUpsertable(d, MashEndpointSystemDomains),
		TrafficManagerDomain:                       ExtractString(d, MashEndpointTrafficManagerDomain, ""),
		UseSystemDomainCredentials:                 ExtractBool(d, MashEndpointUseSystemDomainCredentials, false),
		SystemDomainCredentialKey:                  ExtractStringPointer(d, MashEndpointSystemDomainCredentialKey),
		SystemDomainCredentialSecret:               ExtractStringPointer(d, MashEndpointSystemDomainCredentialSecret),

		ParentServiceId: serviceIdent,
	}

	return rv, serviceIdent, rvd
}

func (sem *ServiceEndpointMapperImpl) persistCache(cache *masherytypes.Cache) []interface{} {
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

func (sem *ServiceEndpointMapperImpl) persistCors(cors *masherytypes.Cors) []interface{} {
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

func (sem *ServiceEndpointMapperImpl) persistProcessor(processor *masherytypes.Processor) []interface{} {
	// We will persist the processor in the state only if will contain useful information.
	// Mashery V3 API will return processor data structure also when the call transformation is not enabled.
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

		return rv
	} else {
		return nil
	}
}

func (sem *ServiceEndpointMapperImpl) persistDomains(d []masherytypes.Domain) []string {
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

func (sem *ServiceEndpointMapperImpl) persistAuthentication(d *masherytypes.SystemDomainAuthentication) []interface{} {
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

func (sem *ServiceEndpointMapperImpl) PersistTyped(endpoint masherytypes.Endpoint, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashEndpointId:                                         endpoint.Id,
		MashEndpointAllowMissingApiKey:                         endpoint.AllowMissingApiKey,
		MashEndpointApiKeyValueLocationKey:                     endpoint.ApiKeyValueLocationKey,
		MashEndpointApiKeyValueLocations:                       endpoint.ApiKeyValueLocations,
		MashEndpointApiMethodDetectionKey:                      endpoint.ApiMethodDetectionKey,
		MashEndpointApiMethodDetectionLocations:                endpoint.ApiMethodDetectionLocations,
		MashEndpointCache:                                      sem.persistCache(endpoint.Cache),
		MashEndpointConnectionTimeoutForSystemDomainRequest:    endpoint.ConnectionTimeoutForSystemDomainRequest,
		MashEndpointConnectionTimeoutForSystemDomainResponse:   endpoint.ConnectionTimeoutForSystemDomainResponse,
		MashEndpointCookiesDuringHttpRedirectsEnabled:          endpoint.CookiesDuringHttpRedirectsEnabled,
		MashEndpointCors:                                       sem.persistCors(endpoint.Cors),
		MashObjCreated:                                         endpoint.Created.ToString(),
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
		MashEndpointProcessor:                                  sem.persistProcessor(endpoint.Processor),
		MashEndpointPublicDomains:                              sem.persistDomains(endpoint.PublicDomains),
		MashEndpointRequestAuthenticationType:                  endpoint.RequestAuthenticationType,
		MashEndpointRequestPathAlias:                           endpoint.RequestPathAlias,
		MashEndpointRequestProtocol:                            endpoint.RequestProtocol,
		MashEndpointOauthGrantTypes:                            endpoint.OAuthGrantTypes,
		MashEndpointStringsToTrimFromApiKey:                    endpoint.StringsToTrimFromApiKey,
		MashEndpointSupportedHttpMethods:                       endpoint.SupportedHttpMethods,
		MashEndpointSystemDomainAuthentication:                 sem.persistAuthentication(endpoint.SystemDomainAuthentication),
		MashEndpointSystemDomains:                              sem.persistDomains(endpoint.SystemDomains),
		MashEndpointTrafficManagerDomain:                       endpoint.TrafficManagerDomain,
		MashObjUpdated:                                         endpoint.Updated.ToString(),
		MashEndpointUseSystemDomainCredentials:                 endpoint.UseSystemDomainCredentials,
		MashEndpointSystemDomainCredentialKey:                  endpoint.SystemDomainCredentialKey,
		MashEndpointSystemDomainCredentialSecret:               endpoint.SystemDomainCredentialSecret,
	}

	return sem.persistMap(endpoint.Identifier(), data, d)
}

// init
func init() {
	ServiceEndpointMapper = &ServiceEndpointMapperImpl{
		ResourceMapperImpl: ResourceMapperImpl{
			schema:       EndpointSchema,
			v3ObjectName: "endpoint",
			v3Identity: func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				rv := masherytypes.ServiceEndpointIdentifier{}
				rvd := diag.Diagnostics{}

				if !CompoundIdFrom(&rv, d.Id()) {
					rvd = append(rvd, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "endpoint identifier is incomplete",
					})
				}

				return rv, rvd
			},

			upsertFunc: func(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics) {
				return ServiceEndpointMapper.UpsertableTyped(d)
			},
			persistFunc: func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
				ptr := rv.(*masherytypes.Endpoint)
				return ServiceEndpointMapper.PersistTyped(*ptr, d)
			},
		},
	}
}
