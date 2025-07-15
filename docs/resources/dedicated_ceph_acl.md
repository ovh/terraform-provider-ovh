---
subcategory : "Cloud Disk Array"
---

# ovh_dedicated_ceph_acl

Add a new access ACL for the given network/mask.

## Example Usage

```terraform
data "ovh_dedicated_ceph" "my_ceph" {
  service_name = "94d423da-0e55-45f2-9812-836460a19939"
}

resource "ovh_dedicated_ceph_acl" "my_acl" {
  service_name = data.ovh_dedicated_ceph.my_ceph.id
  network      = "1.2.3.4"
  netmask      = "255.255.255.255"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your dedicated CEPH
* `network` - (Required) The network IP to authorize
* `netmask` - (Required) The network mask to apply

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `network` - See Argument Reference above.
* `netmask` - See Argument Reference above.
* `family` - IP family. `IPv4` or `IPv6`
