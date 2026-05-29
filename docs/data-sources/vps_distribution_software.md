---
subcategory : "VPS"
---

# ovh_vps_distribution_software (Data Source)

Use this data source to list the software currently installed on one of your
VPS. The list of software ids is fetched from
`GET /vps/{serviceName}/distribution/software`, and detail for each id is then
fetched in parallel (capped at 4 concurrent requests).

## Example Usage

```terraform
data "ovh_vps_distribution_software" "installed" {
  service_name  = "vpsXXXXX.ovh.net"
  type_filter   = "webserver"
  status_filter = "stable"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.
* `type_filter` - (Optional) Restrict results to a software type. One of
  `database`, `environment`, `webserver`.
* `status_filter` - (Optional) Restrict results to a software status. One of
  `stable`, `testing`, `deprecated`.

## Attributes Reference

* `software_ids` - Sorted list of installed software ids matching the filters.
* `software` - Detail for each installed software id matching the filters:
  * `id` - Software id.
  * `name` - Software name.
  * `type` - Software type.
  * `status` - Software status.

## Compatibility

This data source is part of the OVHcloud public VPS API schema, but the
underlying endpoint may not be implemented for every VPS lineup. Live
testing on a `vps-le-2-2-40` (2019v1) US instance returns
`404: Got an invalid (or empty) URL`. The endpoint is
exposed on EU and CA regions; on the US API it is only present for legacy VPS plans (pre-2019).

If you receive a 404 from this data source, your VPS plan likely does
not expose this endpoint. See [the OVHcloud VPS API console](https://api.us.ovhcloud.com/console/#/vps)
for what's available on your specific service.
