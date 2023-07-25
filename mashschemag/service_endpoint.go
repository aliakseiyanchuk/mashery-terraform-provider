package mashschemag

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var ServiceEndpointResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.ServiceEndpointIdentifier]{
		IdentityFunc: func() masherytypes.ServiceEndpointIdentifier {
			return masherytypes.ServiceEndpointIdentifier{}
		},
	})

// Service endpoint parent identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.ServiceIdentifier]{
		Key: mashschema.MashSvcRef,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Service reference, to which this plan belongs",
		},
		IdentityFunc: func() masherytypes.ServiceIdentifier {
			return masherytypes.ServiceIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.ServiceIdentifier) bool {
			return len(inp.ServiceId) > 0
		},
	}

	ServiceEndpointResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

// Read-only fields
func init() {
	ServiceEndpointResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashServiceEndpointId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 identifier of this endpoint",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashObjCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashObjUpdated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	})
}

// Read-write fields
func init() {
	ServiceEndpointResourceSchemaBuilder.Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.AllowMissingApiKey
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointAllowMissingApiKey,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to true to allow calls without API keys",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.ApiKeyValueLocationKey
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointApiKeyValueLocationKey,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key in request parameters that identifies developer's api key",
			},
		},
	}).Add(&tfmapper.StringArrayFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *[]string {
			return &in.ApiKeyValueLocations
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointApiKeyValueLocations,
			Schema: &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				Description: "Locations where the developer should place key",
				// Probably would be worth-while adding defaults in the description
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				vals := mashschema.ExtractStringArray(in, key, &mashschema.EmptyStringArray)
				for _, val := range vals {
					if dg := mashschema.ValidateStringValueInSet(val, cty.GetAttrPath(key), &mashschema.DeveloperAPIKeyLocationsEnum); dg.HasError() {
						return false, dg[0].Detail
					}
				}

				// If
				if len(vals) == 2 {
					// Validate that only request-parameters and request-body are specified
					if mashschema.FindInArray(mashschema.RequestBodyVal, &vals) < 0 ||
						mashschema.FindInArray(mashschema.RequestParametersVal, &vals) < 0 {
						return false, fmt.Sprintf("%s and %s is an illegal combination, only %s and %s can be specified together",
							vals[0], vals[1],
							mashschema.RequestBodyVal, mashschema.RequestParametersVal,
						)
					}
				} else if len(vals) > 2 {
					return false, fmt.Sprintf("too many options specified (%d), maximum possible is 2", len(vals))
				}

				return true, ""
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.ApiMethodDetectionKey
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointApiMethodDetectionKey,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method v3Identity, aka the string that uniquely identifies the incoming method to the Traffic Manager. I",
			},
		},
	}).Add(&tfmapper.StringArrayFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *[]string {
			return &in.ApiMethodDetectionLocations
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointApiMethodDetectionLocations,
			Schema: &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Locations to derive api method from. Valid options are: request-header, request-body, request-parameters, request-path, and custom",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: mashschema.StringHashcode,
			},
			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				vals := mashschema.ExtractStringArray(in, key, &mashschema.EmptyStringArray)
				for _, val := range vals {
					if dg := mashschema.ValidateStringValueInSet(val, cty.GetAttrPath(key), &mashschema.MethodLocationsEnum); dg.HasError() {
						return false, dg[0].Detail
					}
				}

				// If
				if len(vals) == 2 {
					// Validate that only request-parameters and request-body are specified
					if mashschema.FindInArray(mashschema.RequestBodyVal, &vals) < 0 ||
						mashschema.FindInArray(mashschema.RequestParametersVal, &vals) < 0 {
						return false, fmt.Sprintf("%s and %s is an illegal combination, only %s and %s can be specified together",
							vals[0], vals[1],
							mashschema.RequestBodyVal, mashschema.RequestParametersVal,
						)
					}
				} else if len(vals) > 2 {
					return false, fmt.Sprintf("too many options specified (%d), maximum possible is 2", len(vals))
				}

				return true, ""
			},
		},
	}).Add(&tfmapper.PluggableFiledMapperBase[masherytypes.Endpoint]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointCache,
			Schema: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Set:      cacheHashCode,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mashschema.MashEndpointCacheClientSurrogateControlEnabled: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						mashschema.MashEndpointCacheContentCacheKeyHeaders: {
							Type:     schema.TypeSet,
							Required: true,
							Set:      mashschema.StringHashcode,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						mashschema.MashEndpointCacheTTLOverride: {
							Type:     schema.TypeInt,
							Optional: true,
						},
						mashschema.MashEndpointCacheIncludeApiKeyInContentCacheKey: {
							Type:     schema.TypeBool,
							Optional: true,
						},
						mashschema.MashEndpointCacheRespondFromStaleCacheEnabled: {
							Type:     schema.TypeBool,
							Optional: true,
						},
						mashschema.MashEndpointCacheResponseCacheControlEnabled: {
							Type:     schema.TypeBool,
							Optional: true,
						},
						mashschema.MashEndpointCacheVaryHeaderEnabled: {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
		RemoteToSchemaFunc: func(remote *masherytypes.Endpoint, key string, state *schema.ResourceData) *diag.Diagnostic {
			var v []interface{}

			if remote.Cache != nil && !remote.Cache.IsEmpty() {
				mp := map[string]interface{}{
					mashschema.MashEndpointCacheClientSurrogateControlEnabled:  remote.Cache.ClientSurrogateControlEnabled,
					mashschema.MashEndpointCacheContentCacheKeyHeaders:         remote.Cache.ContentCacheKeyHeaders,
					mashschema.MashEndpointCacheTTLOverride:                    int(remote.Cache.CacheTTLOverride),
					mashschema.MashEndpointCacheIncludeApiKeyInContentCacheKey: remote.Cache.IncludeApiKeyInContentCacheKey,
					mashschema.MashEndpointCacheRespondFromStaleCacheEnabled:   remote.Cache.RespondFromStaleCacheEnabled,
					mashschema.MashEndpointCacheResponseCacheControlEnabled:    remote.Cache.ResponseCacheControlEnabled,
					mashschema.MashEndpointCacheVaryHeaderEnabled:              remote.Cache.VaryHeaderEnabled,
				}
				v = append(v, mp)
			}

			return tfmapper.SetKeyWithDiag(state, key, v)
		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.Endpoint) {
			if cacheSet, exists := state.GetOk(key); exists {
				if remote.Cache == nil {
					remote.Cache = &masherytypes.Cache{}
				}

				tfCache := mashschema.UnwrapStructFromTerraformSet(cacheSet)

				schemaMapToCache(tfCache, remote.Cache)
			}
		},
	}).Add(&tfmapper.IntFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *int {
			return &in.ConnectionTimeoutForSystemDomainRequest
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointConnectionTimeoutForSystemDomainRequest,
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "Timeout to connect to the back-end",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateIntValueInSet(i, path, &mashschema.ConnectionTimeoutForSystemDomainRequestEnum)
				},
			},
		},
	}).Add(&tfmapper.IntFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *int {
			return &in.ConnectionTimeoutForSystemDomainResponse
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointConnectionTimeoutForSystemDomainResponse,
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Timeout to receive response from the back-end",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateIntValueInSet(i, path, &mashschema.ConnectionTimeoutForSystemDomainResponseEnum)
				},
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.CookiesDuringHttpRedirectsEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointCookiesDuringHttpRedirectsEnabled,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Set to true to enable cookies during redirects",
			},
		},
	}).Add(&tfmapper.PluggableFiledMapperBase[masherytypes.Endpoint]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointCors,
			Schema: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Set:      corsHashCode,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mashschema.MashEndpointCorsAllDomainsEnabled: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						mashschema.MashEndpointCorsSubDomainMatchingAllowed: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						mashschema.MashEndpointCorsCookiesAllowed: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						mashschema.MashEndpointCorsMaxAge: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  300,
						},
						mashschema.MashEndpointCorsAllowedDomains: {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      mashschema.StringHashcode,
							Elem:     mashschema.StringElem(),
						},
						mashschema.MashEndpointCorsAllowedHeaders: {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      mashschema.StringHashcode,
							Elem:     mashschema.StringElem(),
						},
						mashschema.MashEndpointCorsExposedHeaders: {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      mashschema.StringHashcode,
							Elem:     mashschema.StringElem(),
						},
					},
				},
			},
		},
		RemoteToSchemaFunc: func(remote *masherytypes.Endpoint, key string, state *schema.ResourceData) *diag.Diagnostic {
			var v []interface{}

			if remote.Cors != nil {
				mp := map[string]interface{}{
					mashschema.MashEndpointCorsAllDomainsEnabled:        remote.Cors.AllDomainsEnabled,
					mashschema.MashEndpointCorsSubDomainMatchingAllowed: remote.Cors.SubDomainMatchingAllowed,
					mashschema.MashEndpointCorsCookiesAllowed:           remote.Cors.CookiesAllowed,
					mashschema.MashEndpointCorsMaxAge:                   remote.Cors.MaxAge,
					mashschema.MashEndpointCorsExposedHeaders:           remote.Cors.ExposedHeaders,
				}

				if tfmapper.IsNonEmptyStringArray(&remote.Cors.AllowedDomains) {
					mp[mashschema.MashEndpointCorsAllowedDomains] = remote.Cors.AllowedDomains
				}
				if tfmapper.IsNonEmptyStringArray(&remote.Cors.AllowedHeaders) {
					mp[mashschema.MashEndpointCorsAllowedHeaders] = remote.Cors.AllowedHeaders
				}
				if tfmapper.IsNonEmptyStringArray(&remote.Cors.ExposedHeaders) {
					mp[mashschema.MashEndpointCorsExposedHeaders] = remote.Cors.ExposedHeaders
				}

				v = append(v, mp)
			}

			return tfmapper.SetKeyWithDiag(state, key, v)
		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.Endpoint) {
			if cacheSet, exists := state.GetOk(key); exists {
				if remote.Cors == nil {
					remote.Cors = &masherytypes.Cors{}
				}

				schemaMapToCORS(mashschema.UnwrapStructFromTerraformSet(cacheSet), remote.Cors)
			}
		},
	}).Add(&tfmapper.StringPtrFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) **string {
			return &in.CustomRequestAuthenticationAdapter
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointCustomRequestAuthenticationAdapter,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom adapter for authentication for an endpoint",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.DropApiKeyFromIncomingCall
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointDropApiKeyFromIncomingCall,
			Schema: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.ForceGzipOfBackendCall
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointForceGzipOfBackendCall,
			Schema: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.GzipPassthroughSupportEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointGzipPassthroughSupportEnabled,
			Schema: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}).Add(&tfmapper.StringArrayFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *[]string {
			return &in.HeadersToExcludeFromIncomingCall
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointHeadersToExcludeFromIncomingCall,
			Schema: &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "HTTP Headers to Drop from Incoming Call",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: mashschema.StringHashcode,
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.HighSecurity
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointHighSecurity,
			Schema: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.HostPassthroughIncludedInBackendCallHeader
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointHostPassthroughIncludedInBackendCallHeader,
			Schema: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.InboundSslRequired
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointInboundSslRequired,
			Schema: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.InboundMutualSslRequired
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointInboundMutualSslRequired,
			Schema: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.JsonpCallbackParameter
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointJsonpCallbackParameter,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parameter is used by Traffic Manager while handling the JSON responses.",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.JsonpCallbackParameterValue
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointJsonpCallbackParameterValue,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default parameter value can be set to be used by Traffic Manager to effectively handle the JSON response.",
			},
		},
	}).Add(&tfmapper.StringArrayFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *[]string {
			return &in.ForwardedHeaders
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointForwardedHeaders,
			Schema: &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Specific Mashery headers to be forwarded to the back-end",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: mashschema.StringHashcode,
			},
			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				vals := mashschema.ExtractStringArray(in, key, &mashschema.EmptyStringArray)
				for _, val := range vals {
					if dg := mashschema.ValidateStringValueInSet(val, cty.GetAttrPath(key), &mashschema.ForwardedHeadersEnum); dg.HasError() {
						return false, dg[0].Detail
					}
				}

				return true, ""
			},
		},
	}).Add(&tfmapper.StringArrayFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *[]string {
			return &in.ReturnedHeaders
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointReturnedHeaders,
			Schema: &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Specific Mashery headers to be forwarded to the back-end",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: mashschema.StringHashcode,
			},
			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				vals := mashschema.ExtractStringArray(in, key, &mashschema.EmptyStringArray)
				for _, val := range vals {
					if dg := mashschema.ValidateStringValueInSet(val, cty.GetAttrPath(key), &mashschema.ReturnedHeadersEnum); dg.HasError() {
						return false, dg[0].Detail
					}
				}

				return true, ""
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Endpoint name",
			},
		},
	}).Add(&tfmapper.IntFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *int {
			return &in.NumberOfHttpRedirectsToFollow
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointNumberOfHttpRedirectsToFollow,
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of HTTP redirects to follow",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.OutboundRequestTargetPath
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointOutboundRequestTargetPath,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/",
				Description: "Root context of the outbound domain to serve this data from",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.OutboundRequestTargetQueryParameters
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointOutboundRequestTargetQueryParameters,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query string to be added to the original request",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.OutboundTransportProtocol
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointOutboundTransportProtocol,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https",
				Description: "Outbound request protocol, defaults to https",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(i, path, &mashschema.OutboundTransportProtocolEnum)
				},
			},
		},
	}).Add(&tfmapper.PluggableFiledMapperBase[masherytypes.Endpoint]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointProcessor,
			Schema: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Set:      processorHashCode,
				Elem: &schema.Resource{
					Schema: mashschema.EndpointProcessorSchema,
				},
			},
			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				if cfgRAW, ok := in.GetOk(key); ok {
					cfg := mashschema.UnwrapStructFromTerraformSet(cfgRAW)
					processor := masherytypes.Processor{}

					schemaMapToProcessor(cfg, &processor)

					if len(processor.Adapter) == 0 {
						return false, "empty adapter string is pointless"
					} else if !processor.PreProcessEnabled && !processor.PostProcessEnabled {
						return false, "an adapter configuration that enables neither pre- nor post-processing is pointless"
					}
				}

				return true, ""
			},
		},
		RemoteToSchemaFunc: remoteProcessorToSchema,
		SchemaToRemoteFunc: schemaProcessorToRemote,
	}).Add(&tfmapper.PluggableFiledMapperBase[masherytypes.Endpoint]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointPublicDomains,
			Schema: &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: mashschema.StringHashcode,
			},
		},
		RemoteToSchemaFunc: func(remote *masherytypes.Endpoint, key string, state *schema.ResourceData) *diag.Diagnostic {
			v := make([]string, len(remote.PublicDomains))

			for i, d := range remote.PublicDomains {
				v[i] = d.Address
			}

			return tfmapper.SetKeyWithDiag(state, key, v)
		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.Endpoint) {
			tfDomans := mashschema.ExtractStringArray(state, key, &[]string{})

			remote.PublicDomains = make([]masherytypes.Domain, len(tfDomans))
			for i, v := range tfDomans {
				remote.PublicDomains[i] = masherytypes.Domain{
					Address: v,
				}
			}
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.RequestAuthenticationType
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointRequestAuthenticationType,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authentication type for the endpoint",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(i, path, &mashschema.RequestAuthenticationTypeEnum)
				},
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.RequestPathAlias
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointRequestPathAlias,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Request path",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.RequestProtocol
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointRequestProtocol,
			Schema: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "rest",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(i, path, &mashschema.RequestProtocolEnum)
				},
			},
		},
	}).Add(&tfmapper.StringArrayFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *[]string {
			return &in.OAuthGrantTypes
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointOauthGrantTypes,
			Schema: &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "OAuth grant types supported at this endpoint",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				vals := mashschema.ExtractStringArray(in, key, &mashschema.EmptyStringArray)
				for _, val := range vals {
					if dg := mashschema.ValidateStringValueInSet(val, cty.GetAttrPath(key), &mashschema.OAuthGrantTypesEnum); dg.HasError() {
						return false, dg[0].Detail
					}
				}

				return true, ""
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.StringsToTrimFromApiKey
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointStringsToTrimFromApiKey,
			Schema: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}).Add(&tfmapper.StringArrayFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *[]string {
			return &in.SupportedHttpMethods
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointSupportedHttpMethods,
			Schema: &schema.Schema{
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Methods this endpoint will support",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: mashschema.StringHashcode,
			},
			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				vals := mashschema.ExtractStringArray(in, key, &mashschema.EmptyStringArray)
				for _, val := range vals {
					if dg := mashschema.ValidateStringValueInSet(val, cty.GetAttrPath(key), &mashschema.HttoMethodsEnum); dg.HasError() {
						return false, dg[0].Detail
					}
				}

				return true, ""
			},
		},
	}).Add(&tfmapper.PluggableFiledMapperBase[masherytypes.Endpoint]{
		// This needs to be extracted into a separate mapper
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointSystemDomainAuthentication,
			Schema: &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Set:      systemDomainAuthHashCode,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mashschema.MashEndpointSystemDomainAuthenticationType: {
							Type:     schema.TypeString,
							Required: true,
						},
						mashschema.MashEndpointSystemDomainAuthenticationUsername: {
							Type:     schema.TypeString,
							Optional: true,
						},
						mashschema.MashEndpointSystemDomainAuthenticationCertificate: {
							Type:     schema.TypeString,
							Optional: true,
						},
						mashschema.MashEndpointSystemDomainAuthenticationPassword: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
		RemoteToSchemaFunc: remoteEndpointSystemAuthenticationToSchema,
		SchemaToRemoteFunc: schemaEndpointSystemAuthenticationToRemote,
	}).Add(&tfmapper.PluggableFiledMapperBase[masherytypes.Endpoint]{
		// This needs to be extracted into a separate mapper
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointSystemDomains,
			Schema: &schema.Schema{
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "The domain name of the client API server",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: mashschema.StringHashcode,
			},
		},
		RemoteToSchemaFunc: func(remote *masherytypes.Endpoint, key string, state *schema.ResourceData) *diag.Diagnostic {
			v := make([]string, len(remote.SystemDomains))

			for i, d := range remote.SystemDomains {
				v[i] = d.Address
			}

			return tfmapper.SetKeyWithDiag(state, key, v)
		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.Endpoint) {
			tfDomans := mashschema.ExtractStringArray(state, key, &[]string{})

			remote.SystemDomains = make([]masherytypes.Domain, len(tfDomans))
			for i, v := range tfDomans {
				remote.SystemDomains[i] = masherytypes.Domain{
					Address: v,
				}
			}
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.TrafficManagerDomain
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointTrafficManagerDomain,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Traffic Manager internal hostname (domain name) to which the requested public hostname is CNAMED.",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) *bool {
			return &in.UseSystemDomainCredentials
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointUseSystemDomainCredentials,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "To suit to the client server's requirement, Mashery can swap the API credentials such as API keys and send the swapped Mashery credentials to the client API server.",
				Optional:    true,
				Default:     false,
			},
		},
	}).Add(&tfmapper.StringPtrFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) **string {
			return &in.SystemDomainCredentialKey
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointSystemDomainCredentialKey,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key to use when making call to the client API server",
			},
		},
	}).Add(&tfmapper.StringPtrFieldMapper[masherytypes.Endpoint]{
		Locator: func(in *masherytypes.Endpoint) **string {
			return &in.SystemDomainCredentialSecret
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointSystemDomainCredentialSecret,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secret to use when making call to the client API server.",
			},
		},
	})
}

