data "mashery_package" "root_pack" {
  search = {
    name = var.package_name
  }
}

output "package_id" {
  value = data.mashery_package.root_pack.id
}