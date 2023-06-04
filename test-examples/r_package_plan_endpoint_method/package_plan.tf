
resource "mashery_package" "oauth" {
  name_prefix="demo_oauth"
  description="Package configuration for OAuth having shared secret which uses"
  near_quota_threshold = 80
  key_length=24
  shared_secret_length=12
}

resource "mashery_package_plan" "default" {
  package_id = mashery_package.oauth.id
  name = "Default"
  admin_provisioning = true
}

resource "mashery_package_plan_service" "svc" {
  plan_id = mashery_package_plan.default.id
  service_id = mashery_service.srv.id
}

resource "mashery_package_plan_service_endpoint" "endp" {
  plan_service_id = mashery_package_plan_service.svc.id
  endpoint_id = mashery_service_endpoint.endp.id
}

resource "mashery_package_plan_service_endpoint_method" "meth_abc" {
  package_plan_service_endpoint_id = mashery_package_plan_service_endpoint.endp.id
  method_id = mashery_endpoint_method.meth_abc.id
  service_filter_id = mashery_endpoint_method_filter.abc_filter.id
}