---
subcategory: "Instances"
---

# ovh_cloud_instance_group (Data Source)

Use this data source to retrieve information about an instance group (placement group) in a public cloud project. An instance group is immutable once created: there is no update route and its membership is fixed at instance-creation time.

## Example Usage

```terraform
data "ovh_cloud_instance_group" "group" {
  service_name = "<Public cloud project id>"
  id           = "<instance group id>"
}

# Membership is fixed at instance-creation time through the instance's
# group_id; this is the only way to read it back.
output "group_member_ids" {
  value = [
    for member in data.ovh_cloud_instance_group.group.current_state.members :
    member.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `id` - (Required) Unique identifier of the instance group to look up.

## Attributes Reference

The following attributes are exported:

* `region` - Region where the instance group, and therefore its member instances, are placed. Immutable after creation: a group cannot be moved to another region, and its members must be created in this same location.
* `name` - Display name of the instance group. Immutable after creation, as instance groups cannot be updated.
* `policy` - Placement policy applied to the group's member instances: `AFFINITY` packs members onto the same hypervisor, `ANTI_AFFINITY` spreads members across distinct hypervisors. Maps to the underlying OpenStack/Nova server group policy.
* `checksum` - Computed hash of the current target specification, used for optimistic concurrency control. Because an instance group has no update route, this value never changes after creation.
* `created_at` - Timestamp at which the instance group was created, in RFC 3339 format.
* `resource_status` - Instance group readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`). `OUT_OF_SYNC` means the group has drifted from its target specification; a group managed by [`ovh_cloud_instance_group`](../resources/cloud_instance_group.md) is destroyed and recreated on the next `terraform apply` to converge.
* `current_state` - State of the instance group as observed on the backend:
  * `name` - Display name of the instance group as reported by the backend. Fixed for the lifetime of the group, since instance groups cannot be updated.
  * `policy` - Placement policy currently enforced for the group's members (`AFFINITY`, `ANTI_AFFINITY`). Mirrors the underlying OpenStack/Nova server group policy and is fixed at creation.
  * `location` - Region (and, where applicable, availability zone) where the instance group and its member instances are deployed, as observed on the backend:
    * `region` - Region code.
    * `availability_zone` - Availability zone within the region.
  * `members` - Instances currently belonging to this group. Membership is determined at instance-creation time via the instance's `group` field and cannot be changed afterwards. Empty when the group has no member instances:
    * `id` - Identifier of the member instance.

~> **NOTE** There is no `updated_at` attribute. An instance group has no update route, so its target specification never changes and the API envelope deliberately omits the field.
