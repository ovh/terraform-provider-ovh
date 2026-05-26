---
subcategory : "VPS"
---

# ovh_vps_backup_ftp_authorizable_blocks (Data Source)

Use this data source to list CIDR blocks that can be granted backup FTP
access on a VPS.

## Example Usage

```terraform
data "ovh_vps_backup_ftp_authorizable_blocks" "blocks" {
  service_name = "vpsXXXXXX.ovh.net"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your VPS.

## Attributes Reference

* `id` - Set to the `service_name`.
* `blocks` - List of CIDR blocks authorized to be granted backup FTP access.

## Compatibility

This data source wraps `GET /vps/{serviceName}/backupftp/authorizableBlocks`. Live cross-region probing on 2026-05-26 shows
the endpoint is present in the **EU** and **CA** API schemas (`eu.api.ovh.com`,
`ca.api.ovh.com`) but **NOT** in the **US** schema (`api.us.ovhcloud.com`).

On a US-region VPS the OVHcloud API returns
`404: Got an invalid (or empty) URL`. Use this data source on EU or CA accounts,
or wait for OVH to expose this endpoint on US.
