---
subcategory : "VPS"
---

# ovh_vps_template (Data Source)

Use this data source to retrieve details for a single legacy VPS install
template. Pair it with `ovh_vps_reinstall` (legacy `/reinstall` flow) when
you already know the template id.

## Example Usage

```terraform
data "ovh_vps_template" "tpl" {
  service_name = "vps-xxxxxx.vps.ovh.net"
  template_id  = 123
}
```

## Argument Reference

* `service_name` - (Required) The VPS service_name.
* `template_id` - (Required) The template id (int).

## Attributes Reference

* `name` - Template name.
* `distribution` - Distribution string.
* `bit_format` - 32 or 64.
* `available_language` - List of supported language codes.
* `locale` - Default locale.
* `software_ids` - List of software ids installable on top of this template
  (one extra GET on `/templates/{id}/software`).

## Compatibility

This data source is part of the OVHcloud public VPS API schema, but the
underlying endpoint may not be implemented for every VPS lineup. Live
testing on a `vps-le-2-2-40` (2019v1) US instance returns
`404: Got an invalid (or empty) URL`. The endpoint is
exposed on EU and CA regions; on the US API it is only present for legacy VPS plans (pre-2019).

If you receive a 404 from this data source, your VPS plan likely does
not expose this endpoint. See [the OVHcloud VPS API console](https://api.us.ovhcloud.com/console/#/vps)
for what's available on your specific service.
