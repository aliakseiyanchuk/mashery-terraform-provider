---
subcategory: "mashery"
layout: "mashery"
page_title: "Mashery: mashery_service"
description: |-
Defines an endpoint within Mashery service
---

# Resource: `mashery_endpoint`

An endpoint is a fairly heavy-weight object representing [`/services/{serviceId}/endpoints/{endpointId}`](https://developer.mashery.com/docs/read/mashery_api/30/resources/services/endpoints)
Mashery V3 API resource.

## Example Usage

```hcl
resource "mashery_endpoint" "terraform-endpoint" {
  service_id = mashery_service.terraform-service.id
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
* `allow_missing_api_key`
* `api_key_value_location_key`
* `api_key_value_locations`
* `api_method_detection_key`
* `api_method_detection_locations`
* `cache` endpoint cache configuration block, providing the following key:
  * `client_surrogate_control_enabled`
  * `content_cache_key_headers`
  * `cache_ttl_override`
  * `include_api_key_in_content_cache_key`
  * `respond_from_stale_cache_enabled`
  * `response_cache_control_enabled`
  * `vary_header_enabled`
* `connection_timeout_for_system_domain_request`
* `connection_timeout_for_system_domain_response`
* `cookies_during_http_redirects_enabled`
* `cors` cross-origin resource sharing settings specification block, if Mashery is to include these in the response.
  The block has the following parameters:
  * `all_domains_enabled` whether all domains are enabled
  * `sub_domain_matching_allowed` whether subdomain matching is allowed
  * `cookies_allowed` whether cookies are allowed in processing CORS requests
  * `allowed_domains` a set of domains where CORS is allowed
  * `allowed_headers` a set of headers that are allowed in CORS requests
  * `exposed_headers` a set of headers that are exposed (returned) back to the calling browser
  * `max_age` CORS max age to be set
* `custom_request_authentication_adapter`
* `drop_api_key_from_incoming_call`
* `force_gzip_of_backend_call`
* `gzip_passthrough_support_enabled`
* `headers_to_exclude_from_incoming_call`
* `high_security`
* `host_passthrough_included_in_backend_call_header`
* `inbound_ssl_required`
* `inbound_mutual_ssl_required`
* `jsonp_callback_parameter`
* `jsonp_callback_parameter_value`
* `forwarded_headers`
* `returned_headers`
* `name`
* `number_of_http_redirects_to_follow`
* `outbound_request_target_path`
* `outbound_request_target_query_parameters`
* `outbound_transport_protocol`
* `processor` configures the endpoint processor
  * `adapter` adapter class
  * `pre_process_enabled` whether to enable pre-processors
  * `pre_config` pre-configuration key-value map as required by the processor
  * `post_process_enabled` whether to enable post-processors
  * `post_config` post-configuration key-value map as required by the processor
* `public_domains`
* `request_authentication_type`
* `request_path_alias`
* `request_protocol`
* `oauth_grant_types`
* `strings_to_trim_from_api_key`
* `supported_http_methods`
* `system_domain_authentication`
  * `type` type of the system domain authentication: `httpBasic` or `clientSslCert`
  * `username` username to login
  * `password` password
  * `certificate` certificate name to use for authentication
* `system_domains`
* `traffic_manager_domain`
* `use_system_domain_credentials`
* `system_domain_credential_key`
* `system_domain_credential_secret`

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