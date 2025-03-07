resource "ovh_me_api_oauth2_client" "my_oauth2_client_client_creds" {
  name = "client credentials service account"
  description = "An OAuth2 client using the client credentials flow for my app"
  flow = "CLIENT_CREDENTIALS"
}
