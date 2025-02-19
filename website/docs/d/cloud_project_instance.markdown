---
subcategory : "Instances"
---

# ovh_cloud_project_instance
**This datasource uses a Beta API**
Use this data source to get the instance of a public cloud project.

## Example Usage

To get information of an instance:

```hcl
data "ovh_cloud_project_instance" "instance" {
  service_name = "YYYY"
  region = "XXXX"
  instance_id = "ZZZZZ"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used
* `region` - (Required) Instance region
* `instance_id` - (Required) Instance id

## Attributes Reference

The following attributes are exported:

* `addresses` - Instance IP addresses
  * `ip` - IP address
  * `version` - IP version
* `attached_volumes` - Volumes attached to the instance
  * `id` - Volume id
* `flavor_id` - Flavor id
* `flavor_name` - Flavor name
* `id` - Instance id
* `name` - Instance name
* `image_id` - Image id
* `task_state` - Instance task state
* `ssh_key` - SSH Keypair
