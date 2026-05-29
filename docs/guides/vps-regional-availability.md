---
subcategory : "VPS"
page_title: "VPS — Regional API surface differences"
---

# VPS — Regional API surface differences

The OVHcloud VPS API is exposed on three regions: **EU** (`eu.api.ovh.com`),
**CA** (`ca.api.ovh.com`), and **US** (`api.us.ovhcloud.com`). The Terraform
provider talks to the region selected by `OVH_ENDPOINT` (or the `endpoint`
provider attribute).

Verified live on 2026-05-26 by fetching the `vps.json` schema from each
region's API host:

| Region | Endpoints in schema |
|---|---|
| EU | 74 |
| CA | 74 (identical to EU) |
| US | 46 (28 endpoints absent) |

## What's missing on US

The following 28 endpoints are present in the EU/CA schemas but absent on US.
Resources and data sources that wrap these endpoints will return
`404: Got an invalid (or empty) URL` if used against a US VPS:

```
/vps/datacenter                                      (top-level catalog)
/vps/{serviceName}/availableUpgrade
/vps/{serviceName}/backupftp                         (and 4 sub-paths)
/vps/{serviceName}/changeContact
/vps/{serviceName}/distribution                      (and 2 sub-paths)
/vps/{serviceName}/migration2016                     (deprecated)
/vps/{serviceName}/migration2018                     (deprecated)
/vps/{serviceName}/models
/vps/{serviceName}/openConsoleAccess                 (used by ovh_vps_vnc)
/vps/{serviceName}/reinstall                         (used by ovh_vps_reinstall)
/vps/{serviceName}/setPassword                       (used by ovh_vps_set_password)
/vps/{serviceName}/status
/vps/{serviceName}/templates                         (and 3 sub-paths)
/vps/{serviceName}/use                               (deprecated)
/vps/{serviceName}/veeam                             (and 4 sub-paths)
```

## Resources affected on US

* `ovh_vps_set_password`, `ovh_vps_reinstall`, `ovh_vps_vnc`,
  `ovh_vps_change_contact`
* `ovh_vps_backup_ftp_access`, `ovh_vps_backup_ftp_password`
* `ovh_vps_veeam_restore`

## Data sources affected on US

* `ovh_vps_models`, `ovh_vps_available_upgrade`, `ovh_vps_status`,
  `ovh_vps_datacenters`, `ovh_vps_current_image`
* `ovh_vps_distribution`, `ovh_vps_distribution_software`,
  `ovh_vps_distribution_software_item`
* `ovh_vps_template`, `ovh_vps_templates`, `ovh_vps_template_software`
* `ovh_vps_backup_ftp`, `ovh_vps_backup_ftp_access`,
  `ovh_vps_backup_ftp_authorizable_blocks`
* `ovh_vps_veeam`, `ovh_vps_veeam_restore_point`,
  `ovh_vps_veeam_restore_points`, `ovh_vps_veeam_restored_backup`

## Resources / data sources that work on every region

All other VPS resources and data sources in this provider — the core
`ovh_vps` lifecycle, snapshots, automated backup, additional disks, IPs
and reverse DNS, secondary DNS, tasks, options, service info, migration
2020→2025, available images, order-rule discovery, and the vRack-VPS
attachment — wrap endpoints present on all three regional schemas and
work uniformly.

## What we did about the difference

* Each affected data source's documentation page carries a `## Compatibility`
  section calling out the regional limitation explicitly.
* Each affected resource's documentation page does the same.
* Acceptance tests for these resources use the `skipIfEndpointMissing` helper
  in `ovh/vps_test_helpers_test.go`, which probes the endpoint at PreCheck
  time and skips the test cleanly with an informative message when run on
  a region where the endpoint isn't exposed (rather than failing with a
  raw `404`).
