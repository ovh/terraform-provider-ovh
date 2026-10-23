---
page_title: "Migrating Public Cloud resources from ovh_cloud_project_* to ovh_cloud_*"
---

# Migrating Public Cloud resources from `ovh_cloud_project_*` to `ovh_cloud_*`

The provider ships a new family of Public Cloud resources and data sources, named `ovh_cloud_*`, built on the Public Cloud API v2. They replace the legacy `ovh_cloud_project_*` family, built on the API v1, for instances, block, file and object storage, private networks, gateways, load balancers and SSH keys.

The legacy resources keep working but are deprecated: `terraform plan` prints a warning for each of them. This guide explains how to move an existing infrastructure to the new resources **without recreating anything**: the move is a Terraform state operation, nothing is touched on the OpenStack side.

## What changes

| Legacy resource | New resource |
|---|---|
| `ovh_cloud_project_instance` | `ovh_cloud_instance` |
| `ovh_cloud_project_instance_snapshot` | `ovh_cloud_instance_snapshot` |
| `ovh_cloud_project_volume` | `ovh_cloud_storage_block_volume` |
| `ovh_cloud_project_volume_backup` | `ovh_cloud_storage_block_volume_backup` |
| `ovh_cloud_project_file_storage_share` | `ovh_cloud_storage_file_share` |
| `ovh_cloud_project_file_storage_share_network` | `ovh_cloud_storage_file_share_network` |
| `ovh_cloud_project_storage` | `ovh_cloud_storage_object_bucket` |
| `ovh_cloud_project_network_private` | `ovh_cloud_network_private_vrack` (one per region) |
| `ovh_cloud_project_network_private_subnet`, `ovh_cloud_project_network_private_subnet_v2` | `ovh_cloud_network_private_vrack_subnet` |
| `ovh_cloud_project_region_network` | `ovh_cloud_network_private_vrack` + `ovh_cloud_network_private_vrack_subnet` |
| `ovh_cloud_project_gateway`, `ovh_cloud_project_gateway_interface` | `ovh_cloud_gateway` |
| `ovh_cloud_project_loadbalancer` | `ovh_cloud_loadbalancer` + `ovh_cloud_loadbalancer_listener` + `ovh_cloud_loadbalancer_pool` + `ovh_cloud_loadbalancer_pool_member` |
| `ovh_cloud_project_ssh_key` | `ovh_cloud_ssh_key` |
| `ovh_cloud_project_failover_ip_attach` | the `networks` attribute of `ovh_cloud_instance` |

