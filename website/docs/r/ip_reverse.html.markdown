---
layout: "ovh"
page_title: "OVH: ovh_ip_reverse"
sidebar_current: "docs-ovh-resource-ip-reverse"
description: |-
    Provides a OVH IP reverse resource.
---

# ovh_ip_reverse

Provides a OVH IP reverse.

## Example Usage

```hcl
# Set the reverse of an IP
resource "ovh_ip_reverse" "test" {
    ip = "192.0.2.0/24"
    ipreverse = "192.0.2.1"
    reverse = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `ip` - (Required) The IP block to which the IP belongs
* `reverse` - (Required) The value of the reverse
* `ipreverse` - (Optional) The IP to set the reverse of, default to `ip` if `ip` is a /32 (IPv4) or a /128 (IPv6)

## Attributes Reference

The following attributes are exported:

* `ipreverse` - The IP to set the reverse of
* `reverse` - The value of the reverse
