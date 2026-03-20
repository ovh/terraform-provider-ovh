---
subcategory : "Cloud Instances"
---

# ovh_cloud_instance_group

Creates an instance group (server group) in a public cloud project. Instance groups define placement policies that control how instances are distributed across compute hosts.

~> **Note:** Instance groups are immutable after creation. Changing any attribute will destroy and recreate the resource.

## Example Usage

```terraform
resource "ovh_cloud_instance_group" "anti_affinity" {
  service_name = "xxxxxxxxxx"
  name         = "web-anti-affinity"
  policy       = "ANTI_AFFINITY"
  region       = "GRA1"
}
```

### With instances

```terraform
resource "ovh_cloud_instance_group" "web_group" {
  service_name = "xxxxxxxxxx"
  name         = "web-anti-affinity"
  policy       = "ANTI_AFFINITY"
  region       = "GRA1"
}

resource "ovh_cloud_instance" "web_1" {
  service_name = "xxxxxxxxxx"
  name         = "web-1"
  flavor_id    = "068a57cf-b1b4-428f-9b17-4f32a526390c"
  image_id     = "8d75e170-1ef9-4e25-8fc8-d231929e56e8"
  region       = "GRA1"
  group_id     = ovh_cloud_instance_group.web_group.id
}

resource "ovh_cloud_instance" "web_2" {
  service_name = "xxxxxxxxxx"
  name         = "web-2"
  flavor_id    = "068a57cf-b1b4-428f-9b17-4f32a526390c"
  image_id     = "8d75e170-1ef9-4e25-8fc8-d231929e56e8"
  region       = "GRA1"
  group_id     = ovh_cloud_instance_group.web_group.id
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `name` - (Required) Instance group name. **Changing this value recreates the resource.**
* `policy` - (Required) Placement policy. Valid values are `AFFINITY` (place instances on the same host) and `ANTI_AFFINITY` (place instances on different hosts). **Changing this value recreates the resource.**
* `region` - (Required) Region where the instance group will be created. **Changing this value recreates the resource.**

## Attributes Reference

The following attributes are exported:

* `id` - Instance group ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the instance group.
* `updated_at` - Last update date of the instance group.
* `resource_status` - Instance group readiness in the system (`CREATING`, `DELETING`, `ERROR`, `READY`).
* `current_state` - Current state of the instance group:
  * `name` - Instance group name.
  * `policy` - Placement policy.
  * `region` - Region.
  * `members` - List of instances in this group:
    * `id` - Instance identifier.

## Import

An instance group can be imported using the `service_name` and `id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_instance_group.group
  id = "<service_name>/<instance_group_id>"
}
```

```bash
$ terraform plan -generate-config-out=instance_group.tf
$ terraform apply
```
