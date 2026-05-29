---
subcategory : "VPS"
---

# ovh_vps_snapshot (Data Source)

Read information about the current snapshot held by an OVHcloud VPS.

A VPS holds at most one snapshot at a time. If no snapshot exists, the
underlying API call returns an error.

## Example Usage

```terraform
data "ovh_vps_snapshot" "snap" {
  service_name = "vps-xxxxxx.vps.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.

## Attributes Reference

The following attributes are exported:

* `id` - Same as `service_name`.
* `description` - Description set on the snapshot.
* `creation_date` - RFC3339 timestamp of when the snapshot was created.
* `region` - OVHcloud region where the snapshot is stored.
