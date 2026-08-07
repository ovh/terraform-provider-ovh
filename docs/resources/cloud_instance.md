---
subcategory: "Instances"
---

# ovh_cloud_instance

Creates an instance in a public cloud project.

~> **WARNING** Changing `image_id` rebuilds the instance and **wipes the root disk**. Back up any data on the root disk before changing the image.

## Example Usage

```terraform
resource "ovh_cloud_instance" "instance" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  ssh_key_name = "my-ssh-key"

  networks = [
    # Public interface with a platform-assigned Ext-Net IP.
    { auto_assign_public_ip = true },
  ]
}
```

### Resolving the flavor, the image and the SSH key

`flavor_id` and `image_id` are opaque per-region identifiers. Rather than
hardcoding them, resolve them from the
[`ovh_cloud_instance_flavors`](../data-sources/cloud_instance_flavors.md) and
[`ovh_cloud_instance_images`](../data-sources/cloud_instance_images.md) data
sources, and point `ssh_key_name` at an
[`ovh_cloud_ssh_key`](cloud_ssh_key.md) resource.

```terraform
data "ovh_cloud_instance_flavors" "catalog" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
}

data "ovh_cloud_instance_images" "catalog" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
}

resource "ovh_cloud_ssh_key" "deploy" {
  service_name = "<Public cloud project id>"
  name         = "my-deploy-key"
  public_key   = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample user@host"
}

resource "ovh_cloud_instance" "from_catalog" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"

  name = "my-catalog-instance"

  flavor_id = one([
    for flavor in data.ovh_cloud_instance_flavors.catalog.flavors :
    flavor.id if flavor.name == "b3-8"
  ])

  image_id = one([
    for image in data.ovh_cloud_instance_images.catalog.images :
    image.id if image.name == "Debian 12"
  ])

  ssh_key_name = ovh_cloud_ssh_key.deploy.name

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

### Network interfaces

A `networks[]` entry supports four shapes. Entries keep the order they are
written in: the API returns them sorted by network id and the provider
re-orders them back to the configuration.

```terraform
resource "ovh_cloud_instance" "networks" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-multi-homed-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"

  networks = [
    # 1. Public interface with a public IP assigned by the platform.
    #    At most one entry may set auto_assign_public_ip.
    { auto_assign_public_ip = true },

    # 2. A public IP the project already owns: an additional IP, or an Ext-Net
    #    IP of the project in the instance's region. Several such entries are
    #    allowed, and they may coexist with auto_assign_public_ip.
    { ip = "203.0.113.10" },

    # 3. Private interface with an address picked by IPAM.
    {
      network_id = "<private network id>"
      subnet_id  = "<subnet id>"
    },

    # 4. Private interface with an explicit address: the port's fixed IP is
    #    pinned to `ip` when it falls inside the subnet CIDR, otherwise the
    #    existing floating IP with that address is associated.
    {
      network_id = "<private network id>"
      subnet_id  = "<subnet id>"
      ip         = "10.0.0.42"
    },
  ]
}
```

### Private-network-only instance

```terraform
resource "ovh_cloud_network_private_vrack" "private" {
  service_name = "<Public cloud project id>"
  name         = "my-private-network"
  region       = "GRA11"
}

resource "ovh_cloud_network_private_vrack_subnet" "private" {
  service_name = ovh_cloud_network_private_vrack.private.service_name
  network_id   = ovh_cloud_network_private_vrack.private.id
  name         = "my-subnet"
  cidr         = "10.0.0.0/24"
  region       = "GRA11"
}

# No public interface at all: the instance is only reachable from the vRack.
# Add an ovh_cloud_gateway on the same subnet to give it egress.
resource "ovh_cloud_instance" "private_only" {
  service_name = ovh_cloud_network_private_vrack.private.service_name
  region       = "GRA11"
  name         = "my-private-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"

  networks = [
    {
      network_id = ovh_cloud_network_private_vrack.private.id
      subnet_id  = ovh_cloud_network_private_vrack_subnet.private.id
    },
  ]
}
```

### Boot from volume

Omit `image_id` and pass a bootable volume in `volume_ids`. The instance's
`current_state.image` stays null for a boot-from-volume instance.

```terraform
resource "ovh_cloud_storage_block_volume" "boot" {
  service_name = "<Public cloud project id>"
  name         = "my-boot-volume"
  size         = 50
  region       = "GRA11"
  volume_type  = "HIGH_SPEED_GEN2"

  create_from = {
    image_id = "<image id>"
  }
}

