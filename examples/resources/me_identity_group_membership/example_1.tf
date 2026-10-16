resource "ovh_me_identity_group" "my_group" {
  description = "My custom group"
  name        = "my_group"
  role        = "NONE"
}

resource "ovh_me_identity_user" "my_user" {
  description = "My custom user"
  email       = "my_user@example.com"
  group       = "DEFAULT"
  login       = "my_user"
  password    = "super-s3cr3t!password"
}

resource "ovh_me_identity_group_membership" "my_membership" {
  login = ovh_me_identity_user.my_user.login
  group = ovh_me_identity_group.my_group.name
}
