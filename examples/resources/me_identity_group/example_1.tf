resource "ovh_me_identity_group" "my_group" {
  description = "Some custom description"
  name        = "my_group_name"
  role        = "NONE"
}
