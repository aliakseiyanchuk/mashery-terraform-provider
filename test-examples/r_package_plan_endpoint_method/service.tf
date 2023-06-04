
resource "mashery_service" "srv" {
  name_prefix="tf-oauth-service"
}

resource "mashery_service_oauth" "srv" {
  service_id = mashery_service.srv.id
}

resource "mashery_unique_path" "endpointPrefix" {
  prefix = "terrafom"
}

resource "mashery_service_endpoint" "endp" {
  # An endpoint belongs to the service
  service_id = mashery_service.srv.id
  name = "service-endpoint-1"
  request_authentication_type = "apiKey"
  api_key_value_locations = ["request-header"]
  request_path_alias = "${mashery_unique_path.endpointPrefix.path}/my-debug-endpoint"
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

resource "mashery_endpoint_method" "meth_abc" {
  endpoint_id = mashery_service_endpoint.endp.id
  name = "do something good"
  sample_json = file("${path.module}/meth_abc.json")
}

resource "mashery_endpoint_method_filter" "abc_filter" {
  method_id = mashery_endpoint_method.meth_abc.id
  name = "abc filter"
  json_fields = "/a"
}