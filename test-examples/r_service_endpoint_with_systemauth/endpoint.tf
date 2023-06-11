# Endpoint belongs to the service, so we have to identify one.
resource "mashery_service" "srv" {
  name_prefix = "tf-systemauth-endpoint"
}

resource "mashery_service_endpoint" "endp" {
  # An endpoint belongs to the service
  service_ref                 = mashery_service.srv.id
  name                        = "service-endpoint-1"
  request_authentication_type = "apiKey"
  api_key_value_locations     = ["request-header"]
  request_path_alias          = "/my-demo-systemauth-processor"
  supported_http_methods      = ["get"]
  system_domains              = ["171.21.35.46"]
  public_domains              = ["171.22.40.90"]

  traffic_manager_domain = var.traffic_manager_domain

  inbound_mutual_ssl_required = false
  inbound_ssl_required        = false

  outbound_request_target_path             = "/backend/api"
  outbound_request_target_query_parameters = "a=b"
  outbound_transport_protocol              = "http"

  system_domain_authentication {
    type = "httpBasic"
    username = "adsfasdf"
    password = "asdfasdfasdf"
  }
}

resource "mashery_service_endpoint" "endp2" {
  # An endpoint belongs to the service
  service_ref                 = mashery_service.srv.id
  name                        = "service-endpoint-2"
  request_authentication_type = "apiKey"
  api_key_value_locations     = ["request-header"]
  request_path_alias          = "/my-demo-systemauth-endpoint/with-cert"
  supported_http_methods      = ["get"]
  system_domains              = ["171.21.35.46"]
  public_domains              = ["171.22.40.90"]

  traffic_manager_domain = var.traffic_manager_domain

  inbound_mutual_ssl_required = false
  inbound_ssl_required        = false

  outbound_request_target_path             = "/backend/api"
  outbound_request_target_query_parameters = "a=b"
  outbound_transport_protocol              = "http"

  system_domain_authentication {
    type = "clientSslCert"
    certificate = "abc"
    password = "def"
  }
}

