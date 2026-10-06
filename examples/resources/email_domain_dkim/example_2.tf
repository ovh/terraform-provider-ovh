resource "ovh_email_domain_dkim" "my_dkim" {
  domain = "example.com"
}

# The selectors are only known once DKIM is enabled, and for_each needs its keys
# at plan time. OVHcloud allocates exactly two and the list is sorted, so a fixed
# count works on the first apply too.
resource "cloudflare_dns_record" "dkim_selectors" {
  count = 2

  zone_id = var.cloudflare_zone_id
  name    = "${ovh_email_domain_dkim.my_dkim.selectors[count.index].selector_name}._domainkey.example.com"
  type    = "CNAME"
  ttl     = 1
  proxied = false

  # cname is a full zone-file line, the target is its last field
  content = reverse(split(" ", ovh_email_domain_dkim.my_dkim.selectors[count.index].cname))[0]
}
