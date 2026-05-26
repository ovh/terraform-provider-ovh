---
subcategory : "VPS"
---

# ovh_vps_template_software (Data Source)

Use this data source to retrieve details about a single software entry
attached to a legacy VPS install template.

## Example Usage

```terraform
data "ovh_vps_template_software" "sw" {
  service_name = "vps-xxxxxx.vps.ovh.net"
  template_id  = 123
  software_id  = 456
}
```

## Argument Reference

* `service_name` - (Required) The VPS service_name.
* `template_id` - (Required) The template id (int).
* `software_id` - (Required) The software id (int).

## Attributes Reference

* `name` - Software name.
* `type` - One of `database`, `environment`, `webserver`.
* `status` - One of `deprecated`, `stable`, `testing`.

## Compatibility

This data source is part of the OVHcloud public VPS API schema, but the
underlying endpoint may not be implemented for every VPS lineup. Live
testing on a `vps-le-2-2-40` (2019v1) US instance returns
`404: Got an invalid (or empty) URL`. The endpoint is
exposed on EU and CA regions; on the US API it is only present for legacy VPS plans (pre-2019).

If you receive a 404 from this data source, your VPS plan likely does
not expose this endpoint. See [the OVHcloud VPS API console](https://api.us.ovhcloud.com/console/#/vps)
for what's available on your specific service.
