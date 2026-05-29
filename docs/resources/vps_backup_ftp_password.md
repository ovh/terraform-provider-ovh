---
subcategory : "VPS"
---

# ovh_vps_backup_ftp_password

Rotates the Backup FTP password attached to a VPS service.

This is a one-shot resource: each apply with a changed `triggers` map (or a
fresh resource block) issues a `POST /vps/{service_name}/backupftp/password`
call. The new password is delivered out-of-band by OVHcloud (email / control
panel). Deleting the resource from state does not "unrotate" anything.

~> **WARNING** The OVH API returns a `dedicated.server.Task` for this
endpoint. Polling is best-effort across the `/vps/{sn}/tasks/{id}` and
`/dedicated/server/{sn}/task/{id}` paths; if polling cannot converge before
the timeout, the apply still succeeds with `task_state = "unknown"`.

## Example Usage

```terraform
resource "ovh_vps_backup_ftp_password" "rotate" {
  service_name = "vpsXXXXXX.vps.ovh.net"

  triggers = {
    # bump this value to rotate the password again
    rotation = "1"
  }
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS.
* `triggers` - (Optional, ForceNew) Arbitrary map of strings; changing any
  value forces a new resource (and therefore a fresh password rotation).

## Attributes Reference

* `id` - The task id returned by the API.
* `task_id` - The numeric id of the dedicated server task.
* `task_state` - Final state of the rotation task (`done`, `error`, or
  `unknown` if best-effort polling could not converge).

## Compatibility

This resource wraps `POST /vps/{serviceName}/backupftp/password`. Live cross-region probing on 2026-05-26 shows
the endpoint is present in the **EU** and **CA** API schemas (`eu.api.ovh.com`,
`ca.api.ovh.com`) but **NOT** in the **US** schema (`api.us.ovhcloud.com`).

On a US-region VPS the OVHcloud API returns
`404: Got an invalid (or empty) URL`. Use this resource on EU or CA accounts,
or wait for OVH to expose this endpoint on US.
