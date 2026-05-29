---
subcategory : "VPS"
---

# ovh_vps_snapshot_revert

Reverts a VPS to its current snapshot. This is a one-shot action resource: applying
it issues `POST /vps/{service_name}/snapshot/revert`, polls the returned task to
completion, and stores the result. The VPS must have a snapshot already taken
(via `ovh_vps_snapshot`).

To re-run the revert later — for example after taking a fresh snapshot — change
any value in `triggers`; Terraform will destroy and recreate the resource, which
runs the revert again.

`terraform destroy` only drops the resource from state. There is no API to undo
a revert.

## Example Usage

```hcl
resource "ovh_vps_snapshot_revert" "rollback" {
  service_name = "vpsXXXXX.ovh.net"

  triggers = {
    after_snapshot_id = ovh_vps_snapshot.daily.id
  }
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS service.
* `triggers` - (Optional, ForceNew) Map of arbitrary string values; changing any
  value re-runs the revert.

## Attributes Reference

* `task_id` - ID of the OVH task that performed the revert.
* `task_state` - Terminal state of the task.
* `reverted_at` - RFC3339 timestamp of when the revert was issued.
