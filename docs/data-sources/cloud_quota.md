---
subcategory : "Cloud Project"
---

# ovh_cloud_quota (Data Source)

Fetch read-only quota information for a public cloud project: target quota profile per region, manual-quota flag, list of available profiles with their caps, and per-region applied profile + usage reported by OpenStack.

## Example Usage

```terraform
data "ovh_cloud_quota" "quota" {
  service_name = "XXXXXX"
}
```

To narrow the per-region quota state to a single region, pass `region`:

```terraform
data "ovh_cloud_quota" "gra11_quota" {
  service_name = "XXXXXX"
  region       = "GRA11"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `region` - (Optional) If set, restricts the per-region quota state to this single region. Otherwise all configured regions are returned.

## Attributes Reference

The following attributes are exported:

* `id` - Resource identifier (the project id).
* `resource_status` - Readiness of the resource.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date (RFC3339).
* `updated_at` - Last update date (RFC3339).
* `target_spec` - Desired quota specification:
  * `manual_quota` - When true, automatic quota upgrades are disabled for this project.
  * `regions` - Target quota profile per region:
    * `region` - Region where the profile applies.
    * `profile` - Quota profile to apply in this region.
* `current_state` - Current quota state:
  * `manual_quota` - When true, automatic quota upgrades are disabled for this project.
  * `available_profiles` - List of available quota profiles with their caps:
    * `name` - Profile name.
    * `compute` - Compute caps: `cores`, `instances`, `memory`.
    * `volume` - Block storage caps: `backup_size_total`, `backups`, `size_total`, `snapshots`, `volumes`.
    * `network` - Networking caps: `floating_ips`, `gateways`, `networks`, `security_group_rules`, `security_groups`, `subnets`.
    * `loadbalancer` - Load balancer caps: `health_monitors`, `l7_policies`, `l7_rules`, `listeners`, `loadbalancers`, `members`, `pools`.
    * `key_manager` - Key manager caps: `containers`, `secrets`.
    * `share` - Shared file system caps: `backup_size_total`, `backups`, `shares`, `size_total`, `snapshots`.
    * `keypair` - Keypair caps: `keypairs`.
  * `regions` - Per-region quota state:
    * `region` - Region name (e.g. `GRA11`).
    * `profile` - Currently applied quota profile name in this region.
    * `compute` - Compute usage: `cores`, `instances`, `memory`. Each entry exposes `limit`, `used`, `unit`.
    * `volume` - Block storage usage: `backup_size_total`, `backups`, `size_total`, `snapshots`, `volumes` (each with `limit`, `used`, `unit`) and `per_volume_size` (`limit`, `unit`).
    * `network` - Networking usage: `floating_ips`, `gateways`, `networks`, `security_group_rules`, `security_groups`, `subnets`.
    * `loadbalancer` - Load balancer usage: `health_monitors`, `l7_policies`, `l7_rules`, `listeners`, `loadbalancers`, `members`, `pools`.
    * `key_manager` - Key manager usage: `containers`, `secrets`.
    * `share` - Shared file system usage: `backup_size_total`, `backups`, `share_networks`, `shares`, `size_total`, `snapshot_size_total`, `snapshots` (each with `limit`, `used`, `unit`) and `per_share_size` (`limit`, `unit`).
    * `keypair` - Keypair usage: `keypairs`.

Each usage entry carries:

* `limit` - Maximum authorized value for this quota.
* `used` - Current usage reported by OpenStack. `null` when the underlying service does not expose usage.
* `unit` - Unit of the limit and used values (e.g. `count`, `GB`, `MB`).

Limit-only entries (such as `per_volume_size` and `per_share_size`) carry:

* `limit` - Maximum authorized value for this limit.
* `unit` - Unit of the limit value.
