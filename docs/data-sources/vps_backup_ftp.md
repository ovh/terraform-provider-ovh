---
subcategory : "VPS"
---

# ovh_vps_backup_ftp (Data Source)

Use this data source to retrieve information about the backup FTP storage
attached to a VPS.

## Example Usage

```terraform
data "ovh_vps_backup_ftp" "backup" {
  service_name = "vpsXXXXXX.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.

## Attributes Reference

* `id` - Set to the `service_name`.
* `ftp_backup_name` - Name of the backup FTP server.
* `type` - Type of backup FTP offer.
* `read_only_date` - Date at which the backup FTP will be set in read-only
  mode (RFC3339), empty if not applicable.
* `quota` - Map describing the storage quota:
  * `quota.unit` - Unit of the quota value (e.g. `GB`).
  * `quota.value` - Numeric quota value.
* `usage` - Map describing the current storage usage:
  * `usage.unit` - Unit of the usage value (e.g. `GB`).
  * `usage.value` - Numeric usage value.

## Compatibility

This data source wraps `GET /vps/{serviceName}/backupftp`. Live cross-region probing on 2026-05-26 shows
the endpoint is present in the **EU** and **CA** API schemas (`eu.api.ovh.com`,
`ca.api.ovh.com`) but **NOT** in the **US** schema (`api.us.ovhcloud.com`).

On a US-region VPS the OVHcloud API returns
`404: Got an invalid (or empty) URL`. Use this data source on EU or CA accounts,
or wait for OVH to expose this endpoint on US.