resource "ovh_cloud_instance" "boot_from_volume" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-boot-from-volume-instance"
  flavor_id    = "<flavor id>"
  volume_ids   = [ovh_cloud_storage_block_volume.boot.id]

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

### Attaching additional block volumes

`volume_ids` is mutable: adding or removing an id attaches or detaches the
volume in place, without recreating the instance.

```terraform
resource "ovh_cloud_storage_block_volume" "data" {
  service_name = "<Public cloud project id>"
  name         = "my-data-volume"
  size         = 100
  region       = "GRA11"
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_instance" "with_volumes" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance-with-volumes"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  volume_ids   = [ovh_cloud_storage_block_volume.data.id]

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

### Security groups

`security_group_ids` is optional and computed:

* **omitted** — the platform applies the project's `default` security group and
  writes it into the target spec, so after the first apply `security_group_ids`
  holds exactly that one group id. Removing the attribute from a configuration
  that used to set it therefore keeps the groups already attached; it does not
  detach them.
* **explicit empty list** (`security_group_ids = []`) — no security group is
  applied at all. The instance accepts **no inbound traffic**; use it only when
  the filtering is handled elsewhere.
* **explicit list of ids** — exactly those groups are applied to every interface.

```terraform
resource "ovh_cloud_security_group" "web" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-web-security-group"
  description  = "Allow SSH and HTTPS"

  rule = [
    {
      direction        = "INGRESS"
      ethernet_type    = "IPV4"
      protocol         = "TCP"
      port_range_min   = 22
      port_range_max   = 22
      remote_ip_prefix = "0.0.0.0/0"
      description      = "SSH"
    },
  ]
}

resource "ovh_cloud_instance" "filtered" {
  service_name       = "<Public cloud project id>"
  region             = "GRA11"
  name               = "my-filtered-instance"
  flavor_id          = "<flavor id>"
  image_id           = "<image id>"
  security_group_ids = [ovh_cloud_security_group.web.id]

  networks = [
    { auto_assign_public_ip = true },
  ]
}

# No security group at all: the instance accepts no inbound traffic.
resource "ovh_cloud_instance" "no_filtering" {
  service_name       = "<Public cloud project id>"
  region             = "GRA11"
  name               = "my-unfiltered-instance"
  flavor_id          = "<flavor id>"
  image_id           = "<image id>"
  security_group_ids = []

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

### Attaching file storage shares

```terraform
resource "ovh_cloud_storage_file_share" "share" {
  service_name     = "<Public cloud project id>"
  name             = "my-share"
  size             = 150
  region           = "GRA11"
  protocol         = "NFS"
  share_type       = "STANDARD_1AZ"
  share_network_id = "<share network id>"
}

# No ovh_cloud_storage_file_share_acl resource targets this share: the instance
# below owns its access rules.
resource "ovh_cloud_instance" "with_shares" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance-with-shares"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"

  shares = [
    {
      id = ovh_cloud_storage_file_share.share.id
    },
    {
      id           = "<read only share id>"
      access_level = "READ_ONLY"
    },
  ]

  networks = [
    {
      network_id = "<private network id>"
      subnet_id  = "<subnet id>"
    },
  ]
}
```

### Power state

```terraform
# The instance stays provisioned and keeps its disks, ports and IPs, but the
# VM is powered off.
resource "ovh_cloud_instance" "stopped" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-stopped-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  power_state  = "SHUTOFF"

  networks = [
    { auto_assign_public_ip = true },
  ]
}

# Powered off and released on the hypervisor. Switching power_state back to
# ACTIVE unshelves it in place.
resource "ovh_cloud_instance" "shelved" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-shelved-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  power_state  = "SHELVED"

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

### Joining an instance group

Group membership is only settable here, through the instance's `group_id`. See
[`ovh_cloud_instance_group`](cloud_instance_group.md).

```terraform
resource "ovh_cloud_instance_group" "spread" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-anti-affinity-group"
  policy       = "ANTI_AFFINITY"
}

