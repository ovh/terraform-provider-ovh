---
subcategory : "VPS"
---

# ovh_vps_backup_ftp_access (Data Source)

Use this data source to read a single Backup FTP ACL entry attached to a
VPS service.

## Example Usage

```terraform
data "ovh_vps_backup_ftp_access" "entry" {
  service_name = "vpsXXXXXX.vps.ovh.net"
  ip_block     = "203.0.113.0/24"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.
* `ip_block` - (Required) CIDR-formatted IP block identifying the ACL entry.

## Attributes Reference

`id` is set to `<service_name>|<ip_block>`.

In addition, the following attributes are exported:

* `ftp` - Whether FTP access is granted to this IP block.
* `cifs` - Whether CIFS / SMB access is granted to this IP block.
* `nfs` - Whether NFS access is granted to this IP block.
* `is_applied` - Whether the ACL entry has been applied on the backup FTP server.
* `last_update` - Timestamp of the last ACL update.
