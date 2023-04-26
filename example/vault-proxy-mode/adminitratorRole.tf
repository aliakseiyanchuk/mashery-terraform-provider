variable "role_name" {
  type = string
  default = "Administrator"
  description = "Role to identify"
}

data "mashery_role" "my_role" {
  search = {
    name: var.role_name
  }
}

output "roleId" {
  value = data.mashery_role.my_role.id
}