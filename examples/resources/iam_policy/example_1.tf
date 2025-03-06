data "ovh_me" "account" {}

resource "ovh_me_identity_group" "my_group" {
  name        = "my_group"
  description = "my_group created in Terraform"
}

resource "ovh_iam_policy" "manager" {
  name        = "allow_ovh_manager"
  description = "Users are allowed to use the OVH manager"
  identities  = [ovh_me_identity_group.my_group.urn]
  resources   = [data.ovh_me.account.urn]
  # these are all the actions
  allow = [
    "account:apiovh:me/get",
    "account:apiovh:me/supportLevel/get",
    "account:apiovh:me/certificates/get",
    "account:apiovh:me/tag/get",
    "account:apiovh:services/get",
    "account:apiovh:*",
  ]
}