Data sources follow the same renaming; see [Data sources](#data-sources) at the end.

What does **not** change:

- `service_name` is still the Public Cloud project ID, and still falls back to the `OVH_CLOUD_PROJECT_SERVICE` environment variable on most resources.
- `region` is still a flat string, for example `"GRA11"`. The `region_name` argument of some legacy resources becomes `region`.
- The resource `id` is still the OpenStack identifier of the object, so the import IDs below can be read from the legacy state with `terraform state show`.

What changes on every new resource:

- **Enum values are upper snake case:** `"classic"` becomes `"CLASSIC"`, `"roundRobin"` becomes `"ROUND_ROBIN"`, `"enabled"` becomes `"ENABLED"`.
- **The arguments describe the target state, the `current_state` attribute holds what OpenStack reports.** Values you used to read at the root (`addresses`, `status`, `vip_address`, `virtual_host`...) now live under `current_state`, next to `resource_status`, `created_at`, `updated_at` and `checksum`.
- **Most arguments update in place.** The legacy resources recreated the object on almost every change; the new ones resize, rename, retype, attach and detach without recreation. Each resource page lists the arguments that still force a replacement.
- **One resource per concern.** The legacy instance could create its network, gateway, floating IP and SSH key inline; the legacy load balancer embedded its listeners, pools and members. The new resources reference each other by ID instead.
- **Import IDs no longer carry the region:** `<service_name>/<id>`, and `<service_name>/<parent_id>/<id>` for child resources.

## Before you start

- Terraform **1.7 or later** for the `import` and `removed` blocks used below. With Terraform 1.5 or 1.6 use the `terraform import` and `terraform state rm` commands shown at the end of the procedure.
- Upgrade the provider to a version that ships every resource you need; check the [changelog](https://github.com/ovh/terraform-provider-ovh/blob/master/CHANGELOG.md).
- Back up your state: `terraform state pull > backup.tfstate`.
- Migrate one resource at a time and run `terraform plan` after each step. A successful migration ends with a plan that shows **no changes**.

!> **Never run `terraform destroy`, and never delete a legacy block without a `removed` block or a `terraform state rm`.** Removing a legacy resource block from the configuration schedules the destruction of the real instance, volume or bucket.

-> Terraform can draft the new block for you: declare the `import` block without the matching `resource` block, then run `terraform plan -generate-config-out=generated.tf`. Review the generated file (it contains every computed attribute) and keep only the arguments you want to manage.

## The procedure

Every migration follows the same five steps. The example moves an instance; the per-resource sections below only give the mapping and the import ID.

**1. Write the new resource block** next to the legacy one, with the same values translated through the mapping table of the resource.

```terraform
# the existing and deprecated resource
resource "ovh_cloud_project_instance" "web" {
  service_name   = "<public cloud project ID>"
  region         = "GRA11"
  name           = "web-1"
  billing_period = "hourly"

  boot_from {
    image_id = "<image ID>"
  }
  flavor {
    flavor_id = "<flavor ID>"
  }
  ssh_key {
    name = "deploy"
  }
  network {
    public = true
  }
}

# the new resource type (API v2)
resource "ovh_cloud_instance" "web" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  name         = "web-1"
  flavor_id    = "<flavor ID>"
  image_id     = "<image ID>"
  ssh_key_name = "deploy"

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

**2. Replace the legacy block with an `import` block and a `removed` block.** The `import` block hands the existing object to the new resource; the `removed` block tells Terraform to forget the legacy address without destroying anything. The import ID is `<service_name>/<id>`; `terraform state show ovh_cloud_project_instance.web` prints the `id`.

```terraform
import {
  to = ovh_cloud_instance.web
  id = "<public cloud project ID>/<instance ID>"
}

removed {
  from = ovh_cloud_project_instance.web

  lifecycle {
    destroy = false
  }
}
```

Delete the legacy `resource "ovh_cloud_project_instance" "web"` block in the same edit. Terraform refuses a `removed` block whose address is still declared (`Error: Removed resource still exists`), so this cannot wait for the cleanup step.

**3. Plan and check.** The plan must announce one import and one forgotten object, and **no update or replacement**:

```sh
$ terraform plan
ovh_cloud_instance.web: Preparing import... [id=<public cloud project ID>/<instance ID>]
ovh_cloud_instance.web: Refreshing state...

Terraform will perform the following actions:

  # ovh_cloud_instance.web will be imported
    resource "ovh_cloud_instance" "web" {
        ...
    }

  # ovh_cloud_project_instance.web will no longer be managed by Terraform, but will not be destroyed
  # (destroy = false is set in the configuration)
  . resource "ovh_cloud_project_instance" "web" {
        ...
    }

Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.

Warning: Some objects will no longer be managed by Terraform

If you apply this plan, Terraform will discard its tracking information for
the following objects, but it will not delete them:
 - ovh_cloud_project_instance.web
```

If the plan shows an update or a replacement, an argument of the new block does not match the live object. Fix the block before applying: the mapping tables below list the usual culprits, and `terraform state show ovh_cloud_instance.web` after a `terraform plan -generate-config-out` gives the live values.

**4. Apply.** Nothing is created, changed or destroyed: the apply only writes the state.

```sh
$ terraform apply
...
ovh_cloud_instance.web: Importing... [id=<public cloud project ID>/<instance ID>]
ovh_cloud_instance.web: Import complete [id=<instance ID>]

Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.
```

**5. Clean up.** The `import` and `removed` blocks are one-shot instructions: delete them. Update every reference to the legacy resource (`ovh_cloud_project_instance.web.id` becomes `ovh_cloud_instance.web.id`, computed attributes move under `current_state`, see the mapping tables). Run `terraform plan` once more: it must show no changes and no deprecation warning for the migrated resource.

### Without `import` and `removed` blocks (Terraform 1.5 and 1.6)

The two blocks of step 2 become two commands, run after deleting the legacy block from the configuration. There is no plan preview: `terraform import` writes the state as soon as it runs, so the check of step 3 happens after the import instead of before.

```sh
# step 2: delete the legacy block from the configuration, then
terraform state rm ovh_cloud_project_instance.web
terraform import ovh_cloud_instance.web "<public cloud project ID>/<instance ID>"
# steps 3 and 4 are already done: check that nothing is pending
terraform plan
```

Steps 1 and 5 do not change, except that there are no `import` and `removed` blocks to delete.

## Instances

`ovh_cloud_project_instance` → `ovh_cloud_instance`. Import ID: `<service_name>/<instance_id>`.

| Legacy | New |
|---|---|
| `billing_period` | removed; the new resource bills hourly |
| `flavor { flavor_id }` | `flavor_id`; a change resizes in place |
| `boot_from { image_id }` | `image_id`; a change rebuilds the instance and wipes its root disk |
| `boot_from { volume_id }` | `volume_ids = ["<volume ID>"]` with no `image_id`; the first bootable volume is the root disk |
| `ssh_key { name }` | `ssh_key_name` |
| `ssh_key_create { name, public_key }` | an `ovh_cloud_ssh_key` resource referenced through `ssh_key_name` |
| `user_data` (plain text) | `user_data = base64encode(...)`, sensitive and write-only |
| `auto_backup { cron, rotation }` | removed |
| `group { group_id }` | `group_id` |
| `network { public = true }` | `networks = [{ auto_assign_public_ip = true }]` |
| `network { private { network { id, subnet_id } ip } }` | `networks = [{ network_id, subnet_id, ip }]`, one entry per interface |
| `network { private { floating_ip { id } } }` | `networks = [{ network_id, subnet_id, ip = "<floating IP address>" }]`; the floating IP is referenced by address |
| `network { private { floating_ip_create } }` | an `ovh_cloud_floating_ip` resource, its `id` is the address to put in `networks[*].ip` |
| `network { private { gateway, gateway_create } }` | an `ovh_cloud_gateway` resource with the subnet in `subnet_ids` |
| `network { private { network_create } }` | `ovh_cloud_network_private_vrack` and `ovh_cloud_network_private_vrack_subnet` resources |
| `timeouts {}` | removed |
| `addresses[*].{ip, version}` | `current_state.networks[*].addresses[*].{ip, version, mac, type}` |
| `attached_volumes[*].id` | `current_state.volumes[*].{id, name, size}` |
| `flavor_name` | `current_state.flavor.name` |
| `status` | `current_state.power_state` (OpenStack) and `resource_status` (API) |

New arguments with no legacy counterpart: `power_state` (`ACTIVE`, `SHUTOFF`, `SHELVED`), `security_group_ids`, `shares`, `volume_ids` for data volumes.

~> **`user_data` after an import.** The import does not read the existing cloud-init script, so the imported state has no `user_data`. If the new block sets `user_data`, the first `terraform apply` reinstalls the instance and wipes its root disk. Leave `user_data` out of the block until the migration is done, or set it only when a reinstall is acceptable.

~> **Security groups.** Leaving `security_group_ids` out keeps the project's `default` security group. Setting it to `[]` removes every group and blocks all inbound traffic.

Before, an instance on a private network with a pinned IP and a public interface:

```terraform
resource "ovh_cloud_project_instance" "app" {
  service_name   = "<public cloud project ID>"
  region         = "GRA11"
  name           = "app-1"
  billing_period = "hourly"

  boot_from {
    image_id = "<image ID>"
  }
  flavor {
    flavor_id = "<flavor ID>"
  }
  ssh_key {
    name = "deploy"
  }
  network {
    public = true
    private {
      network {
        id        = ovh_cloud_project_network_private.net.regions_attributes[0].openstackid
        subnet_id = ovh_cloud_project_network_private_subnet.subnet.id
      }
      ip = "10.0.0.10"
    }
  }
}
```

After:

```terraform
resource "ovh_cloud_instance" "app" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  name         = "app-1"
  flavor_id    = "<flavor ID>"
  image_id     = "<image ID>"
  ssh_key_name = "deploy"

  networks = [
    { auto_assign_public_ip = true },
    {
      network_id = ovh_cloud_network_private_vrack.net.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.subnet.id
      ip         = "10.0.0.10"
    },
  ]
}
```

## Instance snapshots

`ovh_cloud_project_instance_snapshot` → `ovh_cloud_instance_snapshot`. Import ID: `<service_name>/<snapshot_id>`.

| Legacy | New |
|---|---|
| (computed `region`) | `region`, now a **required** argument |
| `instance_id`, `name` | unchanged |
| `creation_date` | `created_at` |
| `status`, `min_disk`, `min_ram`, `visibility` | `current_state.{status, min_disk, min_ram, visibility}` |
| `size` in GiB | `current_state.size` in bytes |
| `flavor_type`, `plan_code`, `tags`, `type`, `user` | removed |

```terraform
resource "ovh_cloud_project_instance_snapshot" "snap" {
  service_name = "<public cloud project ID>"
  instance_id  = ovh_cloud_project_instance.app.id
  name         = "app-1-before-upgrade"
}

resource "ovh_cloud_instance_snapshot" "snap" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  instance_id  = ovh_cloud_instance.app.id
  name         = "app-1-before-upgrade"
}
```

## Block storage volumes

`ovh_cloud_project_volume` → `ovh_cloud_storage_block_volume`. Import ID: `<service_name>/<volume_id>`.

| Legacy | New |
|---|---|
| `region_name` | `region` |
| `name`, `size` | unchanged, both **required** now |
| `description` | removed |
| `type = "classic"` | `volume_type = "CLASSIC"`; a change retypes in place |
| `type = "high-speed"` | `volume_type = "HIGH_SPEED"` |
| `type = "high-speed-gen2"` | `volume_type = "HIGH_SPEED_GEN2"` |
| `type = "<type>-luks"` | `volume_type = "<TYPE>"` plus `encryption = { enabled = true }` |
| `type = "classic-multiattach"` | `volume_type = "CLASSIC"`; in a 3-AZ region the API picks the multi-attach variant itself |
| `encryption { encrypted, kms { domain_id, service_key_id } }` | `encryption = { enabled, kms = { domain_id, service_key_id } }` |
| `image_id`, `snapshot_id` | `create_from = { image_id }`, `create_from = { snapshot_id }`; `create_from = { backup_id }` restores a backup |
| attachment to an instance | the `volume_ids` list of `ovh_cloud_instance` |
| `region_name` (computed) | `current_state.location.region` |

```terraform
resource "ovh_cloud_project_volume" "data" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA11"
  name         = "app-data"
  size         = 50
  type         = "high-speed"
}

