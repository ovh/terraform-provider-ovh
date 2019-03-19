---
layout: "ovh"
page_title: "OVH: domain_zone"
sidebar_current: "docs-ovh-datasource-domain-zone"
description: |-
  Get information & status of a domain zone.
---

# ovh_domain_zone

Use this data source to retrieve information about a domain zone.

## Example Usage

```hcl
data "ovh_domain_zone" "rootzone" {
    name = "mysite.ovh"
}
```

## Argument Reference

* `name` - (Required) The name of the domain zone.

## Attributes Reference

`id` is set to the domain zone name.
In addition, the following attributes are exported:

* `last_update` - Last update date of the DNS zone
* `has_dns_anycast` - hasDnsAnycast flag of the DNS zone
* `name_servers` - Name servers that host the DNS zone
* `dnssec_supported` - Is DNSSEC supported by this zone
