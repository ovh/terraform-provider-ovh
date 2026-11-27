---
page_title: "Migrating VPS configurations to version 3.0.0"
---

# Migrating VPS configurations to version 3.0.0

Version 3.0.0 of the provider rewrites the `ovh_vps` resource and the `ovh_vps` and `ovh_vpss` data sources on the schema of the VPS API v2. A configuration that uses them must be updated: the first `terraform plan` after the upgrade fails on every attribute that was renamed or removed.

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

Reading the region calls `GET /v2/location`, which the previous version did not call. When your credentials are not allowed to call it, for instance a consumer key whose access rules only allow `GET /vps/*`, `location.region` and `current_state.location.region` are null and the read emits a warning that gives the reason.

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

`vps` is sorted by `id`, as `result` was. A VPS whose service is expired is left out of the list, with a warning that names it.

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
  value = data.ovh_vpss.all.vps[*].id
}
```

## `ovh_vps` resource

The state of an `ovh_vps` resource is migrated automatically by the first `terraform plan` after the upgrade: the VPS keeps its `id`, its service name, and does not need to be imported again. The migration keeps `id`, `plan_code`, `order.duration`, `order.pricing_mode`, `order.order_id`, `order.options`, `display_name`, `installation.ssh_key` and `installation.send_password_by_email`; the refresh that follows reads the other attributes from the VPS.

| Previous attribute | New attribute |
|---|---|
| `service_name`, `id`, `name` | `id`. The import identifier is unchanged: the service name |
| `plan[0].plan_code` | `plan_code` |
| `plan[0].configuration` `vps_datacenter` (or `zone`, `zone_type`) | `location.region`, an OVHcloud region code such as `eu-west-gra` |
| `plan[0].configuration` `vps_os` | `installation.from.image_name` |
| `image_id` | `installation.from.image_id` |
| `plan[0].duration`, `plan[0].pricing_mode` | `order.duration`, `order.pricing_mode` |
| `plan[0].item_id`, `plan[0].quantity` | Removed |
| `ovh_subsidiary` | Removed: the subsidiary of your account is used |
| `plan_option` | `order.options`, the plan codes of the options: `plan_option[*].plan_code`. Each option is ordered once, with the duration and pricing mode of the VPS. The options of the VPS are read in `current_state.options` |
| `order.order_id` | `order.order_id` |
| `order.date`, `order.expiration_date`, `order.details` | Removed |
| `display_name` | `display_name` for the name you requested, `current_state.display_name` for the current name |
| `public_ssh_key` | `installation.ssh_key` |
| `do_not_send_password` | `installation.send_password_by_email`, inverted: `do_not_send_password = true` becomes `send_password_by_email = false` |
| `state` | `current_state.power_state`, and `power_state` for the state you requested, with the values of the `ovh_vps` data source above |
| `netboot_mode` | Removed. A VPS in rescue mode reads `power_state = "RESCUE"` |
| `vcore`, `model.vcore` | `current_state.model.vcpus` |
| `memory_limit`, `model.memory` | `current_state.model.ram` |
| `model.disk` | `current_state.model.disk` |
| `zone` | `current_state.location.region` |
| `model.name`, `model.offer`, `model.version`, `model.available_options`, `model.datacenter`, `model.maximum_additionnal_ip` | Removed. `plan_code` holds the plan of the VPS |
| `iam` | Removed |
| `cluster`, `keymap`, `sla_monitoring`, `monitoring_ip_blocks`, `offer_type` | Removed |

Before:

```terraform
resource "ovh_vps" "web" {
  display_name         = "web"
  ovh_subsidiary       = data.ovh_me.my_account.ovh_subsidiary
  public_ssh_key       = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample user@laptop"
  do_not_send_password = true

  plan = [
    {
      duration     = "P1M"
      plan_code    = "vps-2025-model1"
      pricing_mode = "default"

      configuration = [
        {
          label = "vps_datacenter"
          value = "GRA"
        },
        {
          label = "vps_os"
          value = "Debian 12"
        }
      ]
    }
  ]
}
```

After:

```terraform
resource "ovh_vps" "web" {
  plan_code    = "vps-2025-model1"
  display_name = "web"

  location = {
    region = "eu-west-gra"
  }

  installation = {
    from = {
      image_name = "Debian 12"
    }
    ssh_key                = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample user@laptop"
    send_password_by_email = false
  }
}
```

### Check the first plan

Check the first `terraform plan` after the migration, and after an import:

- When `plan_code` or `location.region` in your configuration differ from the ones of the VPS, Terraform plans to replace it, which terminates the VPS and orders a new one. This happens when the VPS changed plan since its order, or when the region you wrote is not the one of the VPS. Set them to the values of the state, shown by `terraform state show`.
- Any other difference, such as `display_name`, `power_state`, `options.auto_backup.cron`, `installation.from` or `installation.ssh_key`, is shown as an update in place, and `terraform apply` applies the change.
- A difference of `installation.from`, `installation.ssh_key` or `installation.post_installation_script` reinstalls the VPS, which erases its disk: `terraform plan` warns `This apply reinstalls VPS …, which erases its disk`. `installation.from.image_name` is the image the VPS runs: to keep the VPS as it is, set it to the value of the state.
- Changes of `timeouts`, `order` and `installation.send_password_by_email` are only recorded in the state. So are `installation.ssh_key` and `installation.post_installation_script` when the state has none, since they are never read back, and `plan_code`, `location.region`, `power_state` and `installation.from` when they could not be read from the VPS.

A plan with no change means the configuration matches the VPS.

### After a delivery timeout

`terraform apply` waits for the delivery of the VPS for at most `timeouts.create` (60 minutes by default). When the delivery takes longer, the order is paid and kept in the state, in `order.order_id`, and Terraform marks the resource as tainted: the next `apply` would terminate it and order a new VPS. Run `terraform untaint` on it:

```sh
terraform untaint ovh_vps.web
```

The same goes when `terraform apply` is interrupted while it waits, or fails on any other error after the payment. The next `terraform plan` or `terraform apply` checks the delivery once, without waiting, and fills `id` once the VPS is delivered. When the order is cancelled, they fail with `VPS order … is cancelled`: remove the resource from the state with `terraform state rm`.
