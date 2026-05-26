---
subcategory : "VPS"
---

# ovh_vps_snapshot_download (Data Source)

Retrieve a short-lived signed URL to **download** the current snapshot of an
OVHcloud VPS.

The URL is sensitive and time-limited; treat it like a credential.

## Example Usage

```terraform
data "ovh_vps_snapshot_download" "dl" {
  service_name = "vps-xxxxxx.vps.ovh.net"
}

output "snapshot_size_bytes" {
  value = data.ovh_vps_snapshot_download.dl.size
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.

## Attributes Reference

The following attributes are exported:

* `id` - Same as `service_name`.
* `url` - Short-lived signed download URL (sensitive).
* `size` - Size of the snapshot in bytes.
