---
subcategory : "VPS"
---

# ovh_vps_tasks (Data Source)

Use this data source to list the IDs of tasks running on a VPS associated with your OVHcloud Account, optionally filtering by state and/or type.

## Example Usage

```terraform
data "ovh_vps_tasks" "tasks" {
  service_name = "vps-XXXXXX.vps.ovh.net"
  state_filter = "doing"
  type_filter  = "rebootVm"
}
```

## Argument Reference

* `service_name` - (Required) The service name of your VPS (ex: "vps-123456.vps.ovh.net").
* `state_filter` - (Optional) Filter tasks by state. One of: `blocked`, `cancelled`, `doing`, `done`, `error`, `paused`, `todo`, `waitingAck`.
* `type_filter` - (Optional) Filter tasks by type. One of: `addVeeamBackupJob`, `changeRootPassword`, `createSnapshot`, `deleteSnapshot`, `deliverVm`, `getConsoleUrl`, `internalTask`, `migrate`, `openConsoleAccess`, `provisioningAdditionalIp`, `reOpenVm`, `rebootVm`, `reinstallVm`, `removeVeeamBackup`, `rescheduleAutoBackup`, `restoreFullVeeamBackup`, `restoreVeeamBackup`, `restoreVm`, `revertSnapshot`, `setMonitoring`, `setNetboot`, `startVm`, `stopVm`, `upgradeVm`.

## Attributes Reference

The following attributes are exported:

* `task_ids` - The sorted list of matching task IDs.