resource "ovh_cloud_storage_block_volume" "data" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  name         = "app-data"
  size         = 50
  volume_type  = "HIGH_SPEED"
}
```

~> A volume attached to an instance stays attached through the migration. To manage the attachment, add the volume to the instance's `volume_ids` in the same change; the plan must show no update on the instance, since the volume is already attached.

## Block storage backups

`ovh_cloud_project_volume_backup` → `ovh_cloud_storage_block_volume_backup`. Import ID: `<service_name>/<backup_id>`.

| Legacy | New |
|---|---|
| `region_name` | `region` |
| `volume_id` | unchanged |
| `name` | unchanged, **required** now; `description` is new |
| `creation_date` | `created_at` |
| `size`, `region` | `current_state.size`, `current_state.location.region` |
| `status` | `resource_status` |

```terraform
resource "ovh_cloud_project_volume_backup" "data" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA11"
  volume_id    = ovh_cloud_project_volume.data.id
  name         = "app-data-weekly"
}

resource "ovh_cloud_storage_block_volume_backup" "data" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  volume_id    = ovh_cloud_storage_block_volume.data.id
  name         = "app-data-weekly"
}
```

## File storage shares

`ovh_cloud_project_file_storage_share` → `ovh_cloud_storage_file_share`, and `ovh_cloud_project_file_storage_share_network` → `ovh_cloud_storage_file_share_network`. Import IDs: `<service_name>/<share_id>` and `<service_name>/<share_network_id>`.

The legacy share accepted `network_id` and `subnet_id` inline and created a share network behind the scenes. The new share only takes a `share_network_id`: migrate the share network first (or import the one the legacy share created, whose ID is the legacy computed `share_network_id`), then the share.

Share network:

| Legacy | New |
|---|---|
| `region_name` | `region` |
| `name` | unchanged, **required** now |
| `network_id`, `subnet_id`, `description` | unchanged |
| `availability_zone` | unchanged |
| `cidr`, `network_type` | removed |

Share:

| Legacy | New |
|---|---|
| `region_name` | `region` |
| `name` | unchanged, **required** now |
| (computed `protocol`) | `protocol = "NFS"`, now a **required** argument |
| `type = "standard-1az"` | `share_type = "STANDARD_1AZ"` |
| `network_id`, `subnet_id` | removed; use `share_network_id` |
| `share_network_id` | unchanged |
| `snapshot_id` | `create_from = { snapshot_id }` |
| `availability_zone` | unchanged |
| `size` in GB | unchanged, resizes in place |
| `is_public`, `status` | removed; see `resource_status` |
| (export path) | `current_state.export_locations[*].path` |

```terraform
resource "ovh_cloud_project_file_storage_share" "share" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA11"
  name         = "app-shared"
  size         = 150
  type         = "standard-1az"
  network_id   = "<network ID>"
  subnet_id    = "<subnet ID>"
}

