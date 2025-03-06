resource "ovh_me_identity_user" "my_user" {
  description = "Some custom description"
  email       = "my_login@example.com"
  group       = "DEFAULT"
  login       = "my_login"
  password    = "super-s3cr3t!password"
}
