---
subcategory : "Domain names"
---

# ovh_domain_zone_dnssec (Data Source)

Use this data source to retrieve information about a domain zone DNSSEC status.

## Example Usage

```hcl
data "ovh_domain_zone_dnssec" "dnssec" {
  zone_name = "mysite.ovh"
}
```

## Argument Reference

* `zone_name` - (Required) The name of the domain zone

## Attributes Reference

* `status` - DNSSEC status (`disableInProgress`, `disabled`, `enableInProgress` or `enabled`)
