---
subcategory : "Managed Private Registry (MPR)"
---

# ovh_cloud_project_containerregistry_user

Creates a user for a container registry associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_containerregistry" "registry" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

resource "ovh_cloud_project_containerregistry_user" "user" {
  service_name = ovh_cloud_project_containerregistry.registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.registry.id
  email        = "foo@bar.com"
  login        = "foobar"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `registry_id` - (Required) Registry ID
* `login` - (Required) User name
* `email` - (Required) User email

## Attributes Reference

The following attributes are exported:

* `email` - User email
* `id` - User ID
* `password` - (Sensitive) User password
* `user` - User name (same as `login`)

## Import

OVHcloud Managed Private Registry user can be imported using the `service_name`, `registry_id` and `id` of the user, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_containerregistry_user.my_user service_name/registry_id/user_id
```
