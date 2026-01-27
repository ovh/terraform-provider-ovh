resource "ovh_me_identity_user" "user" {
  description = "User description"
  email       = "user.email@example.com"
  group       = "DEFAULT"
  login       = "user_login"
  password    = "SecretPassword123"
}

resource "ovh_me_identity_user_token" "token" {
  user_login  = ovh_me_identity_user.user.login
  name        = "token_name"
  description = "Token description"
  expires_at  = "2030-01-01T00:00:00Z"
}
