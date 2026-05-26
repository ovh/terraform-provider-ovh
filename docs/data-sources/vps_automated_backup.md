---
subcategory : "VPS"
---

# ovh_vps_automated_backup (Data Source)

Use this data source to retrieve the automated backup configuration of a VPS.

## Example Usage

```terraform
data "ovh_vps_automated_backup" "ab" {
  service_name = "vpsXXXXXX.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.

## Attributes Reference

* `state` - The current state of automated backups (`enabled` or `disabled`).
* `schedule` - The time of day at which the automated backup runs (`HH:MM:SS`).
* `rotation` - The number of automated backups kept on rotation.
* `service_resource_name` - The resource name of the backup service attached to the VPS.
