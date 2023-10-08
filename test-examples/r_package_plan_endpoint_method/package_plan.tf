
resource "mashery_package" "oauth" {
  name_prefix="demo_oauth"
  description="Package configuration for OAuth having shared secret which uses"
  near_quota_threshold = 80
  key_length=24
  shared_secret_length=12
}

resource "mashery_package_plan" "default" {
  package_ref = mashery_package.oauth.id
  name = "Default"
  admin_provisioning = true
}

resource "mashery_package_plan_service" "svc" {
  plan_ref = mashery_package_plan.default.id
  service_ref = mashery_service.srv.id
}

resource "mashery_package_plan_service_endpoint" "endp" {
  plan_service_ref = mashery_package_plan_service.svc.id
  endpoint_ref = mashery_service_endpoint.endp.id
}

resource "mashery_package_plan_service_endpoint_method" "meth_abc" {
  package_plan_service_endpoint_ref = mashery_package_plan_service_endpoint.endp.id
  method_ref = mashery_service_endpoint_method.meth_abc.id
  service_filter_ref = mashery_service_endpoint_method_filter.abc_filter.id
}