func schemaMapToCache(tfCache map[string]interface{}, v *masherytypes.Cache) {
	mashschema.ExtractKeyFromMap(tfCache, mashschema.MashEndpointCacheClientSurrogateControlEnabled, &v.ClientSurrogateControlEnabled)
	mashschema.ExtractKeyFromMap(tfCache, mashschema.MashEndpointCacheIncludeApiKeyInContentCacheKey, &v.IncludeApiKeyInContentCacheKey)
	mashschema.ExtractKeyFromMap(tfCache, mashschema.MashEndpointCacheRespondFromStaleCacheEnabled, &v.RespondFromStaleCacheEnabled)
	mashschema.ExtractKeyFromMap(tfCache, mashschema.MashEndpointCacheResponseCacheControlEnabled, &v.ResponseCacheControlEnabled)
	mashschema.ExtractKeyFromMap(tfCache, mashschema.MashEndpointCacheVaryHeaderEnabled, &v.VaryHeaderEnabled)

	var ttl int
	mashschema.ExtractKeyFromMap(tfCache, mashschema.MashEndpointCacheTTLOverride, &ttl)
	v.CacheTTLOverride = float64(ttl)

	mashschema.ExtractSchemaSetKeyFromMap(tfCache, mashschema.MashEndpointCacheContentCacheKeyHeaders, &v.ContentCacheKeyHeaders)
}

