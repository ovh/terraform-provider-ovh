---
subcategory : "VPS"
---

# ovh_vps_distribution_software_item (Data Source)

Use this data source to retrieve details about a single piece of software
installed on a VPS. The list of installed software ids is exposed by the
`ovh_vps_distribution_software` data source.

## Example Usage

```terraform
data "ovh_vps_distribution_software_item" "nginx" {
  service_name = "vpsXXXXX.ovh.net"
  software_id  = 42
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.
* `software_id` - (Required) Software id as exposed by
  `/vps/{serviceName}/distribution/software`.

## Attributes Reference

* `name` - Software name.
* `type` - Software type (`database`, `environment`, `webserver`).
* `status` - Software status (`stable`, `testing`, `deprecated`).

## Compatibility

This data source is part of the OVHcloud public VPS API schema, but the
underlying endpoint may not be implemented for every VPS lineup. Live
testing on a `vps-le-2-2-40` (2019v1) US instance returns
`404: Got an invalid (or empty) URL`. The endpoint is
exposed on EU and CA regions; on the US API it is only present for legacy VPS plans (pre-2019).

If you receive a 404 from this data source, your VPS plan likely does
not expose this endpoint. See [the OVHcloud VPS API console](https://api.us.ovhcloud.com/console/#/vps)
for what's available on your specific service.
