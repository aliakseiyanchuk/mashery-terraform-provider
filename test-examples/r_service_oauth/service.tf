resource "mashery_service" "srv" {
  name_prefix="tf-debug"
}

resource "mashery_service_oauth" "srv_oauth" {
  service_id = mashery_service.srv.id
}