func schemaMapToCORS(tfCors map[string]interface{}, corsObject *masherytypes.Cors) {
	mashschema.ExtractKeyFromMap(tfCors, mashschema.MashEndpointCorsAllDomainsEnabled, &corsObject.AllDomainsEnabled)
	mashschema.ExtractKeyFromMap(tfCors, mashschema.MashEndpointCorsSubDomainMatchingAllowed, &corsObject.SubDomainMatchingAllowed)
	mashschema.ExtractKeyFromMap(tfCors, mashschema.MashEndpointCorsCookiesAllowed, &corsObject.CookiesAllowed)
	mashschema.ExtractKeyFromMap(tfCors, mashschema.MashEndpointCorsMaxAge, &corsObject.MaxAge)

	mashschema.ExtractSchemaSetKeyFromMap(tfCors, mashschema.MashEndpointCorsAllowedDomains, &corsObject.AllowedDomains)
	mashschema.ExtractSchemaSetKeyFromMap(tfCors, mashschema.MashEndpointCorsAllowedHeaders, &corsObject.AllowedHeaders)
	mashschema.ExtractSchemaSetKeyFromMap(tfCors, mashschema.MashEndpointCorsExposedHeaders, &corsObject.ExposedHeaders)
}

// Implementation function
func remoteProcessorToSchema(remote *masherytypes.Endpoint, key string, state *schema.ResourceData) *diag.Diagnostic {
	v := []interface{}{}

	if remote.Processor != nil && !remote.Processor.IsEmpty() {
		processorSchema := map[string]interface{}{
			mashschema.MashEndpointProcessorAdapter:            remote.Processor.Adapter,
			mashschema.MashEndpointProcessorPreProcessEnabled:  remote.Processor.PreProcessEnabled,
			mashschema.MashEndpointProcessorPostProcessEnabled: remote.Processor.PostProcessEnabled,
		}

		if len(remote.Processor.PreInputs) > 0 {
			processorSchema[mashschema.MashEndpointProcessorPreConfig] = remote.Processor.PreInputs
		}
		if len(remote.Processor.PostInputs) > 0 {
			processorSchema[mashschema.MashEndpointProcessorPostConfig] = remote.Processor.PostInputs
		}

		v = append(v, processorSchema)
	}

	return tfmapper.SetKeyWithDiag(state, key, v)
}

