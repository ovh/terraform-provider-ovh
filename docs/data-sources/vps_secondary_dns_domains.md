---
subcategory : "VPS"
---

# ovh_vps_secondary_dns_domains (Data Source)

Lists all secondary DNS domains slaved on a VPS.

## Example Usage

```terraform
data "ovh_vps_secondary_dns_domains" "all" {
  service_name = "vpsXXXXX.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.

## Attributes Reference

* `result` - Sorted list of slaved domains.
