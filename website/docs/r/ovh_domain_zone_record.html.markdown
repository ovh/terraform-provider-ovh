---
layout: "ovh"
page_title: "OVH: ovh_domain_zone_record"
sidebar_current: "docs-ovh-resource-domain-zone-record"
description: |-
  Provides a OVHcloud domain zone resource.
---

# ovh_domain_zone_record

Provides a OVHcloud domain zone record.

## Example Usage

```hcl
# Add a record to a sub-domain
resource "ovh_domain_zone_record" "test" {
  zone      = "testdemo.ovh"
  subdomain = "test"
  fieldtype = "A"
  ttl       = "3600"
  target    = "0.0.0.0"
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The domain to add the record to
* `subdomain` - (Required) The name of the record
* `target` - (Required) The value of the record
* `fieldtype` - (Required) The type of the record
* `ttl` - (Optional) The TTL of the record, it shall be >= to 60.


## Attributes Reference

The following attributes are exported:

* `id` - The record ID
* `zone` - The domain to add the record to
* `subDomain` - The name of the record
* `target` - The value of the record
* `fieldType` - The type of the record
* `ttl` - The TTL of the record

## Import
OVHcloud domain zone record can be imported using the `id`, which can be retrieved by using [OVH API portal](https://api.ovh.com/console/#/domain/zone/%7BzoneName%7D/record~GET), and the `zone`, separated by "." E.g.,


```bash
$ terraform import ovh_domain_zone_record.test id.zone
```
