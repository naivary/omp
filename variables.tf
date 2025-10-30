variable "keycloak_url" {
  description = "URL of the keycloak instance"
  type        = string
}

variable "keycloak_username" {
  type = string
}

variable "keycloak_password" {
  type      = string
  sensitive = true
}

variable "keycloak_omp_oidc_client_secret" {
  type      = string
  sensitive = true
}
