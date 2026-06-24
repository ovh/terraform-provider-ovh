---
subcategory : "Cloud Project"
---

# ovh_cloud_quota (Resource)

Manages the public cloud project quota: the applied quota profile per region
and the `prevent_automatic_quota_upgrade` toggle. There is exactly one quota
envelope per project (it is lazily initialized on the first read), so this
resource is a singleton — the project's service name is used as the import id.

Regions omitted from `regions` are left unchanged upstream.

## Example Usage

```hcl
resource "ovh_cloud_quota" "this" {
  service_name                    = <Public cloud project id>
  prevent_automatic_quota_upgrade = false
  regions = [
    {
      region  = "GRA11"
      profile = "50vms"
    },
    {
      region  = "DE1"
      profile = "default"
    },
  ]
}
```

## Argument Reference

- `service_name` (Required, ForceNew) — Service name of the cloud project.
- `prevent_automatic_quota_upgrade` (Required, bool) — When true, automatic
  quota upgrades are disabled for this project.
- `regions` (Required, list) — Target quota profile per region:
  - `region` (Required, string) — Region where the profile applies
    (e.g. `GRA11`).
  - `profile` (Required, string) — Quota profile to apply. Available values
    are exposed live in `current_state.available_profiles`.

## Attributes Reference

- `id` — Resource identifier (the project id).
- `resource_status` — Quota readiness in the system (`CREATING`, `UPDATING`,
  `DELETING`, `OUT_OF_SYNC`, `READY`, `ERROR`, `SUSPENDED`, `UNKNOWN`).
- `checksum` — Computed hash of the current target specification, used for
  optimistic concurrency on updates.
- `created_at`, `updated_at` — Envelope timestamps (RFC3339).
- `current_state` — Live state of the quota:
  - `prevent_automatic_quota_upgrade`
  - `available_profiles` — All quota profiles offered to the project, with
    their per-service caps (`compute`, `volume`, `network`, `loadbalancer`,
    `key_manager`, `share`, `keypair`).
  - `regions` — Per-region live usage (`limit` / `used` / `unit`) and the
    currently applied profile.

## Import

The quota envelope is a singleton per project — import by service name only:

```bash
terraform import ovh_cloud_quota.this <service_name>
```