resource "ovh_cloud_storage_file_share_network" "share" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  name         = "app-shared"
  network_id   = "<network ID>"
  subnet_id    = "<subnet ID>"
}

resource "ovh_cloud_storage_file_share" "share" {
  service_name     = "<public cloud project ID>"
  region           = "GRA11"
  name             = "app-shared"
  size             = 150
  protocol         = "NFS"
  share_type       = "STANDARD_1AZ"
  share_network_id = ovh_cloud_storage_file_share_network.share.id
}
```

Access rules and snapshots have their own resources, `ovh_cloud_storage_file_share_acl` and `ovh_cloud_storage_file_share_snapshot`; an instance can also mount a share through its `shares` attribute.

## Object storage buckets

`ovh_cloud_project_storage` → `ovh_cloud_storage_object_bucket`. Import ID: `<service_name>/<bucket_id>`, where the bucket ID is the bucket name on a single-region deployment and `<REGION>_<name>` (for example `GRA_my-bucket`) on a multi-region one.

| Legacy | New |
|---|---|
| `region_name` | `region`, upper case |
| `name`, `tags` | unchanged |
| `owner_id` (number) | `owner_user_id` (string) |
| `encryption { sse_algorithm }` | `encryption = { algorithm }`, `"AES256"` or `"PLAINTEXT"` |
| `versioning { status = "enabled" }` | `versioning = { status = "ENABLED" }` |
| `object_lock { status = "enabled", rule { mode = "governance", period = "P30D" } }` | `object_lock = { mode = "GOVERNANCE", retention_days = 30 }`; the block's presence replaces `status` |
| `replication { rules [...] }` | not supported |
| `hide_objects`, `limit`, `marker`, `prefix`, `objects[*]` | removed; the resource no longer lists objects |
| `virtual_host`, `objects_count`, `objects_size` | `current_state.{virtual_host, objects_count, objects_size}` |

```terraform
resource "ovh_cloud_project_storage" "assets" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA"
  name         = "app-assets"
  versioning = {
    status = "enabled"
  }
}

resource "ovh_cloud_storage_object_bucket" "assets" {
  service_name = "<public cloud project ID>"
  region       = "GRA"
  name         = "app-assets"
  versioning = {
    status = "ENABLED"
  }
}
```

~> **Object lock recreates the bucket.** Adding, changing or removing the `object_lock` block on the new resource replaces the bucket. Translate the legacy block exactly, including the retention period converted to days, and check that the plan after import shows no change before applying.

~> Lifecycle configurations (`ovh_cloud_project_storage_object_bucket_lifecycle_configuration`) and replication jobs keep their legacy resources. They reference the bucket by name, so they keep working after the bucket itself is migrated.

## Private networks

`ovh_cloud_project_network_private` → `ovh_cloud_network_private_vrack`. Import ID: `<service_name>/<network_id>`, where the network ID is the **regional OpenStack UUID** (the legacy `regions_attributes[*].openstackid`), not the legacy `pn-xxxx` ID.

The legacy resource spanned several regions. The new one is regional: **one resource per region**, sharing the same `vlan_id`.

| Legacy | New |
|---|---|
| `regions = ["GRA11", "BHS5"]` | one resource per region, `region = "GRA11"`, `region = "BHS5"` |
| `vlan_id` (default `0`) | `vlan_id`, **set it explicitly** on every regional resource: when omitted the API assigns one; a change recreates the network |
| `name` | unchanged, updates in place; `description` is new |
| `id` = `pn-xxxx` | `id` = regional OpenStack UUID |
| `regions_attributes[*].openstackid`, `regions_openstack_ids` | the `id` of each regional resource |
| `regions_attributes[*].status`, `regions_status`, `status` | `resource_status` |
| `type` | removed |

```terraform
resource "ovh_cloud_project_network_private" "net" {
  service_name = "<public cloud project ID>"
  name         = "backend"
  vlan_id      = 42
  regions      = ["GRA11", "BHS5"]
}

