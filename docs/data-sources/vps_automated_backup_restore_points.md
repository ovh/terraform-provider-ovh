---
subcategory : "VPS"
---

# ovh_vps_automated_backup_restore_points (Data Source)

Use this data source to list automated backup restore points for a VPS.

## Example Usage

```terraform
data "ovh_vps_automated_backup_restore_points" "rp" {
  service_name = "vpsXXXXXX.ovh.net"
  state        = "available"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.
* `state` - (Optional) Filter restore points by state. One of `available`, `restored`, `restoring`.

## Attributes Reference

* `restore_points` - List of restore points as RFC3339 datetime strings.