// Error set initialization
func init() {
	ServiceEndpointResourceSchemaBuilder.Add(&tfmapper.PluggableFiledMapperBase[masherytypes.Endpoint]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointErrorSetRef,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Error set applied to this endpoint",
			},

			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				errIdentStr := mashschema.ExtractString(in, key, "")
				if len(errIdentStr) == 0 {
					return true, ""
				}

				var endpIdent masherytypes.ServiceEndpointIdentifier
				var errSetIdent masherytypes.ErrorSetIdentifier

				// Parse and control the identifiers. We need to check that the error set identifier
				// belongs to the same service as this endpoint.
				_ = tfmapper.UnwrapJSON(in.Id(), &endpIdent)
				if err := tfmapper.UnwrapJSON(errIdentStr, &errSetIdent); err != nil {
					return false, fmt.Sprintf("error set identifier is malformed: %s", err.Error())
				}

				if endpIdent.ServiceId != errSetIdent.ServiceId {
					return false, fmt.Sprintf("error set identifier is from service %s, this endpoint is from %s", errSetIdent.ServiceId, endpIdent.ServiceId)
				}

				return true, ""
			},
		},
		RemoteToSchemaFunc: func(remote *masherytypes.Endpoint, key string, state *schema.ResourceData) *diag.Diagnostic {
			val := ""
			if remote.ErrorSet != nil {

				var endpIdent masherytypes.ServiceEndpointIdentifier
				_ = tfmapper.UnwrapJSON(state.Id(), &endpIdent)

				errorSetId := masherytypes.ErrorSetIdentifier{
					ErrorSetId:        remote.ErrorSet.Id,
					ServiceIdentifier: endpIdent.ServiceIdentifier,
				}
				val = tfmapper.WrapJSON(errorSetId)
			}

			return tfmapper.SetKeyWithDiag(state, key, val)
		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.Endpoint) {
			errIdentStr := mashschema.ExtractString(state, key, "")
			if len(errIdentStr) > 0 {
				var errSetIdent masherytypes.ErrorSetIdentifier
				_ = tfmapper.UnwrapJSON(errIdentStr, &errSetIdent)

				remote.ErrorSet = &masherytypes.AddressableV3Object{
					Id: errSetIdent.ErrorSetId,
				}
			}
		},
		NilRemoteToSchemaFunc: func(key string, state *schema.ResourceData) *diag.Diagnostic {
			return tfmapper.SetKeyWithDiag(state, key, "")
		},
	})

	ServiceEndpointResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointUserControlledErrorLocation,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User controlled error format location",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(i, path, &mashschema.UserControlledErrorFormatEnum)
				},
			},
		},
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.UserControlledErrorLocation
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Endpoint]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Endpoint]{
			Key: mashschema.MashEndpointUserControlledErrorLocationKey,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Description: "User controlled error format location key",
				Optional:    true,
			},
		},
		Locator: func(in *masherytypes.Endpoint) *string {
			return &in.UserControlledErrorLocationKey
		},
	})
}

