terraform {
  required_providers {
    mashery = {
      version = "0.0.1"
      source = "yanchuk.nl/aliakseiyanchuk/mashery"
    }
  }
}


provider "mashery" {
  log_file="./log"
}