resource "ovh_me_api_oauth2_client" "my_oauth2_client_no_secret" {
  name                  = "service account without client secret"
  description           = "An OAuth2 client whose secret is not kept in the Terraform state"
  flow                  = "CLIENT_CREDENTIALS"
  discard_client_secret = true
}
