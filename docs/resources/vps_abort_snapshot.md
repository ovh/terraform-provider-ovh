---
subcategory : "VPS"
---

# ovh_vps_abort_snapshot

Cancel an in-flight snapshot or automated-backup operation on a VPS by POSTing
`/vps/{service_name}/abortSnapshot`. This is a one-shot, fire-and-forget
resource: the OVH API does not return a task, the abort happens server-side
asynchronously, and Terraform only records the timestamp at which the call was
issued.

Lifecycle:

* `Create` POSTs the endpoint and writes the current RFC3339 timestamp to
  `aborted_at`.
* `Read` is a no-op.
* `Update` is never invoked — every field is `ForceNew`.
* `Delete` simply drops the resource id; the abort itself cannot be undone.

Change any value in `triggers` to re-run the abort. If no snapshot or
automated-backup operation is currently in progress, the OVH API will return
an error and creation will fail with a user-friendly message.

## Example Usage

```terraform
resource "ovh_vps_abort_snapshot" "abort" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  # Change any value here to re-run the abort.
  triggers = {
    run = "2026-05-25T00:00:00Z"
  }
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS service.
* `triggers` - (Optional, ForceNew) Map of arbitrary string values; changing
  any value re-runs the abort.

## Attributes Reference

* `id` - A composite identifier of the form `{service_name}/abortSnapshot/{timestamp}`.
* `aborted_at` - RFC3339 timestamp of when the abort was issued.
