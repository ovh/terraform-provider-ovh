resource "ovh_email_domain_redirection" "postmaster" {
  domain = "example.com"
  from   = "postmaster@example.com"
  to     = "admin@example.com"
}
