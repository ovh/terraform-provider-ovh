---
subcategory : "Instances"
---

# ovh_cloud_project_instances

**This datasource uses a Beta API**

Use this data source to get the list of instances in a region of a public cloud project.

## Example Usage

To list your instances:

```terraform
data "ovh_cloud_project_instances" "instance" {
  service_name = "YYYY"
  region = "XXXX"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `region` - (Required) Instance region.

## Attributes Reference

The following attributes are exported:
* `instances` - List of instances
  * `addresses` - Instance IP addresses
    * `ip` - IP address
    * `version` - IP version
  * `attached_volumes` - Volumes attached to the instance
    * `id` - Volume id
  * `availability_zone` - Availability zone of the instance
  * `flavor_id` - Flavor id
  * `flavor_name` - Flavor name
  * `id` - Instance id
  * `name` - Instance name
  * `image_id` - Image id
  * `task_state` - Instance task state
  * `ssh_key` - SSH Keypair
