resource "ovh_email_domain_dkim" "my_dkim" {
  domain = "example.com"
}

resource "ovh_domain_zone_record" "dkim_selectors" {
  for_each = {
    for selector in ovh_email_domain_dkim.my_dkim.selectors : selector.selector_name => selector
  }

  zone      = "example.com"
  subdomain = "${each.key}._domainkey"
  fieldtype = "CNAME"
  ttl       = 3600

  # cname is a full zone-file line, the target is its last field
  target = element(split(" ", each.value.cname), length(split(" ", each.value.cname)) - 1)
}
