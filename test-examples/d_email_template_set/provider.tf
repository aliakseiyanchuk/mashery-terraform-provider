terraform {
  required_providers {
    mashery = {
      version = "0.5"
      source = "github.com/aliakseiyanchuk/mashery"
    }
  }
}

provider "mashery" {
  vault_addr = var.vault_url
  # For the test examples, the Mashery secret engine needs to be mounted on mash-auth path
  vault_mount = "mash-auth"
  role = var.vault_role
  qps = 1
}