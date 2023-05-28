terraform {
  required_providers {
    mashery = {
      version = "0.4"
      source = "github.com/aliakseiyanchuk/mashery"
    }
  }
}

provider "mashery" {
  log_file="./log"
  vault_addr = var.vault_url
  # For the test examples, the Mashery secret engine needs to be mounted on mash-auth path
  vault_mount = "mash-auth"
  vault_role = var.vault_role
}