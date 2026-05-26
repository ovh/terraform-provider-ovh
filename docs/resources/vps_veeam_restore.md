---
subcategory: "VPS"
page_title: "OVHcloud: ovh_vps_veeam_restore"
description: |-
  Triggers a Veeam restore operation for an OVHcloud VPS.
---

# ovh_vps_veeam_restore

Triggers a Veeam restore operation for an OVHcloud VPS.

## Example Usage

```hcl
resource "ovh_vps_veeam_restore" "restore" {
  service_name     = "vps-xxxxxxxx.vps.ovh.net"
  restore_point_id = 12345
  type             = "full"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service_name of your VPS.
* `restore_point_id` - (Required) The ID of the Veeam restore point to restore from.
* `type` - (Required) The type of restore (file┃full).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the restore operation.
* `access_ip` - IP address to access the restored backup (for file restores).
* `nfs_url` - NFS URL for the restored backup (for file restores).
* `state` - State of the restore operation.
