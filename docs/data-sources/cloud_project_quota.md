---
subcategory : "Cloud Project"
---

# ovh_cloud_project_quota (Data Source)

Fetch read-only quota information for a public cloud project: currently applied profile, list of available profiles with their caps, and per-region usage reported by OpenStack.

## Example Usage

```terraform
data "ovh_cloud_project_quota" "quota" {
  service_name = "XXXXXX"
}
```

To narrow the per-region usage to a single region, pass `region`:

```terraform
data "ovh_cloud_project_quota" "gra7_quota" {
  service_name = "XXXXXX"
  region       = "GRA7"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `region` - (Optional) If set, restricts the per-region quota usage to this single region. Otherwise all configured regions are returned.

## Attributes Reference

The following attributes are exported:

* `id` - Resource identifier (the project id).
* `resource_status` - Readiness of the resource.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date (RFC3339).
* `updated_at` - Last update date (RFC3339).
* `target_spec` - Desired quota specification:
  * `profile` - Name of the quota profile applied to the project.
* `current_state` - Current quota state:
  * `profile` - Name of the currently applied quota profile.
  * `available_profiles` - List of available quota profiles with their caps:
    * `name` - Profile name.
    * `compute` - Compute caps: `instances`, `cores`, `ram`, `security_groups`, `security_group_rules`, `server_groups`, `server_group_members`.
    * `block_storage` - Block storage caps: `volumes`, `gigabytes`, `snapshots`, `backups`, `backup_gigabytes`.
    * `network` - Networking caps: `networks`, `subnets`, `floating_ips`, `gateways`, `security_groups`, `security_group_rules`.
    * `loadbalancer` - Load balancer caps: `loadbalancers`, `listeners`, `pools`, `members`, `healthmonitors`, `l7_policies`, `l7_rules`.
    * `key_manager` - Key manager caps: `secrets`, `containers`.
    * `share` - Shared file system caps: `shares`, `gigabytes`, `snapshots`, `backups`, `backup_gigabytes`.
    * `keypair` - Keypair caps: `keypairs`.
  * `regions` - Per-region quota usage:
    * `region` - Region name (e.g. `GRA7`).
    * `compute` - Compute usage: `instances`, `cores`, `memory`. Each entry exposes `limit`, `used` and `unit`.
    * `volume` - Block storage usage: `volumes`, `gigabytes`, `snapshots`, `backups`, `backup_gigabytes` (each with `limit`, `used`, `unit`) and `per_volume_size` (`limit`, `unit`).
    * `network` - Networking usage: `networks`, `subnets`, `floating_ips`, `gateways`, `security_groups`, `security_group_rules`.
    * `loadbalancer` - Load balancer usage: `loadbalancers`, `listeners`, `pools`, `members`, `healthmonitors`, `l7_policies`, `l7_rules`.
    * `key_manager` - Key manager usage: `secrets`, `containers`.
    * `share` - Shared file system usage: `shares`, `size_total`, `snapshots`, `snapshot_gigabytes`, `backups`, `backup_gigabytes`, `share_networks` and `per_share_size`.
    * `keypair` - Keypair usage: `keypairs`.

Each usage entry carries:

* `limit` - Maximum authorized value for this quota.
* `used` - Current usage reported by OpenStack. `null` when the underlying service does not expose usage.
* `unit` - Unit of the limit and used values (e.g. `count`, `GB`, `MB`).

Limit-only entries (such as `per_volume_size` and `per_share_size`) carry:

* `limit` - Maximum authorized value for this limit.
* `unit` - Unit of the limit value.