resource "ovh_cloud_network_private_vrack" "net_gra" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  name         = "backend"
  vlan_id      = 42
}

resource "ovh_cloud_network_private_vrack" "net_bhs" {
  service_name = "<public cloud project ID>"
  region       = "BHS5"
  name         = "backend"
  vlan_id      = 42
}
```

```terraform
import {
  to = ovh_cloud_network_private_vrack.net_gra
  id = "<public cloud project ID>/<OpenStack network ID in GRA11>"
}

import {
  to = ovh_cloud_network_private_vrack.net_bhs
  id = "<public cloud project ID>/<OpenStack network ID in BHS5>"
}

removed {
  from = ovh_cloud_project_network_private.net

  lifecycle {
    destroy = false
  }
}
```

~> The legacy resource with `vlan_id = 0` (the default) created a network without a vRack VLAN; translate it as `vlan_id = 0`. `vlan_id` is not supported in Local Zone regions: leave it out there.

### Subnets

`ovh_cloud_project_network_private_subnet` and `ovh_cloud_project_network_private_subnet_v2` → `ovh_cloud_network_private_vrack_subnet`. Import ID: `<service_name>/<network_id>/<subnet_id>` with the regional network UUID.

| Legacy v1 | Legacy v2 | New |
|---|---|---|
| `network_id` = `pn-xxxx` | `network_id` = OpenStack UUID | `network_id = ovh_cloud_network_private_vrack.<name>.id` |
| `region` | `region` | `region` |
| `network` | `cidr` | `cidr` |
| `start`, `end` | `allocation_pools` | `allocation_pools = [{ start, end }]`, optional |
| `dhcp` (default `false`) | `dhcp` (default `true`) | `dhcp_enabled`, API default when omitted |
| `no_gateway` | `enable_gateway_ip` | no equivalent; `gateway_ip` is computed when omitted |
| (none) | `gateway_ip`, `dns_nameservers` | unchanged |
| (none) | `name` | `name`, **required** |
| (none) | `host_route` | read-only under `current_state.host_routes` |
| (none) | `use_default_public_dns_resolver` | removed |
| `ip_pools[*]` | (none) | `allocation_pools` |

```terraform
resource "ovh_cloud_project_network_private_subnet" "subnet" {
  service_name = "<public cloud project ID>"
  network_id   = ovh_cloud_project_network_private.net.id
  region       = "GRA11"
  network      = "10.0.0.0/24"
  start        = "10.0.0.10"
  end          = "10.0.0.200"
  dhcp         = true
  no_gateway   = false
}

resource "ovh_cloud_network_private_vrack_subnet" "subnet" {
  service_name = "<public cloud project ID>"
  network_id   = ovh_cloud_network_private_vrack.net_gra.id
  region       = "GRA11"
  name         = "backend-gra"
  cidr         = "10.0.0.0/24"
  dhcp_enabled = true

  allocation_pools = [
    { start = "10.0.0.10", end = "10.0.0.200" },
  ]
}
```

~> Subnets created with v1 have no name on the OpenStack side. `name` is required on the new resource, so the first apply after the import renames the subnet; this is an in-place update with no network impact. Check that `cidr`, `dhcp_enabled` and `allocation_pools` match the live values, since those are the arguments that would otherwise trigger a replacement or an update.

### Region networks

`ovh_cloud_project_region_network` created a network and its subnet in one resource. It becomes an `ovh_cloud_network_private_vrack` and an `ovh_cloud_network_private_vrack_subnet`, imported separately: `<service_name>/<network_id>` and `<service_name>/<network_id>/<subnet_id>`, where `network_id` is the legacy resource `id`.

| Legacy | New |
|---|---|
| `region_name` | `region` on both resources |
| `name`, `vlan_id` | on the network |
| `subnet.cidr`, `subnet.gateway_ip`, `subnet.allocation_pools` | on the subnet, same names |
| `subnet.enable_dhcp` | `dhcp_enabled` |
| `subnet.dns_name_servers` | `dns_nameservers` |
| `subnet.name` | `name`, **required** |
| `subnet.ip_version`, `subnet.enable_gateway_ip`, `subnet.use_default_public_dnsresolver`, `subnet.host_routes` | removed |

```terraform
resource "ovh_cloud_project_region_network" "lz" {
  service_name = "<public cloud project ID>"
  region_name  = "EU-SOUTH-LZ-MAD-A"
  name         = "madrid"
  subnet = {
    cidr        = "10.1.0.0/24"
    enable_dhcp = true
    ip_version  = 4
    name        = "madrid-subnet"
  }
}

resource "ovh_cloud_network_private_vrack" "lz" {
  service_name = "<public cloud project ID>"
  region       = "EU-SOUTH-LZ-MAD-A"
  name         = "madrid"
}

