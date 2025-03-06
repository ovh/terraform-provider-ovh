resource "ovh_me_api_oauth2_client" "my_oauth2_client_auth_code" {
  name = "OAuth2 authorization code service account"
  flow = "AUTHORIZATION_CODE"
  description = "An OAuth2 client using the authorization code flow for my-app.com"
  callback_urls = ["https://my-app.com/callback"]
}
