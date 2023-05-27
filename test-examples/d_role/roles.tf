data "mashery_role" "internal_dev" {
  search = {
    "name": "Terraform Developer"
  }
}

output "role_id" {
  value = data.mashery_role.internal_dev.id
}