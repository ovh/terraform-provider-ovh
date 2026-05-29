---
subcategory : "VPS"
---

# ovh_vps_disk_monitoring (Data Source)

Use this data source to retrieve time-series monitoring statistics for a VPS
disk over a given period.

## Example Usage

```terraform
data "ovh_vps_disk_monitoring" "mon" {
  service_name = "vps-XXXXXX.vps.ovh.net"
  disk_id      = 1234
  period       = "lastday"
  type         = "cpu:used"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of the VPS (e.g.
  `vps-XXXXXX.vps.ovh.net`).
* `disk_id` - (Required) Numeric identifier of the disk attached to the VPS.
* `period` - (Required) Time window for the monitoring data. One of `lastday`,
  `lastweek`, `lastmonth`, `lastyear`, `today`.
* `type` - (Required) Statistic type to query. Accepted values mirror
  `vps.VpsStatisticTypeEnum`: `cpu:iowait`, `cpu:max`, `cpu:nice`, `cpu:sys`,
  `cpu:used`, `cpu:user`, `mem:max`, `mem:used`, `net:rx`, `net:tx`. Only the
  disk-relevant subset is meaningful on this endpoint; the API will reject
  incompatible combinations.

## Attributes Reference

* `unit` - The unit of the returned series (e.g. `%`, `B`, `B/s`).
* `values` - Ordered list of samples in the requested period. Each entry has:
  * `timestamp` - RFC3339 timestamp of the sample.
  * `value` - Numeric value of the sample.
