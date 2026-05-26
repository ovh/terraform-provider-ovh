---
subcategory : "VPS"
---

# ovh_vps_backup_ftp_access

Grants and manages a backup FTP access (ACL) on a VPS for a given CIDR block.
The OVHcloud API exposes one ACL entry per authorized CIDR block and per
protocol toggle (FTP, NFS, CIFS).

~> **WARNING** Backup FTP operations return a `dedicated.server.Task` which
is polled best-effort. If the underlying task endpoint cannot be reached,
the apply will not fail but the resource may show drift on the next refresh.

## Example Usage

```terraform
resource "ovh_vps_backup_ftp_access" "acl" {
  service_name = "vpsXXXXXX.ovh.net"
  ip_block     = "203.0.113.0/24"
  cifs         = true
  nfs          = false
  ftp          = false
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS.
* `ip_block` - (Required, ForceNew) CIDR block to grant access to.
* `cifs` - (Required) Whether CIFS (SMB) protocol is enabled for this block.
* `nfs` - (Required) Whether NFS protocol is enabled for this block.
* `ftp` - (Optional, default `false`) Whether FTP protocol is enabled for
  this block.

## Attributes Reference

* `id` - Composite id `service_name|ip_block`.
* `is_applied` - Whether the ACL is currently applied on the backup FTP
  storage.
* `last_update` - Last ACL update date.

## Import

Backup FTP ACLs can be imported using their composite id:

```bash
$ terraform import ovh_vps_backup_ftp_access.acl "vpsXXXXXX.ovh.net|203.0.113.0/24"
```
