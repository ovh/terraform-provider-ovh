---
subcategory : "VPS"
---

# ovh_vps_task (Data Source)

Use this data source to retrieve information about a specific task running on a VPS associated with your OVHcloud Account.

## Example Usage

```terraform
data "ovh_vps_task" "task" {
  service_name = "vps-XXXXXX.vps.ovh.net"
  id           = 12345
}
```

## Argument Reference

* `service_name` - (Required) The service name of your VPS (ex: "vps-123456.vps.ovh.net").
* `id` - (Required) The numeric ID of the task to look up.

## Attributes Reference

`id` is set to the numeric ID of the task.

The following attributes are exported:

* `date` - The creation date of the task (RFC3339).
* `type` - The task type (ex: "rebootVm", "reinstallVm").
* `state` - The current state of the task. One of: `blocked`, `cancelled`, `doing`, `done`, `error`, `paused`, `todo`, `waitingAck`.
* `progress` - The completion progress, as an integer percentage.
