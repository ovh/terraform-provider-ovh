---
subcategory : "VPS"
---

# ovh_vps_migration

Schedules the migration of an OVHcloud VPS from the 2020 to the 2025 generation.

~> **WARNING** This resource is destructive-ish. Creating it enqueues a migration that, once started, **cannot be rolled back**. Running `terraform destroy` only cancels a queued migration (status `planned` or `toPlan`) — it does **not** revert a completed (`done`) or in-progress (`ongoing`) migration. Once the migration is `done`, removing the resource only stops Terraform from tracking it.

## Example Usage

```hcl
resource "ovh_vps_migration" "now" {
  service_name = "vpsXXXXX.ovh.net"
  target_plan  = "vps-2025-le-2-4-40"
}

resource "ovh_vps_migration" "scheduled" {
  service_name   = "vpsXXXXX.ovh.net"
  target_plan    = "vps-2025-le-2-4-40"
  scheduled_date = "2025-12-31T22:00:00Z"
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) VPS service name (e.g. `vpsXXXXX.ovh.net`).
* `target_plan` - (Required) Target plan code for the migration. Must be one of the plan codes exposed in `available_plans`. Changing this value cancels the queued migration and enqueues a new one with the new plan.
* `scheduled_date` - (Optional) ISO-8601 datetime at which the migration should run. When omitted, the migration is enqueued for immediate execution by OVHcloud.

## Attributes Reference

* `id` - Set to the `service_name`.
* `current_plan` - The source (2020) plan currently set on the VPS.
* `status` - Current migration status, one of `available`, `done`, `notAvailable`, `ongoing`, `planned`, `toPlan`.
* `position` - Position of the migration in the queue, when applicable.
* `available_plans` - List of plan codes that this VPS can be migrated to.

## Import

A migration tracking record can be imported using the VPS `service_name`:

```sh
terraform import ovh_vps_migration.mig vpsXXXXX.ovh.net
```
