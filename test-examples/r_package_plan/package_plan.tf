data "mashery_role" "content_manager" {
  search = {
    name : "Content Manager"
  }
}

data "mashery_email_template_set" "admin_email_template" {
  search = {
    name : "Terraform X1"
  }
}

data "mashery_email_template_set" "user_email_template" {
  search = {
    name : "Terraform X2"
  }
}


resource "mashery_package" "oauth" {
  name_prefix="demo_oauth"
  description="Package configuration for OAuth having shared secret"
  near_quota_threshold = 80
  key_length=24
  shared_secret_length=12
}

resource "mashery_package_plan" "default" {
  package_id = mashery_package.oauth.id
  name = "Default"
  admin_provisioning = true
  portal_access_roles = toset([data.mashery_role.content_manager.id])

  email_template_set = data.mashery_email_template_set.user_email_template.id
  admin_email_template_set = data.mashery_email_template_set.admin_email_template.id
}