resource "ovh_cloud_network_private_vrack_subnet" "lz" {
  service_name = "<public cloud project ID>"
  network_id   = ovh_cloud_network_private_vrack.lz.id
  region       = "EU-SOUTH-LZ-MAD-A"
  name         = "madrid-subnet"
  cidr         = "10.1.0.0/24"
  dhcp_enabled = true
}
```

## Gateways

`ovh_cloud_project_gateway` and `ovh_cloud_project_gateway_interface` → `ovh_cloud_gateway`. Import ID: `<service_name>/<gateway_id>`.

The gateway's main subnet and every `ovh_cloud_project_gateway_interface` become entries of one `subnet_ids` list, which updates in place. Forget each interface resource with a `removed` block; the import of the gateway picks the interfaces up.

| Legacy | New |
|---|---|
| `name` | unchanged; `description` is new |
| `model = "s"` | `external_gateway = { enabled = true, model = "S" }`; models are `S`, `M`, `L`, `XL`, `2XL`, `3XL` |
| `network_id` | removed |
| `subnet_id` | `subnet_ids[0]` |
| `subnet_id` of each `ovh_cloud_project_gateway_interface` | another entry of `subnet_ids` |
| `external_information[*].ips[*].ip` | `current_state.external_ip` |
| `interfaces[*].subnet_id` | `current_state.subnets[*].id` |
| `interfaces[*].{id, ip, network_id}`, interface `ip`, `network_id` | removed |
| `status` | `current_state.status` |

~> The legacy gateway always had external access. Set `external_gateway.enabled = true` on the new resource, otherwise the plan after import shows a change.

```terraform
resource "ovh_cloud_project_gateway" "gw" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  name         = "backend-gw"
  model        = "s"
  network_id   = ovh_cloud_project_network_private.net.regions_attributes[0].openstackid
  subnet_id    = ovh_cloud_project_network_private_subnet.subnet.id
}

resource "ovh_cloud_project_gateway_interface" "gw_other" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  id           = ovh_cloud_project_gateway.gw.id
  subnet_id    = ovh_cloud_project_network_private_subnet.other.id
}

resource "ovh_cloud_gateway" "gw" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  name         = "backend-gw"

  external_gateway = {
    enabled = true
    model   = "S"
  }

  subnet_ids = [
    ovh_cloud_network_private_vrack_subnet.subnet.id,
    ovh_cloud_network_private_vrack_subnet.other.id,
  ]
}
```

```terraform
import {
  to = ovh_cloud_gateway.gw
  id = "<public cloud project ID>/<gateway ID>"
}

removed {
  from = ovh_cloud_project_gateway.gw
  lifecycle {
    destroy = false
  }
}

removed {
  from = ovh_cloud_project_gateway_interface.gw_other
  lifecycle {
    destroy = false
  }
}
```

## Load balancers

`ovh_cloud_project_loadbalancer` → `ovh_cloud_loadbalancer`, plus one `ovh_cloud_loadbalancer_listener` per listener, one `ovh_cloud_loadbalancer_pool` per pool and one `ovh_cloud_loadbalancer_pool_member` per member. Import IDs:

| Resource | Import ID |
|---|---|
| `ovh_cloud_loadbalancer` | `<service_name>/<loadbalancer_id>` |
| `ovh_cloud_loadbalancer_listener` | `<service_name>/<loadbalancer_id>/<listener_id>` |
| `ovh_cloud_loadbalancer_pool` | `<service_name>/<loadbalancer_id>/<pool_id>` |
| `ovh_cloud_loadbalancer_pool_member` | `<service_name>/<loadbalancer_id>/<pool_id>/<member_id>` |

The legacy state does not expose the listener, pool and member IDs. Read them from the `ovh_cloud_loadbalancer_listeners`, `ovh_cloud_loadbalancer_pools` and `ovh_cloud_loadbalancer_pool_members` data sources, or from the API: `GET /v2/publicCloud/project/{serviceName}/loadbalancer/{loadbalancerId}/listener`.

Load balancer:

| Legacy | New |
|---|---|
| `region_name` | `region` |
| `flavor_id` (UUID) | `flavor_name`: `SMALL`, `MEDIUM`, `LARGE` or `XL`; a change recreates the load balancer |
| `name` | unchanged, **required** now |
| `description` | unchanged |
| `network { private { network { id, subnet_id } } }` | `network = { id, subnet_id }` |
| `network { private { floating_ip { id } } }` | `network = { id, subnet_id, ip = "<floating IP address>" }` |
| `network { private { floating_ip_create } }` | an `ovh_cloud_floating_ip` resource, its `id` is the address to put in `network.ip` |
| `network { private { gateway, gateway_create } }` | an `ovh_cloud_gateway` resource |
| `listeners [...]` | separate resources, see below |
| `vip_address` | the `ip` of `current_state.network.addresses` whose `type` is `FIXED` |
| `floating_ip.ip` | the `ip` of `current_state.network.addresses` whose `type` is `FLOATING` |
| `vip_network_id`, `vip_subnet_id` | `current_state.network.{id, subnet_id}` |
| `operating_status`, `provisioning_status` | `current_state.{operating_status, provisioning_status}` |

Listener (`listeners[*]` → `ovh_cloud_loadbalancer_listener`):

| Legacy | New |
|---|---|
| `port` | `protocol_port` |
| `protocol = "tcp"` | `protocol = "TCP"`; `HTTP`, `HTTPS`, `SCTP`, `TCP`, `TERMINATED_HTTPS`, `UDP`; `prometheus` has no equivalent |
| `name` | unchanged, **required** now |
| `secret_id` | `default_tls_container_ref` (a Key Manager container reference) |
| `pool { ... }` | a `ovh_cloud_loadbalancer_pool` resource referenced through `default_pool_id` |
| `allowed_cidrs`, `description`, `timeout_client_data`, `timeout_member_data`, `tls_versions` | unchanged; timeouts are in milliseconds |

Pool (`listeners[*].pool` → `ovh_cloud_loadbalancer_pool`, attached to the load balancer with `loadbalancer_id`, not to the listener):

| Legacy | New |
|---|---|
| `protocol = "http"` | `protocol = "HTTP"`, **required** |
| `algorithm = "roundRobin"` | `algorithm = "ROUND_ROBIN"`, `LEAST_CONNECTIONS`, `SOURCE_IP`, **required** |
| `session_persistence { type = "httpCookie", cookie_name }` | `persistence = { type = "HTTP_COOKIE", cookie_name }`; `APP_COOKIE`, `SOURCE_IP`; omit the block for `disabled` |
| `health_monitor { monitor_type = "http" }` | `health_monitor = { type = "HTTP" }`; `delay`, `timeout` and `max_retries` are required |
| `health_monitor { http_configuration { url_path, http_method, expected_codes, domain_name, http_version } }` | flattened into `health_monitor = { url_path, http_method, expected_codes, domain_name, http_version }` |

Member (`listeners[*].pool.members[*]` → `ovh_cloud_loadbalancer_pool_member`):

| Legacy | New |
|---|---|
| `address`, `protocol_port`, `name` | unchanged; `address` and `protocol_port` recreate the member |
| `weight` 1-256 | `weight` 0-256 |
| (none) | `pool_id` and `loadbalancer_id` **required**; `subnet_id`, `backup`, `monitor` are new |

Before:

```terraform
resource "ovh_cloud_project_loadbalancer" "lb" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA11"
  flavor_id    = "<small flavor ID>"
  name         = "web-lb"

  network = {
    private = {
      network = {
        id        = ovh_cloud_project_network_private.net.regions_attributes[0].openstackid
        subnet_id = ovh_cloud_project_network_private_subnet.subnet.id
      }
    }
  }

  listeners = [
    {
      name     = "http"
      port     = 80
      protocol = "http"
      pool = {
        name      = "web"
        protocol  = "http"
        algorithm = "roundRobin"
        members = [
          { name = "web-1", address = "10.0.0.11", protocol_port = 8080 },
          { name = "web-2", address = "10.0.0.12", protocol_port = 8080 },
        ]
      }
    },
  ]
}
```

After:

```terraform
resource "ovh_cloud_loadbalancer" "lb" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  flavor_name  = "SMALL"
  name         = "web-lb"

  network = {
    id        = ovh_cloud_network_private_vrack.net_gra.id
    subnet_id = ovh_cloud_network_private_vrack_subnet.subnet.id
  }
}

