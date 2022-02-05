package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashery"
	"testing"
	"time"
)

// BUG: not sending data on concurrent requests.

func TestV3EndpointToResourceDataWithEmptyCache(t *testing.T) {
	source := masherytypes.MasheryEndpoint{
		Cache: &masherytypes.Cache{
			ClientSurrogateControlEnabled: false,
			ContentCacheKeyHeaders:        []string{},
		},
	}

	if _, exists := storeEndpointAndGet(&source, mashery.MashEndpointCache); exists {
		t.Errorf("Cache configuration should not be set of effectively empty configuration")
	}
}

func storeEndpointAndGet(endp *masherytypes.MasheryEndpoint, key string) (interface{}, bool) {
	res := schema.Resource{
		Schema: mashery.EndpointSchema,
	}

	d := res.TestResourceData()
	mashery.V3EndpointToResourceData(endp, d)

	return d.GetOk(key)
}

func TestV3EndpointToResourceDataWithEmptyProcessor(t *testing.T) {
	source := masherytypes.MasheryEndpoint{
		Processor: &masherytypes.Processor{
			PreProcessEnabled:  false,
			PostProcessEnabled: false,
			PostInputs:         []string{},
			PreInputs:          []string{},
			Adapter:            "",
		},
	}

	if _, exists := storeEndpointAndGet(&source, mashery.MashEndpointProcessor); exists {
		t.Errorf("Processor configuration should not be set of effectively empty configuration")
	}
}

