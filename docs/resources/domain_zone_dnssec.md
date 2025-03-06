---
subcategory : "Domain names"
---

# ovh_domain_zone_dnssec

Enable / disable DNSSEC on a domain zone.

## Example Usage

```terraform
resource "ovh_domain_zone_dnssec" "dnssec" {
  zone_name = "mysite.ovh"
}
```

## Argument Reference

* `zone_name` - (Required) The name of the domain zone

## Attributes Reference

* `status` - DNSSEC status (`disableInProgress`, `disabled`, `enableInProgress` or `enabled`)