func schemaProcessorToRemote(state *schema.ResourceData, key string, remote *masherytypes.Endpoint) {
	// Initialize the default processor
	if remote.Processor == nil {
		remote.Processor = &masherytypes.Processor{
			Adapter:            "",
			PreProcessEnabled:  false,
			PostProcessEnabled: false,
			PreInputs:          map[string]string{},
			PostInputs:         map[string]string{},
		}
	}

	if set, ok := state.GetOk(key); ok {
		tfProcMap := mashschema.UnwrapStructFromTerraformSet(set)
		if len(tfProcMap) > 1 {
			schemaMapToProcessor(tfProcMap, remote.Processor)
		}
	}
}

func schemaMapToProcessor(tfProcMap map[string]interface{}, v *masherytypes.Processor) {
	mashschema.ExtractKeyFromMap(tfProcMap, mashschema.MashEndpointProcessorPreProcessEnabled, &v.PreProcessEnabled)
	mashschema.ExtractKeyFromMap(tfProcMap, mashschema.MashEndpointProcessorPostProcessEnabled, &v.PostProcessEnabled)
	mashschema.ExtractKeyFromMap(tfProcMap, mashschema.MashEndpointProcessorAdapter, &v.Adapter)

	v.PreInputs = mashschema.SchemaMapToStringMap(tfProcMap[mashschema.MashEndpointProcessorPreConfig])
	v.PostInputs = mashschema.SchemaMapToStringMap(tfProcMap[mashschema.MashEndpointProcessorPostConfig])
}

