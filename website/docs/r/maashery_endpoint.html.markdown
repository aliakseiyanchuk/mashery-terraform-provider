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

* `service_id` service identifier, to which this endpoint belongs
* `allow_missing_api_key`
* `api_key_value_location_key`
* `api_key_value_locations`
* `api_method_detection_key`
* `api_method_detection_locations`
* `cache`
* `client_surrogate_control_enabled`
* `content_cache_key_headers`
* `connection_timeout_for_system_domain_request`
* `connection_timeout_for_system_domain_response`
* `cookies_during_http_redirects_enabled`
* `cors`
* `all_domains_enabled`
* `max_age`
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
* `processor`
* `adapter`
* `pre_process_enabled`
* `post_process_enabled`
* `pre_config`
* `post_config`
* `public_domains`
* `request_authentication_type`
* `request_path_alias`
* `request_protocol`
* `oauth_grant_types`
* `strings_to_trim_from_api_key`
* `supported_http_methods`
* `system_domain_authentication`
* `system_domains`
* `traffic_manager_domain`
* `use_system_domain_credentials`
* `system_domain_credential_key`
  * `system_domain_credential_secret`

## Attribute Reference

In addition to all arguments above, the following attributes are exposed:

* `endpoint_id` endpoint Id
* `created` date endpoint was created
* `updated` date endpoint was updated