resource "ovh_cloud_loadbalancer_pool" "web" {
  service_name    = "<public cloud project ID>"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  name            = "web"
  protocol        = "HTTP"
  algorithm       = "ROUND_ROBIN"
}

resource "ovh_cloud_loadbalancer_listener" "http" {
  service_name    = "<public cloud project ID>"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  name            = "http"
  protocol        = "HTTP"
  protocol_port   = 80
  default_pool_id = ovh_cloud_loadbalancer_pool.web.id
}

resource "ovh_cloud_loadbalancer_pool_member" "web_1" {
  service_name    = "<public cloud project ID>"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  pool_id         = ovh_cloud_loadbalancer_pool.web.id
  name            = "web-1"
  address         = "10.0.0.11"
  protocol_port   = 8080
}

resource "ovh_cloud_loadbalancer_pool_member" "web_2" {
  service_name    = "<public cloud project ID>"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  pool_id         = ovh_cloud_loadbalancer_pool.web.id
  name            = "web-2"
  address         = "10.0.0.12"
  protocol_port   = 8080
}
```

```terraform
import {
  to = ovh_cloud_loadbalancer.lb
  id = "<public cloud project ID>/<loadbalancer ID>"
}

import {
  to = ovh_cloud_loadbalancer_pool.web
  id = "<public cloud project ID>/<loadbalancer ID>/<pool ID>"
}

import {
  to = ovh_cloud_loadbalancer_listener.http
  id = "<public cloud project ID>/<loadbalancer ID>/<listener ID>"
}

import {
  to = ovh_cloud_loadbalancer_pool_member.web_1
  id = "<public cloud project ID>/<loadbalancer ID>/<pool ID>/<member ID of web-1>"
}

import {
  to = ovh_cloud_loadbalancer_pool_member.web_2
  id = "<public cloud project ID>/<loadbalancer ID>/<pool ID>/<member ID of web-2>"
}

removed {
  from = ovh_cloud_project_loadbalancer.lb
  lifecycle {
    destroy = false
  }
}
```

~> **Changing the flavor recreates the new load balancer**, where the legacy resource resized it in place. Translate `flavor_id` to the `flavor_name` of the live load balancer (`current_state.flavor.id` after import, or the `ovh_cloud_loadbalancer` data source before) so the plan after import shows no replacement.

~> Log subscriptions keep their legacy resource, `ovh_cloud_project_region_loadbalancer_log_subscription`; its `loadbalancer_id` takes the same ID.

## SSH keys

`ovh_cloud_project_ssh_key` → `ovh_cloud_ssh_key`. Import ID: `<service_name>/<name>`.

| Legacy | New |
|---|---|
| `name`, `public_key` | unchanged; both recreate the key |
| `region` | removed; the key is project-wide |
| `finger_print`, `regions` | removed |

```terraform
resource "ovh_cloud_project_ssh_key" "deploy" {
  service_name = "<public cloud project ID>"
  name         = "deploy"
  public_key   = "ssh-ed25519 AAAA... deploy@example"
}