// Test full save/restore to/from the Terraform state.
func TestV3EndpointToResourceData(t *testing.T) {
	tm := masherytypes.MasheryJSONTime(time.Now())

	custAdapter := "custom-adapter"
	user := "user"
	cert := "cert"
	pwd := "pwd"

	sysCredsKey := "sysCredsKey"
	sysCredsPass := "sysCredsPass"

	source := masherytypes.MasheryEndpoint{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:      "endpointId",
			Name:    "endpName",
			Created: &tm,
			Updated: &tm,
		},
		AllowMissingApiKey:          true,
		ApiKeyValueLocationKey:      "api_key",
		ApiKeyValueLocations:        []string{"apiKeyValueLocs"},
		ApiMethodDetectionKey:       "1 2 3",
		ApiMethodDetectionLocations: []string{"methDetectLocs"},
		Cache: &masherytypes.Cache{
			ClientSurrogateControlEnabled: true,
			ContentCacheKeyHeaders:        []string{"header-a", "header-b"},
		},
		ConnectionTimeoutForSystemDomainRequest:  5,
		ConnectionTimeoutForSystemDomainResponse: 6,
		CookiesDuringHttpRedirectsEnabled:        true,
		Cors: &masherytypes.Cors{
			AllDomainsEnabled: true,
			MaxAge:            40,
		},
		CustomRequestAuthenticationAdapter:         &custAdapter,
		DropApiKeyFromIncomingCall:                 true,
		ForceGzipOfBackendCall:                     true,
		GzipPassthroughSupportEnabled:              true,
		HeadersToExcludeFromIncomingCall:           []string{"excl-a", "excl-b"},
		HighSecurity:                               true,
		HostPassthroughIncludedInBackendCallHeader: true,
		InboundSslRequired:                         true,
		JsonpCallbackParameter:                     "jsonp",
		JsonpCallbackParameterValue:                "jsonp-value",
		ScheduledMaintenanceEvent:                  nil,
		ForwardedHeaders:                           []string{"fh1"},
		ReturnedHeaders:                            []string{"rh1"},
		Methods:                                    nil,
		NumberOfHttpRedirectsToFollow:              3,
		OutboundRequestTargetPath:                  "/a",
		OutboundRequestTargetQueryParameters:       "?ef",
		OutboundTransportProtocol:                  "https",
		Processor: &masherytypes.Processor{
			PreProcessEnabled:  true,
			PostProcessEnabled: true,
			PostInputs: []string{
				"a", "b",
			},
			PreInputs: []string{
				"c", "d",
			},
			Adapter: "adapter",
		},
		PublicDomains:             []masherytypes.Domain{{Address: "addr"}},
		RequestAuthenticationType: "req-auth",
		RequestPathAlias:          "req-path-alias",
		RequestProtocol:           "req-proto",
		OAuthGrantTypes:           []string{"oauth-gt"},
		StringsToTrimFromApiKey:   "bba",
		SupportedHttpMethods:      []string{"support-methds"},
		SystemDomainAuthentication: &masherytypes.SystemDomainAuthentication{
			Type:        "ff",
			Username:    &user,
			Certificate: &cert,
			Password:    &pwd,
		},
		SystemDomains:                []masherytypes.Domain{{Address: "dom-addr"}},
		TrafficManagerDomain:         "tm-doamin",
		UseSystemDomainCredentials:   true,
		SystemDomainCredentialKey:    &sysCredsKey,
		SystemDomainCredentialSecret: &sysCredsPass,
	}

	res := schema.Resource{
		Schema: mashery.EndpointSchema,
	}

	d := res.TestResourceData()
	d.SetId("serviceId::endpointId")

	diags := mashery.V3EndpointToResourceData(&source, d)
	if len(diags) > 0 {
		t.Errorf("full conversion has encountered %d errors where none were expected", len(diags))
	}

	reverse, diags := mashery.MashEndpointUpsertable(d)
	if len(diags) > 0 {
		t.Errorf("Reverse conversion has encountered %d errors where none were expected", len(diags))
	}

	// Doing the assertion that the loaded data is the same
	assertSameString(t, "Id", &source.Id, &reverse.Id)
	assertSameString(t, "Name", &source.Name, &reverse.Name)
	assertSameBool(t, "AllowMissingApiKey", &source.AllowMissingApiKey, &reverse.AllowMissingApiKey)
	assertSameString(t, "ApiKeyValueLocationKey", &source.ApiKeyValueLocationKey, &reverse.ApiKeyValueLocationKey)
	assertSameStringArray(t, "ApiKeyValueLocationKey", source.ApiKeyValueLocations, reverse.ApiKeyValueLocations)

	assertSameString(t, "ApiMethodDetectionKey", &source.ApiMethodDetectionKey, &reverse.ApiMethodDetectionKey)
	assertSameStringArray(t, "ApiMethodDetectionLocations", source.ApiMethodDetectionLocations, reverse.ApiMethodDetectionLocations)

	assertDeepEqual(t, "Cache", source.Cache, reverse.Cache)

	assertSameString(t, "CustomRequestAuthenticationAdapter", source.CustomRequestAuthenticationAdapter, reverse.CustomRequestAuthenticationAdapter)

	assertSameBool(t, "DropApiKeyFromIncomingCall", &source.DropApiKeyFromIncomingCall, &reverse.DropApiKeyFromIncomingCall)
	assertSameBool(t, "ForceGzipOfBackendCall", &source.ForceGzipOfBackendCall, &reverse.ForceGzipOfBackendCall)
	assertSameBool(t, "GzipPassthroughSupportEnabled", &source.GzipPassthroughSupportEnabled, &reverse.GzipPassthroughSupportEnabled)

	assertSameSet(t, "HeadersToExcludeFromIncomingCall", &source.HeadersToExcludeFromIncomingCall, &reverse.HeadersToExcludeFromIncomingCall)

	assertSameBool(t, "HighSecurity", &source.HighSecurity, &reverse.HighSecurity)
	assertSameBool(t, "HighSecurity", &source.HostPassthroughIncludedInBackendCallHeader, &reverse.HostPassthroughIncludedInBackendCallHeader)
	assertSameBool(t, "HighSecurity", &source.InboundSslRequired, &reverse.InboundSslRequired)

	assertSameString(t, "JsonpCallbackParameter", &source.JsonpCallbackParameter, &reverse.JsonpCallbackParameter)
	assertSameString(t, "JsonpCallbackParameterValue", &source.JsonpCallbackParameterValue, &reverse.JsonpCallbackParameterValue)

	assertDeepEqual(t, "ForwardedHeaders", source.ForwardedHeaders, reverse.ForwardedHeaders)
	assertDeepEqual(t, "ReturnedHeaders", source.ReturnedHeaders, reverse.ReturnedHeaders)

	assertSameInt(t, "ReturnedHeaders", &source.NumberOfHttpRedirectsToFollow, &reverse.NumberOfHttpRedirectsToFollow)

	assertSameString(t, "OutboundRequestTargetPath", &source.OutboundRequestTargetPath, &reverse.OutboundRequestTargetPath)
	assertSameString(t, "OutboundRequestTargetQueryParameters", &source.OutboundRequestTargetQueryParameters, &reverse.OutboundRequestTargetQueryParameters)
	assertSameString(t, "OutboundTransportProtocol", &source.OutboundTransportProtocol, &reverse.OutboundTransportProtocol)

	assertDeepEqual(t, "Processor", source.Processor, reverse.Processor)
	assertDeepEqual(t, "PublicDomains", source.PublicDomains, reverse.PublicDomains)

	assertSameString(t, "RequestAuthenticationType", &source.RequestAuthenticationType, &reverse.RequestAuthenticationType)
	assertSameString(t, "RequestPathAlias", &source.RequestPathAlias, &reverse.RequestPathAlias)
	assertSameString(t, "RequestProtocol", &source.RequestProtocol, &reverse.RequestProtocol)

	assertDeepEqual(t, "OAuthGrantTypes", source.OAuthGrantTypes, reverse.OAuthGrantTypes)

	assertSameString(t, "StringsToTrimFromApiKey", &source.StringsToTrimFromApiKey, &reverse.StringsToTrimFromApiKey)

	assertDeepEqual(t, "SupportedHttpMethods", source.SupportedHttpMethods, reverse.SupportedHttpMethods)
	assertDeepEqual(t, "SystemDomainAuthentication", source.SystemDomainAuthentication, reverse.SystemDomainAuthentication)
	assertDeepEqual(t, "SystemDomains", source.SystemDomains, reverse.SystemDomains)

	assertSameString(t, "TrafficManagerDomain", &source.TrafficManagerDomain, &reverse.TrafficManagerDomain)
	assertSameBool(t, "UseSystemDomainCredentials", &source.UseSystemDomainCredentials, &reverse.UseSystemDomainCredentials)

	assertSameString(t, "SystemDomainCredentialKey", source.SystemDomainCredentialKey, reverse.SystemDomainCredentialKey)
	assertSameString(t, "SystemDomainCredentialSecret", source.SystemDomainCredentialSecret, reverse.SystemDomainCredentialSecret)
}
