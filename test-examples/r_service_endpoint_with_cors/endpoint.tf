# Endpoint belongs to the service, so we have to identify one.
resource "mashery_service" "srv" {
  name_prefix = "tf-cors-endpoint"
}

resource "mashery_service_endpoint" "endp" {
  # An endpoint belongs to the service
  service_ref                 = mashery_service.srv.id
  name                        = "service-endpoint-1"
  request_authentication_type = "apiKey"
  api_key_value_locations     = ["request-header"]
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

  cors {
    all_domains_enabled = true
    max_age             = 30
    exposed_headers     = toset(["ABC", "DEF"])
  }
}

resource "mashery_service_endpoint" "endp-finegrained-cors" {
  # An endpoint belongs to the service
  service_ref = mashery_service.srv.id
  name = "service-endpoint-2"
  request_authentication_type = "apiKey"
  api_key_value_locations = ["request-header"]
  request_path_alias = "/my-debug-endpoint-cors2"
  supported_http_methods = [ "get" ]
  system_domains = [ "171.21.35.46" ]
  public_domains = [ "171.22.40.90" ]

  traffic_manager_domain = var.traffic_manager_domain

  inbound_mutual_ssl_required = false
  inbound_ssl_required = false

  outbound_request_target_path = "/backend/api"
  outbound_request_target_query_parameters = "a=b"
  outbound_transport_protocol = "http"

  cors {
    all_domains_enabled = false
    sub_domain_matching_allowed = true
    cookies_allowed = false
    max_age = 30
    allowed_domains = toset(["171.21.35.46", "171.21.35.47", "171.21.35.48"])
    allowed_headers = toset(["X-Ping-1", "X-Ping-2"])
    exposed_headers = toset(["X-App-1", "X-App-2"])
  }
}

