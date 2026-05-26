---
subcategory : "VPS"
---

# ovh_vps_datacenters (Data Source)

Use this data source to list the OVHcloud VPS datacenters, optionally filtered by country.

## Example Usage

```terraform
data "ovh_vps_datacenters" "all" {}

data "ovh_vps_datacenters" "us_only" {
  country = "US"
}
```

## Argument Reference

* `country` - (Optional) Two-letter country code used to filter datacenters (e.g. `FR`, `US`, `CA`, `DE`).

## Attributes Reference

* `datacenters` - The list of datacenter codes (e.g. `gra`, `sbg`, `bhs`).

## Compatibility

This data source is part of the OVHcloud public VPS API schema, but the
underlying endpoint may not be implemented for every VPS lineup. Live
testing on a `vps-le-2-2-40` (2019v1) US instance returns
`404: Got an invalid (or empty) URL`. The endpoint is
exposed on EU and CA regions but not on the US API (as of 2026-05-26).

If you receive a 404 from this data source, your VPS plan likely does
not expose this endpoint. See [the OVHcloud VPS API console](https://api.us.ovhcloud.com/console/#/vps)
for what's available on your specific service.
