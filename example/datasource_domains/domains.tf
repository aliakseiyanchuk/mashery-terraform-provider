data mashery_public_domains "myArea" {
}

data mashery_system_domains "myArea" {
}

output "out_public_domains" {
  value = data.mashery_public_domains.myArea.domains
}

output "out_system_domains" {
  value = data.mashery_system_domains.myArea.domains
}