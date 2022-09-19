---
layout: "ovh"
page_title: "OVH: ovh_ip_reverse"
sidebar_current: "docs-ovh-resource-ip-reverse"
description: |-
    Provides a OVHcloud IP reverse resource.
---

# ovh_ip_reverse

Provides a OVHcloud IP reverse.

## Example Usage

```hcl
# Set the reverse of an IP
resource "ovh_ip_reverse" "test" {
    ip = "192.0.2.0/24"
    ip_reverse = "192.0.2.1"
    reverse = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `ip` - (Required) The IP block to which the IP belongs
* `reverse` - (Required) The value of the reverse
* `ip_reverse` - (Required) The IP to set the reverse of

## Attributes Reference

The id is set to the value of ip_reverse.
