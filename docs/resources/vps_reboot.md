---
subcategory : "VPS"
---

# ovh_vps_reboot

Reboot your VPS. This is a one-shot resource: `Create` POSTs `/vps/{service_name}/reboot`
and waits for the resulting task to reach a terminal state. `Read` is a no-op, `Update`
is never invoked (every field is `ForceNew`), and `Delete` simply drops the task id.

Change any value in `triggers` to re-run the reboot.

## Example Usage

```terraform
resource "ovh_vps_reboot" "reboot" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  # Change any value here to re-run the reboot.
  triggers = {
    kernel = "6.6.1"
  }
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS service.
* `triggers` - (Optional, ForceNew) Map of arbitrary string values; changing any
  value re-runs the reboot.

## Attributes Reference

* `id` - The OVH task id (as a string).
* `task_id` - The OVH task id.
* `task_state` - Final state of the OVH task (typically `done`).
