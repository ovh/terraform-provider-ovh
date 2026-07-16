---
subcategory: "Instances"
---

# ovh_cloud_instance_group

Creates an instance group (placement group) in a public cloud project.

~> **WARNING** An instance group is **immutable**: it has no update route. Changing any attribute (`name`, `region`, `policy`) forces the group to be destroyed and re-created.

~> **NOTE** Membership is set only at instance-creation time through the instance's `group_id`; this resource never manages members.

## Example Usage

```terraform
resource "ovh_cloud_instance_group" "group" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance-group"
  policy       = "ANTI_AFFINITY"
}
```

### AFFINITY vs ANTI_AFFINITY

`ANTI_AFFINITY` spreads the members across distinct hypervisors: use it for the
replicas of a highly-available service, so a single host failure cannot take
them all down. `AFFINITY` does the opposite and packs the members onto the same
hypervisor, for tightly coupled instances that benefit from staying on one host.

```terraform
resource "ovh_cloud_instance_group" "spread" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-anti-affinity-group"
  policy       = "ANTI_AFFINITY"
}

resource "ovh_cloud_instance_group" "packed" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-affinity-group"
  policy       = "AFFINITY"
}
```

### Adding member instances

Members join through the instance's `group_id`, never through this resource.
`group_id` is itself immutable, so moving an instance to another group — or out
of its group — recreates the instance. Read the resulting membership back from
`current_state.members`.

```terraform
resource "ovh_cloud_instance_group" "ha" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-ha-group"
  policy       = "ANTI_AFFINITY"
}

resource "ovh_cloud_instance" "replica" {
  count = 3

  service_name = "<Public cloud project id>"
  region       = ovh_cloud_instance_group.ha.region
  name         = "my-replica-${count.index}"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  group_id     = ovh_cloud_instance_group.ha.id

  networks = [
    { auto_assign_public_ip = true },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `region` - (Required) Region where the instance group will be created. **Changing this value recreates the resource.**
* `name` - (Required) Instance group name. **Changing this value recreates the resource.**
* `policy` - (Required) Placement policy applied to the group's member instances (`AFFINITY`, `ANTI_AFFINITY`). **Changing this value recreates the resource.**

## Attributes Reference

The following attributes are exported:

* `id` - Instance group ID.
* `checksum` - Computed hash representing the current target specification value. Because an instance group has no update route, this value never changes after creation.
* `created_at` - Creation date of the instance group.
* `resource_status` - Instance group readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`). If the group drifts from its target specification — for example its underlying Nova server group is changed or removed outside Terraform — it is reported as `OUT_OF_SYNC`; since this resource cannot be updated in place, the next `terraform apply` **destroys and recreates it** to converge back to the Terraform configuration. See [Operational notes](#operational-notes).
* `current_state` - Current state of the instance group:
  * `name` - Instance group name.
  * `policy` - Placement policy.
  * `location` - Location details:
    * `region` - Region.
    * `availability_zone` - Availability zone.
  * `members` - Instances currently belonging to the group:
    * `id` - Instance ID.

~> **NOTE** There is no `updated_at` attribute. An instance group has no update route, so its target specification never changes and the API envelope deliberately omits the field.

## Operational notes

* Create and delete poll the API until the group settles, for up to **20
  minutes** each. The resource offers no `timeouts {}` block, so that budget
  cannot be shortened or extended from the configuration.
* `resource_status = "ERROR"` is terminal: polling stops at once and the
  provider surfaces the summary of the failed task(s) instead of a generic
  unexpected-state error.
* `resource_status = "OUT_OF_SYNC"` is remediated by **replacement**, not by an
  update. A group that drifts is planned for destroy-and-recreate on the next
  `terraform apply`; it never sits drifted. Replacement is the only
  reconciliation path available, because the group has no update route.
* Because the replacement assigns a **new group ID**, every
  `ovh_cloud_instance` whose immutable `group_id` points at this group is
  replaced along with it. Plan a drift remediation as an outage of the whole
  group, not of the placement group alone.

## Import

An instance group in a public cloud project can be imported using the
`service_name` and `instance_group_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_instance_group.group
  id = "<service_name>/<instance_group_id>"
}
```

```bash
$ terraform import ovh_cloud_instance_group.group service_name/instance_group_id
```
