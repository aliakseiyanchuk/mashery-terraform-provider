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
  vault_mount = "mash-auth"
  vault_role = var.vault_role
  qps = 1
}