---
subcategory : "Public IPs"
---

# ovh_cloud_public_ips (Data Source)

Use this data source to list all public IPs (additional, external network and floating IPs) of a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_public_ips" "all" {
  service_name = "<public cloud project ID>"
}

output "public_ips" {
  value = data.ovh_cloud_public_ips.all.public_ips
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.

## Attributes Reference

The following attributes are exported:

* `public_ips` - List of public IPs of the project. Each element exports:
  * `ip` - Public IP address.
  * `type` - Type of the public IP (`ADDITIONAL_IP`, `EXT_NET_IP`, `FLOATING_IP`).