// Implementation function
func remoteEndpointSystemAuthenticationToSchema(remote *masherytypes.Endpoint, key string, state *schema.ResourceData) *diag.Diagnostic {
	v := []interface{}{}

	if remote.SystemDomainAuthentication != nil {
		systemDomainAuthSchema := map[string]interface{}{
			mashschema.MashEndpointSystemDomainAuthenticationType:        remote.SystemDomainAuthentication.Type,
			mashschema.MashEndpointSystemDomainAuthenticationUsername:    remote.SystemDomainAuthentication.Username,
			mashschema.MashEndpointSystemDomainAuthenticationCertificate: remote.SystemDomainAuthentication.Certificate,
			mashschema.MashEndpointSystemDomainAuthenticationPassword:    remote.SystemDomainAuthentication.Password,
		}

		v = append(v, systemDomainAuthSchema)
	}

	return tfmapper.SetKeyWithDiag(state, key, v)
}

func schemaEndpointSystemAuthenticationToRemote(state *schema.ResourceData, key string, remote *masherytypes.Endpoint) {
	if set, ok := state.GetOk(key); ok {
		if remote.SystemDomainAuthentication == nil {
			remote.SystemDomainAuthentication = &masherytypes.SystemDomainAuthentication{}
		}
		v := remote.SystemDomainAuthentication

		tfSysDomain := mashschema.UnwrapStructFromTerraformSet(set)

		// TODO: Unsafe lookups
		schemaMapToSystemDomainAuthentication(tfSysDomain, v)
	}
}

