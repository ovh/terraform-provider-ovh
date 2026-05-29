---
subcategory : "VPS"
---

# ovh_vps_set_password

Reset the root password of your VPS. OVH generates the new password server-side
and, **by default, emails it** to the contact on file. If the VPS was installed
with `do_not_send_password = true`, OVH stores the password without emailing it
and the operator must retrieve it through their own channel.

One-shot resource: `Create` POSTs `/vps/{service_name}/setPassword` and waits for
the resulting task to reach a terminal state. `Read` and `Delete` are no-ops.
Change any value in `triggers` to rotate the password again.

~> **WARNING** The OVH API does not return the new password and does not report
whether it was emailed. The `password_sent_via_email` attribute is informational;
set it to `false` in your Terraform config if your VPS was installed with
`do_not_send_password = true`.

## Example Usage

```terraform
resource "ovh_vps_set_password" "pwd" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  # Re-run the action by changing any value in this map.
  triggers = {
    rotate = "2025-01-01"
  }
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS service.
* `triggers` - (Optional, ForceNew) Map of arbitrary string values; changing any
  value re-runs the password reset.
* `password_sent_via_email` - (Optional) Whether OVH emailed the new root
  password. Defaults to `true`. Set to `false` only if the VPS was installed
  with `do_not_send_password = true`. Informational; the API does not report
  this value.

## Attributes Reference

* `id` - The OVH task id (as a string).
* `task_id` - The OVH task id.
* `task_state` - Final state of the OVH task (typically `done`).
* `password_sent_via_email` - See above.

## Compatibility

This resource wraps `POST /vps/{serviceName}/setPassword`. Live cross-region probing on 2026-05-26 shows
the endpoint is present in the **EU** and **CA** API schemas (`eu.api.ovh.com`,
`ca.api.ovh.com`) but **NOT** in the **US** schema (`api.us.ovhcloud.com`).

On a US-region VPS the OVHcloud API returns
`404: Got an invalid (or empty) URL`. Use this resource on EU or CA accounts,
or wait for OVH to expose this endpoint on US.
