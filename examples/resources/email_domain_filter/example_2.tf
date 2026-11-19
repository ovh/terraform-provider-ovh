variable "mailboxes" {
  type    = set(string)
  default = ["alice", "bob", "carol"]
}

resource "ovh_email_domain_filter" "forged" {
  for_each = var.mailboxes

  domain       = "example.com"
  account_name = each.key
  name         = "Forged example.com"
  priority     = 2
  action       = "account"
  action_param = "quarantine@example.com"

  rules = [
    {
      header  = "From"
      operand = "contains"
      value   = "@example.com"
    },
    {
      header  = "Authentication-Results"
      operand = "contains"
      value   = "dmarc=fail"
    },
  ]
}
