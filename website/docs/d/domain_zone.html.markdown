---
subcategory : "Domain names"
---

# ovh_domain_zone (Data Source)

Use this data source to retrieve information about a domain zone.

## Example Usage

```hcl
data "ovh_domain_zone" "root_zone" {
  name = "mysite.ovh"
}
```

## Argument Reference

* `name` - (Required) The name of the domain zone.

## Attributes Reference

`id` is set to the domain zone name.
In addition, the following attributes are exported:

* `urn` - URN of the DNS zone
* `last_update` - Last update date of the DNS zone
* `has_dns_anycast` - hasDnsAnycast flag of the DNS zone
* `name_servers` - Name servers that host the DNS zone
* `dnssec_supported` - Is DNSSEC supported by this zone
