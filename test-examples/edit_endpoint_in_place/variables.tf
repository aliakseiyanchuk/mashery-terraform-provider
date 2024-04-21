variable "vault_url" {
  default = "http://localhost:8200"
  description = "Vault URL to read data the data from; defaults to the development server."
}

variable "vault_role" {
  default = "sandbox"
  description = "Vault secret engine role to use"
}

variable "service_name" {
  description = "Service name where an endpoint should be edited"
}