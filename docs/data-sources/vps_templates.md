---
subcategory : "VPS"
---

# ovh_vps_templates (Data Source)

Use this data source to list the legacy OS templates installable on a VPS via
the legacy `/vps/{serviceName}/reinstall` flow. Most fleets on the `/rebuild`
flow should use `ovh_vps_available_images` instead; this datasource exists for
existing VPS still on legacy templates.

The list is fetched, then a fan-out GET is issued per template id (capped at
4 concurrent workers) so the returned `templates` block contains the full
template metadata.

## Example Usage

```terraform
data "ovh_vps_templates" "all" {
  service_name = "vps-xxxxxx.vps.ovh.net"
}

data "ovh_vps_templates" "debian_64" {
  service_name        = "vps-xxxxxx.vps.ovh.net"
  distribution_filter = "Debian"
  bit_format_filter   = 64
}
```

## Argument Reference

* `service_name` - (Required) The VPS service_name (ex: `vps-xxxxxx.vps.ovh.net`).
* `distribution_filter` - (Optional) Case-insensitive substring match against `distribution` (ex: `Debian`).
* `bit_format_filter` - (Optional) Filter on bit format. Must be `32` or `64`.

## Attributes Reference

* `template_ids` - List of template IDs matching the filters.
* `templates` - List of template objects. Each has:
  * `id` - Template ID (int).
  * `name` - Template name.
  * `distribution` - Distribution string (ex: `Debian 11`).
  * `bit_format` - 32 or 64 (int — the API returns this as a string, the provider converts it).
  * `available_language` - List of supported language codes.
  * `locale` - Default locale.

## Compatibility

This data source is part of the OVHcloud public VPS API schema, but the
underlying endpoint may not be implemented for every VPS lineup. Live
testing on a `vps-le-2-2-40` (2019v1) US instance returns
`404: Got an invalid (or empty) URL`. The endpoint is
exposed on EU and CA regions; on the US API it is only present for legacy VPS plans (pre-2019).

If you receive a 404 from this data source, your VPS plan likely does
not expose this endpoint. See [the OVHcloud VPS API console](https://api.us.ovhcloud.com/console/#/vps)
for what's available on your specific service.
