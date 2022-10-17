package mashschema_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"terraform-provider-mashery/mashschema"
	"testing"
	"time"
)

// BUG: not sending data on concurrent requests.

func TestV3EndpointToResourceDataWithEmptyCache(t *testing.T) {
	source := masherytypes.Endpoint{
		Cache: &masherytypes.Cache{
			ClientSurrogateControlEnabled: false,
			ContentCacheKeyHeaders:        []string{},
		},
	}

	if _, exists := storeEndpointAndGet(&source, mashschema.MashEndpointCache); exists {
		t.Errorf("Cache configuration should not be set of effectively empty configuration")
	}
}

func storeEndpointAndGet(endp *masherytypes.Endpoint, key string) (interface{}, bool) {
	d := mashschema.ServiceEndpointMapper.TestResourceData()
	mashschema.ServiceEndpointMapper.PersistTyped(*endp, d)

	return d.GetOk(key)
}

func TestV3EndpointToResourceDataWithEmptyProcessor(t *testing.T) {
	source := masherytypes.Endpoint{
		Processor: &masherytypes.Processor{
			PreProcessEnabled:  false,
			PostProcessEnabled: false,
			PostInputs:         []string{},
			PreInputs:          []string{},
			Adapter:            "",
		},
	}

	if _, exists := storeEndpointAndGet(&source, mashschema.MashEndpointProcessor); exists {
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

	source := masherytypes.Endpoint{
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

	d := mashschema.ServiceEndpointMapper.TestResourceData()
	d.SetId("serviceId::endpointId")

	diags := mashschema.ServiceEndpointMapper.PersistTyped(source, d)
	if len(diags) > 0 {
		t.Errorf("full conversion has encountered %d errors where none were expected", len(diags))
	}

	//reverse, diags := mashschema.MashEndpointUpsertable(d)
	if len(diags) > 0 {
		t.Errorf("Reverse conversion has encountered %d errors where none were expected", len(diags))
	}

	// Doing the assertion that the loaded data is the same
	//mashery.assertSameString(t, "Id", &source.Id, &reverse.Id)
	//mashery.assertSameString(t, "Name", &source.Name, &reverse.Name)
	//mashery.assertSameBool(t, "AllowMissingApiKey", &source.AllowMissingApiKey, &reverse.AllowMissingApiKey)
	//mashery.assertSameString(t, "ApiKeyValueLocationKey", &source.ApiKeyValueLocationKey, &reverse.ApiKeyValueLocationKey)
	//mashery.assertSameStringArray(t, "ApiKeyValueLocationKey", source.ApiKeyValueLocations, reverse.ApiKeyValueLocations)
	//
	//mashery.assertSameString(t, "ApiMethodDetectionKey", &source.ApiMethodDetectionKey, &reverse.ApiMethodDetectionKey)
	//mashery.assertSameStringArray(t, "ApiMethodDetectionLocations", source.ApiMethodDetectionLocations, reverse.ApiMethodDetectionLocations)
	//
	//mashery.assertDeepEqual(t, "Cache", source.Cache, reverse.Cache)
	//
	//mashery.assertSameString(t, "CustomRequestAuthenticationAdapter", source.CustomRequestAuthenticationAdapter, reverse.CustomRequestAuthenticationAdapter)
	//
	//mashery.assertSameBool(t, "DropApiKeyFromIncomingCall", &source.DropApiKeyFromIncomingCall, &reverse.DropApiKeyFromIncomingCall)
	//mashery.assertSameBool(t, "ForceGzipOfBackendCall", &source.ForceGzipOfBackendCall, &reverse.ForceGzipOfBackendCall)
	//mashery.assertSameBool(t, "GzipPassthroughSupportEnabled", &source.GzipPassthroughSupportEnabled, &reverse.GzipPassthroughSupportEnabled)
	//
	//mashery.assertSameSet(t, "HeadersToExcludeFromIncomingCall", &source.HeadersToExcludeFromIncomingCall, &reverse.HeadersToExcludeFromIncomingCall)
	//
	//mashery.assertSameBool(t, "HighSecurity", &source.HighSecurity, &reverse.HighSecurity)
	//mashery.assertSameBool(t, "HighSecurity", &source.HostPassthroughIncludedInBackendCallHeader, &reverse.HostPassthroughIncludedInBackendCallHeader)
	//mashery.assertSameBool(t, "HighSecurity", &source.InboundSslRequired, &reverse.InboundSslRequired)
	//
	//mashery.assertSameString(t, "JsonpCallbackParameter", &source.JsonpCallbackParameter, &reverse.JsonpCallbackParameter)
	//mashery.assertSameString(t, "JsonpCallbackParameterValue", &source.JsonpCallbackParameterValue, &reverse.JsonpCallbackParameterValue)
	//
	//mashery.assertDeepEqual(t, "ForwardedHeaders", source.ForwardedHeaders, reverse.ForwardedHeaders)
	//mashery.assertDeepEqual(t, "ReturnedHeaders", source.ReturnedHeaders, reverse.ReturnedHeaders)
	//
	//mashery.assertSameInt(t, "ReturnedHeaders", &source.NumberOfHttpRedirectsToFollow, &reverse.NumberOfHttpRedirectsToFollow)
	//
	//mashery.assertSameString(t, "OutboundRequestTargetPath", &source.OutboundRequestTargetPath, &reverse.OutboundRequestTargetPath)
	//mashery.assertSameString(t, "OutboundRequestTargetQueryParameters", &source.OutboundRequestTargetQueryParameters, &reverse.OutboundRequestTargetQueryParameters)
	//mashery.assertSameString(t, "OutboundTransportProtocol", &source.OutboundTransportProtocol, &reverse.OutboundTransportProtocol)
	//
	//mashery.assertDeepEqual(t, "Processor", source.Processor, reverse.Processor)
	//mashery.assertDeepEqual(t, "PublicDomains", source.PublicDomains, reverse.PublicDomains)
	//
	//mashery.assertSameString(t, "RequestAuthenticationType", &source.RequestAuthenticationType, &reverse.RequestAuthenticationType)
	//mashery.assertSameString(t, "RequestPathAlias", &source.RequestPathAlias, &reverse.RequestPathAlias)
	//mashery.assertSameString(t, "RequestProtocol", &source.RequestProtocol, &reverse.RequestProtocol)
	//
	//mashery.assertDeepEqual(t, "OAuthGrantTypes", source.OAuthGrantTypes, reverse.OAuthGrantTypes)
	//
	//mashery.assertSameString(t, "StringsToTrimFromApiKey", &source.StringsToTrimFromApiKey, &reverse.StringsToTrimFromApiKey)
	//
	//mashery.assertDeepEqual(t, "SupportedHttpMethods", source.SupportedHttpMethods, reverse.SupportedHttpMethods)
	//mashery.assertDeepEqual(t, "SystemDomainAuthentication", source.SystemDomainAuthentication, reverse.SystemDomainAuthentication)
	//mashery.assertDeepEqual(t, "SystemDomains", source.SystemDomains, reverse.SystemDomains)
	//
	//mashery.assertSameString(t, "TrafficManagerDomain", &source.TrafficManagerDomain, &reverse.TrafficManagerDomain)
	//mashery.assertSameBool(t, "UseSystemDomainCredentials", &source.UseSystemDomainCredentials, &reverse.UseSystemDomainCredentials)
	//
	//mashery.assertSameString(t, "SystemDomainCredentialKey", source.SystemDomainCredentialKey, reverse.SystemDomainCredentialKey)
	//mashery.assertSameString(t, "SystemDomainCredentialSecret", source.SystemDomainCredentialSecret, reverse.SystemDomainCredentialSecret)
}
