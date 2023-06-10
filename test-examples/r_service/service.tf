data "mashery_role" "internal_dev" {
  search = {
    "name": "Terraform Developer"
  }
}

data "mashery_organization" "tf_org" {
  search = {
    "name": "Terraform"
  }
}

resource "mashery_service" "srv" {
  name_prefix="tf-debug"
  iodocs_accessed_by = toset([data.mashery_role.internal_dev.id])
  organization = data.mashery_organization.tf_org.id
}

resource "mashery_service_oauth" "srv" {
  service_ref = mashery_service.srv.id
}