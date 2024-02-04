data "mashery_package" "root_pack" {
  search = {
    name = var.package_name
  }
}

data "mashery_package_plan" "root_pack" {
  package_ref = data.mashery_package.root_pack.id
  search = {
    name = var.package_plan_name
  }
}

output "package_id" {
  value = data.mashery_package.root_pack.id
}

output "package_plan_id" {
  value = data.mashery_package_plan.root_pack.id
}