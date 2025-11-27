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

resource "ovh_iam_policy" "ip_prod_access" {
  name        = "ip_prod_access"
  description = "Allow access only from a specific IP to resources tagged prod"
  identities  = [ovh_me_identity_group.my_group.urn]
  resources   = ["urn:v1:eu:resource:vps:*"]

  allow = [
    "vps:apiovh:*",
  ]

  conditions {
    operator = "MATCH"
    values = {
      "resource.Tag(environment)" = "prod"
      "request.IP"                = "192.72.0.1"
    }
  }
}

resource "ovh_iam_policy" "workdays_expiring" {
  name        = "workdays_expiring"
  description = "Allow access only on workdays, expires end of 2026"
  identities  = [ovh_me_identity_group.my_group.urn]
  resources   = ["urn:v1:eu:resource:vps:*"]

  allow = [
    "vps:apiovh:*",
  ]

  conditions {
    operator = "MATCH"
    values = {
      "date(Europe/Paris).WeekDay.In" = "monday,tuesday,wednesday,thursday,friday"
    }
  }

  expired_at = "2026-12-31T23:59:59Z"
}
