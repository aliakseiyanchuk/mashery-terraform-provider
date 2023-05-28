data "mashery_email_template_set" "default" {
  search = {
    "name": "Default"
  }
}

output "set_id" {
  value = data.mashery_email_template_set.default.id
}