resource "ovh_cloud_ssh_key" "deploy" {
  service_name = "<public cloud project ID>"
  name         = "deploy"
  public_key   = "ssh-ed25519 AAAA... deploy@example"
}
```

Instances reference it through `ssh_key_name = ovh_cloud_ssh_key.deploy.name`. The key is only read when an instance is created or rebuilt: migrating it has no effect on running instances.

## Failover IP attachments

`ovh_cloud_project_failover_ip_attach` has no dedicated replacement. An additional IP (the current name of a failover IP) is attached by listing it in the `networks` attribute of the `ovh_cloud_instance` it is routed to, with only `ip` set. Forget the legacy resource with a `removed` block and add the entry to the instance; the plan must show no change on the instance, since the IP is already routed to it.

```terraform
resource "ovh_cloud_project_failover_ip_attach" "ip" {
  service_name = "<public cloud project ID>"
  ip           = "192.0.2.10"
  routed_to    = ovh_cloud_project_instance.app.id
}

resource "ovh_cloud_instance" "app" {
  # ...
  networks = [
    { auto_assign_public_ip = true },
    { ip = "192.0.2.10" },
  ]
}
```

Moving the IP to another instance means removing the entry from one instance and adding it to the other in the same apply. The IP appears under `current_state.networks[*].addresses` with `type = "ADDITIONAL"`.

-> `ovh_cloud_floating_ip` is not the replacement of a failover IP: it manages a floating IP on Ext-Net, what the legacy `floating_ip_create` blocks of the instance and load balancer used to create. It is associated by putting its address in `networks[*].ip` (instance) or `network.ip` (load balancer).

## Data sources

| Legacy | New | Notes |
|---|---|---|
| `ovh_cloud_project_flavor`, `_flavors` | `ovh_cloud_instance_flavor`, `_flavors` | `ram` moves from GiB to MB; `name_filter` is gone, filter with a `for` expression |
| `ovh_cloud_project_image`, `_images` | `ovh_cloud_instance_image`, `_images` | `image_id` becomes `id`; `size` moves from GiB to bytes; `os_type` and `flavor_type` filters are gone |
| `ovh_cloud_project_instance`, `_instances` | `ovh_cloud_instance`, `_instances` | lookup by `service_name` and `id`, no `region`; `flavor_name` becomes `current_state.flavor.name`, `addresses` becomes `current_state.networks[*].addresses`, `status` becomes `current_state.power_state` |
| `ovh_cloud_project_floatingips` | `ovh_cloud_floating_ips` | no `region_name`; `id` is the address, the OpenStack UUID is `current_state.id` |
| `ovh_cloud_project_gateway` | `ovh_cloud_gateway`, `_gateways` | interfaces are `current_state.subnets[*].id`; `ovh_cloud_project_gateway_interface` has no equivalent |
| `ovh_cloud_project_loadbalancer`, `_loadbalancers` | `ovh_cloud_loadbalancer`, `_loadbalancers` | `vip_address` is the `FIXED` entry of `current_state.network.addresses`; `ovh_cloud_project_loadbalancer_flavors` has no equivalent, flavors are `SMALL`, `MEDIUM`, `LARGE`, `XL` |
| `ovh_cloud_project_network_private`, `_privates`, `_private_subnets` | `ovh_cloud_network_private_vrack`, `_vracks`, `_vrack_subnet`, `_vrack_subnets` | regional: `id` is the OpenStack UUID, `regions_openstack_ids` is gone; `ip_pools` becomes `allocation_pools` |
| `ovh_cloud_project_region`, `_regions` | `ovh_cloud_region`, `_regions` | `names` becomes `regions[*].name`; `status` is `ENABLED`, `DISABLED` or `MAINTENANCE`; `type` and `has_services_up` are gone |
| `ovh_cloud_project_ssh_key`, `_ssh_keys` | `ovh_cloud_ssh_key`, `_ssh_keys` | lookup by `name`; `finger_print` and `regions` are gone |
| `ovh_cloud_project_storage`, `_storages` | `ovh_cloud_storage_object_bucket`, `_buckets` | lookup by `id` (`<REGION>_<name>` on multi-region deployments); objects are no longer listed |
| `ovh_cloud_project_volume`, `_volumes` | `ovh_cloud_storage_block_volume`, `_volumes` | lookup by `id`, no `region_name`; `attached_instances[*].id` is new |

New data sources with no legacy counterpart: `ovh_cloud_public_ips` (every public IP of the project with its type), `ovh_cloud_ext_net_ip` and `_ips` (Ext-Net IPs), `ovh_cloud_additional_ip` and `_ips` (additional IPs), `ovh_cloud_quota`.

Before, a flavor lookup by name:

```terraform
data "ovh_cloud_project_flavors" "gra" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
  name_filter  = "b3-8"
}

resource "ovh_cloud_project_instance" "app" {
  # ...
  flavor {
    flavor_id = one(data.ovh_cloud_project_flavors.gra.flavors).id
  }
}
```

After:

```terraform
data "ovh_cloud_instance_flavors" "gra" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
}

resource "ovh_cloud_instance" "app" {
  # ...
  flavor_id = one([
    for flavor in data.ovh_cloud_instance_flavors.gra.flavors :
    flavor.id if flavor.name == "b3-8"
  ])
}
```

## Checklist

- `terraform plan` shows no changes once the `import` and `removed` blocks are applied and deleted.
- No `ovh_cloud_project_*` block remains for a resource that has a new equivalent, and no deprecation warning is printed.
- Every reference to a legacy resource (`.id`, `.regions_attributes`, `.addresses`...) points to the new resource and, where needed, to `current_state`.
- Legacy resources without an equivalent (bucket lifecycle configurations, replication jobs, load balancer log subscriptions) still reference the migrated objects by the same IDs.
