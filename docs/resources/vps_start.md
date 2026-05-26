---
subcategory : "VPS"
---

# ovh_vps_start

Power on your VPS. One-shot resource: `Create` POSTs `/vps/{service_name}/start`
and waits for the task to reach a terminal state. `Read` and `Delete` are no-ops.

Change any value in `triggers` to re-run the action.

## Example Usage

```terraform
resource "ovh_vps_start" "start" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  triggers = {
    nonce = "1"
  }
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS service.
* `triggers` - (Optional, ForceNew) Map of arbitrary string values; changing any
  value re-runs the start.

## Attributes Reference

* `id` - The OVH task id (as a string).
* `task_id` - The OVH task id.
* `task_state` - Final state of the OVH task (typically `done`).
