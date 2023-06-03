resource "mashery_unique_path" "svc-prefix" {
  prefix="tf-demos"
}

# Endpoint belongs to the service, so we have to identify one.
resource "mashery_service" "svc" {
  name_prefix="tf-demo-endpoint"
}

resource "mashery_service_endpoint" "svc-endpoint" {
  # An endpoint belongs to the service
  service_id = mashery_service.svc.id
  name = "service-endpoint-1"
  request_authentication_type = "apiKey"
  api_key_value_locations = ["request-header"]
  request_path_alias = "${mashery_unique_path.svc-prefix.path}/my-debug-endpoint-1"
  supported_http_methods = [ "get" ]
  system_domains = [ "171.21.35.46" ]
  public_domains = [ "171.22.40.90" ]

  traffic_manager_domain = var.traffic_manager_domain

  inbound_mutual_ssl_required = false
  inbound_ssl_required = false

  outbound_request_target_path = "/backend/api"
  outbound_request_target_query_parameters = "a=b"
  outbound_transport_protocol = "http"

}
