---
subcategory : "Managed Private Registry (MPR)"
---

# ovh_cloud_project_containerregistry_user

Creates a user for a container registry associated with a public cloud project.

## Example Usage

```hcl
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


* `service_name` - The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 

* `registry_id` - Registry ID

## Attributes Reference

The following attributes are exported:

* `email` - User email
* `id` - User ID
* `password` - (Sensitive) User password
* `user` - User name
