---
subcategory : "VPS"
---

# ovh_vps_automated_backup_attached (Data Source)

Use this data source to list automated backup restore points currently
attached (mounted) to a VPS, including their access information.

## Example Usage

```terraform
data "ovh_vps_automated_backup_attached" "att" {
  service_name = "vpsXXXXXX.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.

## Attributes Reference

* `attached_backups` - List of currently attached restore points.
  * `restore_point` - The restore point datetime (RFC3339).
  * `nfs` - The NFS access path (if any).
  * `smb` - The SMB access path (if any).
  * `additional_disk` - The additional disk identifier (if any).
