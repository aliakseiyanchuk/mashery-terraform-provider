
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
}