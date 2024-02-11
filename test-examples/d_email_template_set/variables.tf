variable "vault_url" {
  default = "http://localhost:8200"
  description = "Vault URL that has a "
}

variable "vault_role" {
  default = "sandbox"
  description = "Vault secret engine role to use"
}