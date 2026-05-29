---
subcategory : "VPS"
---

# ovh_vps_automated_backup_reschedule

Reschedules the automated backup of an OVHcloud VPS. The automated backup option
must already be enabled on the VPS - this resource only retunes the time of day
at which the daily backup runs; it cannot enable or disable the option.

It wraps the API call `POST /vps/{serviceName}/automatedBackup/reschedule` and
polls the resulting `vps.Task` to completion.

## Example Usage

```terraform
resource "ovh_vps_automated_backup_reschedule" "schedule" {
  service_name = "vpsXXXXXX.vps.ovh.net"
  schedule     = "02:00:00"
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS.
* `schedule` - (Required) Backup schedule expressed as an `HH:MM:SS` time of day
  (UTC), e.g. `"02:00:00"`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `rotation` - Number of backups retained by the automated backup policy.
* `state` - State of the automated backup option (`enabled` or `disabled`).

## Import

The reschedule resource can be imported using the VPS `service_name`:

```sh
terraform import ovh_vps_automated_backup_reschedule.schedule vpsXXXXXX.vps.ovh.net
```

## Notes

`terraform destroy` is a no-op: this endpoint can only retune an existing
schedule, not disable automated backups. To fully disable automated backups
on a VPS, use the VPS option management endpoints or the OVHcloud control
panel.
