resource "mashery_service" "srv" {
  name_prefix="tf-debug"
}

locals {
  package_grant_types = toset(["authorization_code", "password", "client_credentials"])
}

resource "mashery_service_oauth" "svc_oauth" {
  service_ref = mashery_service.srv.id

  access_token_ttl_enabled        = true
  access_token_ttl                = "1h"
  access_token_type               = "bearer"
  allow_multiple_token            = true
  authorization_code_ttl          = "5m"
  forwarded_headers               = toset(["access-token", "client-id", "scope", "user-context"])
  mashery_token_api_enabled       = false
  enable_refresh_token_ttl        = true
  force_oauth_redirect_url        = false
  force_ssl_redirect_url_enabled  = false
  grant_types                     = local.package_grant_types
  secure_tokens_enabled           = false
}


resource "mashery_service_endpoint" "oauth-endpoint" {
  # An endpoint belongs to the service
  service_ref                  = mashery_service.srv.id
  name                         = "oauth-service-loadtest"
  developer_api_key_field_name = "api_key"
  request_authentication_type  = "oauth"
  request_path_alias           = "/tech/opsmw-loadtest/oauth"
  supported_http_methods       = ["get", "post"]
  system_domains               = "171.22.36.47"
  public_domains               = ["171.21.35.46"]

  traffic_manager_domain = var.traffic_manager_domain

  inbound_mutual_ssl_required = false
  inbound_ssl_required        = true
  high_security               = false

  outbound_request_target_path             = "/api/"
  outbound_request_target_query_parameters = "a=b"
  outbound_transport_protocol              = "http"

  oauth_grant_types                     = local.package_grant_types
}

resource "mashery_package" "oauth" {
  name        = "OAuth Package"
  description = "OAuth package managed by Terraform"

  key_length           = 24
  shared_secret_length = 10
}

resource "mashery_package_plan" "oauth-load-test" {
  package_ref        = mashery_package.oauth.id
  name               = "OAuth Package"
  unlimited_qps      = true
  unlimited_quota    = true
  quota_period       = "day"
  admin_provisioning = true
}

resource "mashery_package_plan_service" "oauth-service" {
  package_plan_ref = mashery_package_plan.oauth-load-test.id
  service_ref      = mashery_service.srv.id
}

resource "mashery_package_plan_service_endpoint" "oauth-service-endpoint" {
  package_plan_service_ref = mashery_package_plan_service.oauth-service.id
  service_endpoint_ref     = mashery_service_endpoint.oauth-endpoint.id
}