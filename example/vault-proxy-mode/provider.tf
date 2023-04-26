terraform {
  required_providers {
    mashery = {
      version = "0.0.1"
      source = "yanchuk.nl/aliakseiyanchuk/mashery"
    }
    vault = {
      source = "hashicorp/vault"
      version = "~> 2.17.0"
    }

  }
}

provider "vault" {
  address = var.vault_url
}


provider "mashery" {
  log_file=".out/log_tf_mashery"
  vault_addr = var.vault_url
  vault_role = var.vault_role
}
