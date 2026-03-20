---
subcategory : "Cloud Instances"
---

# ovh_cloud_instance_group (Data Source)

Get information about an instance group in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_instance_group" "group" {
  service_name      = "xxxxxxxxxx"
  instance_group_id = "00000000-0000-0000-0000-000000000001"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `instance_group_id` - (Required) Instance group ID.

## Attributes Reference

The following attributes are exported:

* `id` - Instance group ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the instance group.
* `updated_at` - Last update date of the instance group.
* `resource_status` - Instance group readiness in the system.
* `target_spec` - Target specification:
  * `name` - Instance group name.
  * `policy` - Placement policy (`AFFINITY` or `ANTI_AFFINITY`).
  * `region` - Region.
* `current_state` - Current state of the instance group:
  * `name` - Instance group name.
  * `policy` - Placement policy.
  * `region` - Region.
  * `members` - List of instances in this group:
    * `id` - Instance identifier.
