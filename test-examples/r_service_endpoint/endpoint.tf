# Endpoint belongs to the service, so we have to identify one.
resource "mashery_service" "srv" {
  name_prefix = "tf-demo-endpoint"
}

resource "mashery_service_error_set" "custom" {
  service_ref = mashery_service.srv.id
  name        = "Custom error set 1"
  type        = "application/ld+json"

  error_message {
    error         = "ERR_401_INVALID_NONCE"
    status        = "Really not available"
    detail_header = "our backend is basically awol"
    response_body = "{message: backend is not available, please stop hammering it}"
  }

}

resource "mashery_service_endpoint" "endp" {
  # An endpoint belongs to the service
  service_ref                 = mashery_service.srv.id
  name                        = "service-endpoint-1"
  request_authentication_type = "apiKey"
  developer_api_key_locations = ["request-header"]
  request_path_alias          = "/my-debug-endpoint"
  supported_http_methods      = ["get"]
  system_domains              = ["171.21.35.46"]
  public_domains              = ["171.22.40.90"]

  traffic_manager_domain = var.traffic_manager_domain

  inbound_mutual_ssl_required = false
  inbound_ssl_required        = false

  outbound_request_target_path             = "/backend/api"
  outbound_request_target_query_parameters = "a=b"
  outbound_transport_protocol              = "http"

}

