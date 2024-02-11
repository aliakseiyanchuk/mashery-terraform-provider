variable "vault_url" {
  default = "http://localhost:8200"
  description = "Vault URL to read data the data from; defaults to the development server."
}

variable "vault_role" {
  default = "sandbox"
  description = "Vault secret engine role to use"
}

variable "package_name" {
  default = "test package"
  description = "Name of the package to be found"
}
variable "package_plan_name" {
  default = "Default"
  description = "Name of the plan to find in this package plan"
}