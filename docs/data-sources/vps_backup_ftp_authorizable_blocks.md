---
subcategory : "VPS"
---

# ovh_vps_backup_ftp_authorizable_blocks (Data Source)

Use this data source to list CIDR blocks that can be granted backup FTP
access on a VPS.

## Example Usage

```terraform
data "ovh_vps_backup_ftp_authorizable_blocks" "blocks" {
  service_name = "vpsXXXXXX.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.

## Attributes Reference

* `id` - Set to the `service_name`.
* `blocks` - List of CIDR blocks authorized to be granted backup FTP access.
