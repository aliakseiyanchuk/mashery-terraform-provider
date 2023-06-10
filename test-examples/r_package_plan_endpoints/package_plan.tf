resource "mashery_package" "oauth" {
  name_prefix          = "demo_oauth"
  description          = "Package configuration for OAuth having shared secret"
  near_quota_threshold = 80
  key_length           = 24
  shared_secret_length = 12
}

resource "mashery_package_plan" "default" {
  package_ref        = mashery_package.oauth.id
  name               = "Default"
  admin_provisioning = true
}

resource "mashery_package_plan_service" "oauth-service" {
  plan_ref    = mashery_package_plan.default.id
  service_ref = mashery_service.svc.id
}

resource "mashery_package_plan_service_endpoint" "oauth-service-endpoint" {
  plan_service_ref = mashery_package_plan_service.oauth-service.id
  endpoint_ref     = mashery_service_endpoint.svc-endpoint.id
}