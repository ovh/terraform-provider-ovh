---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_veeam_restored_backup"
description: |-
  Get information about the currently restored Veeam backup for an OVHcloud VPS.
---

# ovh_vps_veeam_restored_backup (Data Source)

Use this data source to retrieve information about the currently restored Veeam backup for an OVHcloud VPS.

## Example Usage

```hcl
data "ovh_vps_veeam_restored_backup" "rb" {
  service_name = "vps-xxxxxxxx.vps.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your VPS.

## Attributes Reference

* `id` - The ID of the restored backup.
* `access_ip` - IP address to access the restored backup.
* `nfs_url` - NFS URL for the restored backup.
* `restore_point_id` - The restore point ID used for this restore.
* `state` - State of the restored backup.