resource "ovh_cloud_instance" "grouped" {
  count = 2

  service_name = "<Public cloud project id>"
  region       = ovh_cloud_instance_group.spread.region
  name         = "my-grouped-instance-${count.index}"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  group_id     = ovh_cloud_instance_group.spread.id

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `region` - (Required) Region where the instance is created. **Changing this value recreates the resource.**
* `name` - (Required) Instance name.
* `flavor_id` - (Required) Flavor ID. Changing it resizes the instance in place.
* `availability_zone` - (Optional) Availability zone of the instance (immutable; assigned by the platform if omitted). **Changing this value recreates the resource** — only when it is set in the configuration; a value assigned by the platform is never treated as a change.
* `ssh_key_name` - (Optional) Name of the SSH key injected at boot (immutable). Point it at the `name` of an [`ovh_cloud_ssh_key`](cloud_ssh_key.md). **Changing this value recreates the resource.**
* `group_id` - (Optional) ID of the placement group the instance belongs to (immutable). This is the only way to make an instance a member of an [`ovh_cloud_instance_group`](cloud_instance_group.md). **Changing this value recreates the resource.**
* `image_id` - (Optional) Image ID to boot from. Omit for a boot-from-volume instance. **WARNING**: changing it rebuilds the instance and **wipes the root disk**.
* `power_state` - (Optional) Desired power state: `ACTIVE`, `SHUTOFF` or `SHELVED`. When omitted, the API applies `ACTIVE` server-side and echoes it back; the provider declares no default of its own.
* `networks` - (Optional) Network interfaces attached to the instance. Entries keep the order they are written in; the API returns them sorted by network id and the provider re-orders them back to the configuration. Four shapes:
  * `auto_assign_public_ip` alone — public interface with a platform-assigned public IP (at most one such entry).
  * `ip` alone — a public IP the project already owns: an additional IP, or an Ext-Net IP of the project in the instance's region. Several are allowed and may coexist with `auto_assign_public_ip`.
  * `network_id` + `subnet_id` — private interface with an address picked by IPAM.
  * `network_id` + `subnet_id` + `ip` — pins the port's fixed address when `ip` is inside the subnet CIDR, otherwise associates the existing floating IP `ip`.

  Each entry supports:
  * `network_id` - (Optional) Private network ID. Omit for a public interface.
  * `subnet_id` - (Optional) Subnet ID within the private network. Required with `network_id`.
  * `ip` - (Optional) IP address of this interface. Without `network_id`: a public IP the project already owns (additional IP, or an Ext-Net IP of the project in the instance's region). With `network_id` + `subnet_id`: pins the port's fixed address when inside the subnet CIDR, otherwise associates the existing floating IP with that address.
  * `auto_assign_public_ip` - (Optional) Attach a public interface with a public IP assigned by the platform. Only valid on an entry with no `network_id` and no `ip`, and on at most one entry.
* `volume_ids` - (Optional) IDs of block-storage volumes attached to the instance.
* `security_group_ids` - (Optional) IDs of security groups applied to all interfaces. Omit it to let the platform apply the project's `default` security group; set an explicit empty list (`[]`) to apply no security group at all (the instance then accepts no inbound traffic).
* `shares` - (Optional) Filesystem shares attached to the instance. Each entry supports:
  * `id` - (Required) Identifier of the file storage share to attach. Adding the reference attaches the share to the instance, removing it detaches it.
  * `access_level` - (Optional) Access level granted to the instance for this share: `READ_ONLY` or `READ_WRITE`. Omit it to let the API apply its default (`READ_WRITE`). When omitted, the attribute stays null in state even though the API echoes `READ_WRITE` back, so an omitted level never shows up as drift.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the instance.
* `checksum` - Computed hash representing the current target specification value. It implements optimistic concurrency control: the value is echoed back on update and the request is rejected when it no longer matches server-side.
* `created_at` - Creation date of the instance, as an RFC 3339 timestamp.
* `updated_at` - Last modification date of the instance, as an RFC 3339 timestamp.
* `resource_status` - Instance readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`). Distinct from `current_state.power_state`, which carries the lower-level OpenStack administrative power state.
* `current_state` - Observed state of the instance as reported by the compute backend, as opposed to the requested specification exposed at root level. Null while the instance is still being created and no backend state is available yet:
  * `name` - Observed display name of the instance.
  * `power_state` - Observed administrative power state of the instance as reported by OpenStack. It may transiently differ from the requested `power_state` while a power transition is in progress.
  * `locked` - Whether the instance is locked against modifications. While locked, mutating actions are refused until it is unlocked.
  * `ssh_key_name` - Name of the SSH key pair injected into the instance at boot, null when none was provided.
  * `host_id` - Opaque identifier of the physical host the instance is running on, as exposed by OpenStack. Null when not available.
  * `project_id` - Identifier of the Public Cloud project the instance belongs to.
  * `user_id` - Identifier of the OpenStack user that owns the instance.
  * `flavor` - Observed flavor of the instance, with its full sizing details:
    * `id` - Unique identifier of the flavor.
    * `name` - Human-readable flavor name (the commercial flavor label).
    * `vcpus` - Number of virtual CPUs provided by the flavor.
    * `ram` - Amount of RAM provided by the flavor, in MB.
    * `disk` - Size of the flavor's local root disk, in GB.
    * `swap` - Size of the flavor's swap space, in MB.
    * `ephemeral` - Size of the flavor's ephemeral disk, in GB.
  * `image` - Observed image the instance was booted from, null for a boot-from-volume instance which has no image:
    * `id` - Unique identifier of the image.
    * `name` - Human-readable image name, null when the backend does not report it.
    * `size` - Size of the image, in bytes. Null when the backend does not report it.
    * `status` - Lifecycle status of the image as reported by Glance.
    * `deprecated` - Whether the image is flagged as deprecated. A deprecated image still boots existing instances but is no longer recommended for new ones.
  * `location` - Observed region and availability zone where the instance is provisioned:
    * `region` - Code of the region where the instance is provisioned (for example GRA11, BHS5).
    * `availability_zone` - Availability zone within the region where the instance is placed, null in regions that have none.
  * `networks` - Observed network interfaces of the instance: one entry per private network plus at most one entry without a network id for the public (Ext-Net) interface. Entries are ordered by network id, so they do not follow the order of the requested `networks`:
    * `id` - Identifier of the network this interface is attached to, null for the public (Ext-Net) interface.
    * `subnet_id` - Identifier of the subnet this interface draws its fixed address from, null for an entry without a network id.
    * `gateway_id` - Identifier of the gateway providing egress for this interface, null when none applies.
    * `addresses` - Addresses observed on this interface: its fixed addresses plus, where applicable, its floating IP and any additional IPs routed to it. Each address carries a type of `FIXED`, `FLOATING` or `ADDITIONAL`:
      * `ip` - IP address assigned to the interface (IPv4 or IPv6).
      * `mac` - MAC address of the interface this IP is bound to. Null when the backend reports no interface for the address, which happens for an additional IP routed to an instance whose public interface has no visible Ext-Net address yet.
      * `type` - How this address reaches the instance: `FIXED` for an address assigned to the interface itself, `FLOATING` for a floating IP NAT'd onto it, `ADDITIONAL` for an additional IP routed to the public interface.
      * `version` - IP version of the address (4 for IPv4, 6 for IPv6).
  * `volumes` - Observed block volumes attached to the instance:
    * `id` - Unique identifier of the attached volume.
    * `name` - Display name of the attached volume.
    * `size` - Size of the attached volume, in GB.
  * `shares` - Observed instance-side share attachments, derived from the Manila access rules that target one of the instance's IPs. Only populated on a single-instance read, never in the list data source:
    * `id` - Identifier of the attached file storage share.
    * `access_level` - Observed access level of the access rule for this instance (`READ_ONLY` or `READ_WRITE`).
    * `access_to` - The instance IP address the Manila access rule targets: its fixed IPv4 address on the share's network.
    * `state` - Observed state of the underlying access rule (`ACTIVE`, `APPLYING`, `DENYING`, `ERROR`). Null while no state has been reported yet.
  * `security_groups` - Security groups currently attached to the instance's ports:
    * `id` - Security group identifier.
  * `group` - Instance (placement) group the instance belongs to, null when it is not part of any group:
    * `id` - Identifier of the instance (placement) group.

`auto_assign_public_ip` is never reported in `current_state` (the platform
cannot distinguish an assigned address from a pinned one). The public interface
is the single observed entry with no `id`.

## Operational notes

* Create, update and delete poll the API until the instance settles, for up to
  **60 minutes** each. The resource offers no `timeouts {}` block, so that
  budget cannot be shortened or extended from the configuration.
* An update re-reads the instance immediately before issuing the `PUT`, so the
  `checksum` it sends is the freshest one. This keeps a concurrent server-side
  change from failing the apply with a `ChecksumMismatch` (HTTP 409).
* `resource_status = "ERROR"` is terminal: polling stops at once and the
  provider surfaces the summary of the failed task(s) instead of a generic
  unexpected-state error.

## Import

An instance in a public cloud project can be imported using the `service_name`
and `instance_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_instance.instance
  id = "<service_name>/<instance_id>"
}
```

```bash
$ terraform import ovh_cloud_instance.instance service_name/instance_id
```
