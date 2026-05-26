---
subcategory : "VPS"
---

# ovh_vps_snapshot (Resource)

Manage a **snapshot** of an OVHcloud VPS.

A VPS can hold at most one snapshot at a time. The targeted VPS must have the
**snapshot** option subscribed; otherwise the OVHcloud API will refuse the
`POST /vps/{serviceName}/createSnapshot` call and the provider will surface a
clear error.

Creating this resource calls `POST /vps/{serviceName}/createSnapshot` and
waits for the resulting task to reach state `done`. Updating the resource
calls `PUT /vps/{serviceName}/snapshot` (only the `description` field is
writable). Deleting calls `DELETE /vps/{serviceName}/snapshot` and waits for
the resulting task to complete.

## Example Usage

```terraform
resource "ovh_vps_snapshot" "snap" {
  service_name = "vps-xxxxxx.vps.ovh.net"
  description  = "snapshot taken before upgrade"
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The internal name of your VPS. The VPS
  must have the *snapshot* option subscribed.
* `description` - (Optional) Human-readable description of the snapshot. This
  is the only field that can be changed in-place.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The internal name of the VPS (same value as `service_name`).
* `creation_date` - RFC3339 timestamp of when the snapshot was created.
* `region` - OVHcloud region where the snapshot is stored.

## Import

VPS snapshot can be imported using the `service_name` of the VPS, e.g.:

```bash
$ terraform import ovh_vps_snapshot.snap vps-xxxxxx.vps.ovh.net
```