func schemaMapToSystemDomainAuthentication(tfSysDomain map[string]interface{}, v *masherytypes.SystemDomainAuthentication) {
	mashschema.ExtractKeyFromMap(tfSysDomain, mashschema.MashEndpointSystemDomainAuthenticationType, &v.Type)
	v.Username = mashschema.SafeLookupStringPointer(&tfSysDomain, mashschema.MashEndpointSystemDomainAuthenticationUsername)
	v.Certificate = mashschema.SafeLookupStringPointer(&tfSysDomain, mashschema.MashEndpointSystemDomainAuthenticationCertificate)
	v.Password = mashschema.SafeLookupStringPointer(&tfSysDomain, mashschema.MashEndpointSystemDomainAuthenticationPassword)
}

func corsHashCode(cors interface{}) int {
	if corsSchemaMap, ok := cors.(map[string]interface{}); !ok {
		return 0
	} else {
		cors := masherytypes.Cors{}
		schemaMapToCORS(corsSchemaMap, &cors)

		hashStr := fmt.Sprintf("%t;%s;%s;%t;%s;%d;%t",
			cors.AllDomainsEnabled,
			mashschema.SortedStringOf(&cors.AllowedDomains),
			mashschema.SortedStringOf(&cors.AllowedHeaders),
			cors.CookiesAllowed,
			mashschema.SortedStringOf(&cors.ExposedHeaders),
			cors.MaxAge,
			cors.SubDomainMatchingAllowed,
		)

		return schema.HashString(hashStr)
	}
}

