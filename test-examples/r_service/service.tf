data "mashery_role" "internal_dev" {
  search = {
    "name": "Terraform Developer"
  }
}

resource "mashery_service" "srv" {
  name_prefix="tf-debug"
  iodocs_accessed_by = toset([data.mashery_role.internal_dev.id])
}

resource "mashery_service_oauth" "srv" {
  service_id = mashery_service.srv.id
}