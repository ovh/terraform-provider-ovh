---
subcategory : "Managed Private Registry (MPR)"
---

# ovh_cloud_project_containerregistry_iam

Creates an IAM configuration in an OVHcloud Managed Private Registry.

## Example Usage

```terraform
resource "ovh_cloud_project_containerregistry_iam" "registry_iam" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"

  #optional field
  delete_users = false
}

output "iam_enabled" {
  value     = ovh_cloud_project_containerregistry_iam.registry_iam.iam_enabled
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `registry_id` - The ID of the Managed Private Registry. **Changing this value recreates the resource.**
* `delete_users` - Delete existing users from Harbor. IAM can't be enabled if there is at least one user already created. This parameter is only used at IAM configuration creation. **Changing this value recreates the resource.**

## Timeouts

```terraform
resource "ovh_cloud_project_containerregistry_iam" "registry_iam" {
  # ...

  timeouts {
    create = "1h"
    update = "45m"
    delete = "50s"
  }
}
```
* `create` - (Default 10m)
* `update` - (Default 10m)
* `delete` - (Default 10m)

## Import

OVHcloud Managed Private Registry IAM can be imported using the tenant `service_name` and registry id `registry_id` separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_containerregistry_iam.my-iam service_name/registry_id
```
