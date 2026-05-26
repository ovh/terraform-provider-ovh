---
subcategory : "VPS"
---

# ovh_vps_automated_backup_restore

Triggers a restore from an automated backup on a VPS via
`/vps/{serviceName}/automatedBackup/restore` and waits for the underlying task
to reach a terminal state.

When the resource is destroyed and the restore type was `file`, the attached
restore point is detached via
`/vps/{serviceName}/automatedBackup/detachBackup`. For `full` restores there
is nothing to detach, so destroy only removes the resource from state.

All inputs are `ForceNew` — changing any of them recreates the resource and
triggers a new restore.

## Example Usage

```terraform
resource "ovh_vps_automated_backup_restore" "r" {
  service_name  = "vpsXXXXXX.ovh.net"
  restore_point = "2024-01-15T02:00:00Z"
  type          = "file"
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS.
* `restore_point` - (Required, ForceNew) The restore point to restore from (RFC3339 datetime).
* `type` - (Required, ForceNew) Restore type. One of `file` (mount as access) or `full` (overwrite VPS).
* `change_password` - (Optional, ForceNew) For `full` restores, whether to also reset the root password.

## Attributes Reference

* `task_id` - The id of the OVHcloud VPS task that performed the restore.
