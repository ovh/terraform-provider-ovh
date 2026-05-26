---
subcategory : "VPS"
---

# ovh_vps_available_images (Data Source)

Use this data source to list the OS images available for (re)installation
on a VPS associated with your OVHcloud account. This is the canonical way
to discover a valid `image_id` for use with `ovh_vps` resources without
having to look the value up manually in the OVHcloud control panel.

## Example Usage

List every available image:

```terraform
data "ovh_vps_available_images" "all" {
  service_name = "vps-XXXXXX.vps.ovh.net"
}

output "image_ids" {
  value = data.ovh_vps_available_images.all.image_ids
}
```

Filter by image name with a regex (e.g. only Debian images):

```terraform
data "ovh_vps_available_images" "debian" {
  service_name = "vps-XXXXXX.vps.ovh.net"
  name_pattern = "(?i)debian"
}
```

## Argument Reference

* `service_name` - (Required) The service name of your VPS (e.g.
  `vps-XXXXXX.vps.ovh.net`).
* `name_pattern` - (Optional) A Go regular expression applied to the
  image `name` field. When set, only matching images are returned.

## Attributes Reference

* `image_ids` - The list of image IDs available for the VPS, after any
  `name_pattern` filtering has been applied.
* `images` - The list of image objects available for the VPS. Each entry
  exposes:
  * `id` - The image ID, suitable for use as `image_id` on the `ovh_vps`
    resource.
  * `name` - The human-readable image name (e.g. `Debian 12`).

If a per-image lookup fails (e.g. transient API error) the offending
entry is skipped rather than failing the whole data source; a warning
is emitted to the Terraform log.
