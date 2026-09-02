variable "aliases" {
  type = map(string)
  default = {
    "postmaster" = "admin@example.com"
    "abuse"      = "admin@example.com"
    "sales"      = "someone.else@example.com"
  }
}

resource "ovh_email_domain_redirection" "aliases" {
  for_each = var.aliases

  domain = "example.com"
  from   = "${each.key}@example.com"
  to     = each.value
}
