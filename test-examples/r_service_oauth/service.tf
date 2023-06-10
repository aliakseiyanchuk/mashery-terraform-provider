resource "mashery_service" "srv" {
  name_prefix="tf-debug"
}

resource "mashery_service_oauth" "svc_oauth" {
  service_ref = mashery_service.srv.id
}
