---
page_title: "Migrating VPS configurations to version 3.0.0"
---

# Migrating VPS configurations to version 3.0.0

Version 3.0.0 of the provider rewrites the `ovh_vps` and `ovh_vpss` data sources on the schema of the VPS API v2. A configuration that reads them must be updated: the first `terraform plan` after the upgrade fails on every attribute that was renamed or removed.

This guide maps each attribute of the previous version to its replacement.

## What changes

- The settings you requested for the VPS are at the root: `display_name`, `power_state`, `location`, `installation` and `options`.
- What the VPS currently runs is in `current_state`: name, power state, image, region, hardware model, IP addresses and options.
- The VPS is identified by `id`, its service name, such as `vps-1a2b3c4d.vps.ovh.net`.
- Reading a VPS that does not exist or whose service is expired fails. The previous version returned an empty `id` instead.

Data sources keep nothing in the state between two runs: once the configuration is updated, no state operation is needed.

## `ovh_vps` data source

| Previous attribute | New attribute |
|---|---|
| `service_name` (argument) | `id` (argument) |
| `id`, `name` | `id` |
| `displayname` | `display_name` for the name you requested, `current_state.display_name` for the current name |
| `state` | `current_state.power_state` (see the values below), and `power_state` for the state you requested |
| `vcore` | `current_state.model.vcpus` |
| `memory` | `current_state.model.ram`, in MB as before |
| `ips` | `current_state.ip_addresses[*].ip`, with the IP version in `current_state.ip_addresses[*].version` |
| `datacenter.name`, `datacenter.longname` | `current_state.location.region`, an OVHcloud region code such as `eu-west-gra` |
| `zone` | `current_state.location.region` |
| `netbootmode` | Removed. A VPS in rescue mode reads `current_state.power_state = "RESCUE"` |
| `model.name`, `model.offer`, `model.version` | Removed. `current_state.model` now holds the hardware of the VPS: `vcpus`, `ram` and `disk` |
| `urn` | Removed |
| `type` | Removed |
| `cluster` | Removed |
| `keymap` | Removed |
| `slamonitoring` | Removed |
| `offertype` | Removed |

The values of `state` map to `current_state.power_state` as follows:

| Previous `state` | `current_state.power_state` | `power_state` |
|---|---|---|
| `running`, `backuping` | `ACTIVE` | `ACTIVE` |
| `stopped` | `SHUTOFF` | `SHUTOFF` |
| `rescued` | `RESCUE` | `RESCUE` |
| `installing`, `rebooting`, `stopping`, `upgrading`, `maintenance` | `UNKNOWN` | null |

The new data source also exposes attributes that had no equivalent: `installation.from` and `current_state.installation.from` (the image of the VPS), `current_state.model.disk`, `options.auto_backup` and `current_state.options` (automated backup, snapshot and additional disk options).

Before:

```terraform
data "ovh_vps" "my_vps" {
  service_name = "vps-1a2b3c4d.vps.ovh.net"
}

output "vps_ips" {
  value = data.ovh_vps.my_vps.ips
}

output "vps_datacenter" {
  value = data.ovh_vps.my_vps.datacenter.name
}
```

After:

```terraform
data "ovh_vps" "my_vps" {
  id = "vps-1a2b3c4d.vps.ovh.net"
}

output "vps_ips" {
  value = data.ovh_vps.my_vps.current_state.ip_addresses[*].ip
}

output "vps_region" {
  value = data.ovh_vps.my_vps.current_state.location.region
}
```

## `ovh_vpss` data source

| Previous attribute | New attribute |
|---|---|
| `result` | `vps[*].id` |

Each element of `vps` now carries the attributes of the `ovh_vps` data source, so reading the details of every VPS no longer needs one `ovh_vps` data source per VPS.

The previous version sorted `result`. The new one lists the VPS in the order the API returns them: use `sort(data.ovh_vpss.all.vps[*].id)` to keep a sorted list. A VPS whose service is expired is left out of the list, with a warning that names it.

Before:

```terraform
data "ovh_vpss" "all" {}

output "vps_ids" {
  value = data.ovh_vpss.all.result
}
```

After:

```terraform
data "ovh_vpss" "all" {}

output "vps_ids" {
  value = sort(data.ovh_vpss.all.vps[*].id)
}
```