func processorHashCode(cors interface{}) int {
	if corsSchemaMap, ok := cors.(map[string]interface{}); !ok {
		return 0
	} else {
		proc := masherytypes.Processor{}
		schemaMapToProcessor(corsSchemaMap, &proc)

		hashStr := fmt.Sprintf("%s;%t;%s;%t;%s",
			proc.Adapter,
			proc.PreProcessEnabled,
			mashschema.SortedMapOf(&proc.PreInputs),
			proc.PostProcessEnabled,
			mashschema.SortedMapOf(&proc.PostInputs),
		)

		return schema.HashString(hashStr)
	}
}
func cacheHashCode(cache interface{}) int {
	if cacheSchemaMap, ok := cache.(map[string]interface{}); !ok {
		return 0
	} else {
		c := masherytypes.Cache{}
		schemaMapToCache(cacheSchemaMap, &c)

		hashStr := fmt.Sprintf("%d;%s;%t;%t;%t;%t",
			int64(c.CacheTTLOverride),
			mashschema.SortedStringOf(&c.ContentCacheKeyHeaders),
			c.VaryHeaderEnabled,
			c.ResponseCacheControlEnabled,
			c.RespondFromStaleCacheEnabled,
			c.IncludeApiKeyInContentCacheKey,
		)

		return schema.HashString(hashStr)
	}
}

func nilOrValue(str *string) string {
	if str == nil {
		return "--NULL--"
	} else {
		return *str
	}
}

func systemDomainAuthHashCode(cache interface{}) int {
	if cacheSchemaMap, ok := cache.(map[string]interface{}); !ok {
		return 0
	} else {
		auth := masherytypes.SystemDomainAuthentication{}
		schemaMapToSystemDomainAuthentication(cacheSchemaMap, &auth)

		hashStr := fmt.Sprintf("%s;%s;%s;%s",
			auth.Type,
			nilOrValue(auth.Username),
			nilOrValue(auth.Password),
			nilOrValue(auth.Certificate),
		)

		return schema.HashString(hashStr)
	}
}
