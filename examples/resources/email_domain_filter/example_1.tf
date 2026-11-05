resource "ovh_email_domain_filter" "viruses" {
  domain       = "example.com"
  account_name = "john.doe"
  name         = "Delete viruses"
  priority     = 1
  action       = "delete"

  rules = [
    {
      header  = "X-Virus-Tag"
      operand = "contains"
      value   = "YES"
    },
  ]
}
