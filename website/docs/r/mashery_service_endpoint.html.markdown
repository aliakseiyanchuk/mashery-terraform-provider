---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: -|
  Defines an endpoint within Mashery service
---

# Resource: `mashery_endpoint`

An endpoint is a fairly heavy-weight object representing [`/services/{serviceId}/endpoints/{endpointId}`](https://developer.mashery.com/docs/read/mashery_api/30/resources/services/endpoints)
Mashery V3 API resource.

## Example Usage

```hcl
resource "mashery_service_endpoint" "terraform-endpoint" {
  service_ref = mashery_service.terraform-service.id
  name ="Terraform Endpoint"

  public_domains = ["api.public.domain.com"]
  request_path_alias = "/path"
  request_authentication_type = "apiKey"
  oauth_grant_types = ["implicit"]
  supported_http_methods = ["get"]

  traffic_manager_domain = "api.public.domain.com"

  api_key_value_locations = ["request-header"]
  headers_to_exclude_from_incoming_call = ["Authorization"]
  forwarded_headers = ["mashery-host", "mashery-message-id", "mashery-service-id"]

  // Where it's gonna be forwarded.
  system_domains = ["api.backend.com"]
  outbound_request_target_path = "/"
}
```

## Argument Reference

* `service_ref` service identifier, to which this endpoint belongs
* `allow_missing_api_key` set `true` to allow calls without API keys. This setting is for the calls made to API Management.
* `developer_api_key_field_name` field name of the developer's API key
* `developer_api_key_locations` one or more options where the Traffic Manager must look for a developer's API key. Only Parameters and Request Body can be chosen at the same time.
   Possible options are given below. Only `request-parameters` and `request-body` can be specified together
  * `request-header`
  * `request-body`
  * `request-parameters`
  * `request-path`
  * `custom`
* `api_method_identifier` method identifier, e.g. `2 4` 
* `api_method_locations` Select one or more options where the Traffic Manager must look for a method name.
   Possible options are given below. Only `request-parameters` and `request-body` can be specified together
  * `request-header`
  * `request-body`
  * `request-parameters`
  * `request-path`
  * `custom`
* `cache` endpoint cache configuration block, providing the following key:
  * `client_surrogate_control_enabled`
  * `content_cache_key_headers`
  * `cache_ttl_override`
  * `include_api_key_in_content_cache_key`
  * `respond_from_stale_cache_enabled`
  * `response_cache_control_enabled`
  * `vary_header_enabled`
* `connection_timeout_for_system_domain_request` timeout with possible options 2, 3, 4, 5, 6, and 7.
* `connection_timeout_for_system_domain_response` timeout for waiting for the response. Possible values 2, 5, 10, 20, 30, 45, 60, 120, 300, 600, 900, and 1200.
  A timeout setting above 60 seconds will require intervention from TIBCO via a Support Case.
* `enable_cookies_during_redirects` set `true` to enable cookies during redirects.
* `cors` cross-origin resource sharing settings specification block, if Mashery is to include these in the response.
  The block has the following parameters:
  * `all_domains_enabled` whether all domains are enabled
  * `sub_domain_matching_allowed` whether subdomain matching is allowed
  * `cookies_allowed` whether cookies are allowed in processing CORS requests
  * `allowed_domains` a set of domains where CORS is allowed
  * `allowed_headers` a set of headers that are allowed in CORS requests
  * `exposed_headers` a set of headers that are exposed (returned) back to the calling browser
  * `max_age` CORS max age to be set
* `custom_request_authentication_adapter` custom request authentication adapter class
* `drop_api_key_from_incoming_call` removes API key and signature query parameters from backend call
* `force_gzip_of_backend_call` set `true` to force GZip directive on the backend server. The Accept-encoding value passed by the client server is overridden.
* `gzip_passthrough_support_enabled` set `true to enable GZIP pass-through.
* `headers_to_exclude_from_incoming_call` a list of headers to drop (such as `Authorization`) 
* `high_security`: whether this endpoint is high security
* `host_passthrough_included_in_backend_call_header` include `X-Mashery-Host` header in the back-end call
* `inbound_ssl_required` requires HTTPS
* `inbound_mutual_ssl_required` require mutual HTTPS
* `jsonp_callback_parameter` JSON callback function parameter.
* `jsonp_callback_parameter_value` default callback function parameter value.
* `forwarded_headers`: forwarded headers toward API back-ends: Possible values are:
  * `mashery-host`
  * `mashery-message-id`
  * `mashery-service-id`
* `returned_headers`: returned headers:
  * `mashery-message-id`
  * `mashery-responder`
* `name`: endpoint name
* `number_of_http_redirects_to_follow`: number of redirects to follow
* `outbound_request_target_path` output path for the back-end call
* `outbound_request_target_query_parameters` outbound query parameters to be added to the call
* `outbound_transport_protocol` outbound protocol. Possible options:
  * `use-inbound` 
  * `http`
  * `https`
* `processor` configures the endpoint processor
  * `adapter` adapter class
  * `pre_process_enabled` whether to enable pre-processors
  * `pre_config` pre-configuration key-value map as required by the processor
  * `post_process_enabled` whether to enable post-processors
  * `post_config` post-configuration key-value map as required by the processor
* `public_domains` list of public domains
* `request_authentication_type` request authentication type, possible options:
  * `apiKey`
  * `apiKeyAndSecret_MD5`
  * `apiKeyAndSecret_SHA256`
  * `secureHash_SHA256`
  * `oauth`
  * `custom`
* `request_path_alias` client facing request path
* `request_protocol` request protocol of this endpoint, possible options:
  * `rest`
  * `soap`
  * `xml-rpc`
  * `json-rpc`
  * `other`
* `oauth_grant_types enabled for this endpoint:
  * `authorization-code`
  * `implicit`
  * `password`
  * `client-credentials`
* `strings_to_trim_from_api_key` strings to trio from API key
* `supported_http_methods` supported HTTP methods on this endpoint
* `system_domain_authentication`: object defining how Mashery should authenticate towards the backend:
  * `type` type of the system domain authentication: `httpBasic` or `clientSslCert`
  * `username` username to login
  * `password` password
  * `certificate` certificate name to use for authentication
* `system_domains` list of system domains to which the authenticated call should be forwarded
* `traffic_manager_domain` traffic manager domain
* `use_system_domain_credentials` To suit to the client server's requirement, Mashery can swap the API credentials such as API keys and send the swapped Mashery credentials to the client API server.
* `system_domain_credential_key`  Key to use when making call to the client API server.
* `system_domain_credential_secret` Secret to use when making call to the client API server.

## CORS Configuration Example

### Minimal enablement

A minimal configuration enables CORS responses for all domains allowing the browser to cache the pre-flight
responses for the given number of minutes. 

```terraform
resource "mashery_service_endpoint" "cors_endpoint" {
  # Configuration as desired

  cors {
    all_domains_enabled = true
    max_age = 30
  }
}
```

## Processor Configuration Example

Processor configuration requires adapter name and indication of call transformations to include (pre-processor and/or
post-processor). The plugin will block pointless configurations (which doesn't define adapter or defines to 
call transformation points). The pre- and post-configuration is optional and depends on the requirements of the
processor being configured.

```terraform
resource "mashery_service_endpoint" "endp" {
  # Configuration as desired

  processor {
    adapter = "com.github.fake-demo-processor"
    pre_process_enabled             = true
    pre_config     = {
      "a": "b",
      "c": "d",
      "e": "FFFddd"
    }
  }
}


```


## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `endpoint_id` Mashery endpoint endpoint Id
* `created` date endpoint was created
* `updated` date endpoint was updated