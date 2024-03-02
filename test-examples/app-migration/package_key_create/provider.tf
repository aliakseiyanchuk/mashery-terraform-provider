terraform {
  required_providers {
    mashery = {
      version = "0.5"
      source  = "github.com/aliakseiyanchuk/mashery"
    }
    vault = {
      source = "hashicorp/vault"
      version = "3.24.0"
    }
  }
}

provider "vault" {
  # Configuration options
  address = var.vault_url
#  skip_child_token = true
}

provider "mashery" {
  vault_addr  = var.vault_url
  vault_mount = "mash-auth"
  role        = var.vault_role
  qps         = 1
}