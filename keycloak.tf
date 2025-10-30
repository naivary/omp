terraform {
  required_providers {
    keycloak = {
      source  = "keycloak/keycloak"
      version = "5.5.0"
    }
  }
}

provider "keycloak" {
  client_id = "admin-cli"
  username  = var.keycloak_username
  password  = var.keycloak_password
  url       = var.keycloak_url
  realm     = "master"
  tls_insecure_skip_verify = true
}

data "keycloak_openid_client" "realm_management" {
  realm_id  = keycloak_realm.omp.id
  client_id = "realm-management"
}

resource "keycloak_realm" "omp" {
  realm        = "omp"
  display_name = "Open Metrics Platform"
  enabled      = true
  // registration settings
  login_with_email_allowed       = true
  registration_email_as_username = true
  verify_email                   = true
  // login settings
  remember_me = true
}

resource "keycloak_realm_user_profile" "userprofile" {
  realm_id = keycloak_realm.omp.id

  attribute {
    name         = "email"
    display_name = "Email"
  }

  attribute {
    name         = "username"
    display_name = "Username"
  }


  attribute {
    name         = "profileID"
    display_name = "Profile ID"
    permissions {
      view = ["admin", "user"]
      edit = ["admin"]
    }
  }


}

resource "keycloak_openid_client" "omp_rest_api" {
  realm_id  = keycloak_realm.omp.id
  client_id = "omp-rest-api"
  name      = "Open Metrics Platform REST API"
  enabled   = true

  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  direct_access_grants_enabled = true
  service_accounts_enabled     = true
  client_secret = var.keycloak_omp_oidc_client_secret

  valid_redirect_uris = [
    "http://localhost:8080/openid-callback"
  ]
}

resource "keycloak_role" "realm_admin" {
  realm_id  = keycloak_realm.omp.id
  client_id = keycloak_openid_client.omp_rest_api.id
  name      = "realm-admin"
}

resource "keycloak_openid_client_service_account_role" "omp_rest_api_svc_acc_admin" {
  realm_id                = keycloak_realm.omp.id
  service_account_user_id = keycloak_openid_client.omp_rest_api.service_account_user_id
  client_id               = data.keycloak_openid_client.realm_management.id
  role                    = keycloak_role.realm_admin.name
}
