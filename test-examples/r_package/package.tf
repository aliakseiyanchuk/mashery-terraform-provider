# Example of mashery package configurations

data "mashery_organization" "tf_org" {
  search = {
    "name": "Terraform"
  }
}


resource "mashery_package" "apikey" {
  name_prefix="demo"
  description="Package configuration for API key-only configurations"
  organization = data.mashery_organization.tf_org.id
}