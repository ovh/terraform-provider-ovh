---
subcategory : "Managed Private Registry (MPR)"
---

# ovh_cloud_project_containerregistry (Data Source)

Use this data source to get information about a container registry associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_containerregistry" "my_registry" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}
```

## Argument Reference

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `registry_id` - Registry ID

## Attributes Reference

The following attributes are exported:

* `created_at` - Registry creation date
* `id` - Registry ID
* `name` - Registry name
* `project_id` - Project ID of your registry
* `region` - Region of the registry
* `size` - Current size of the registry (bytes)
* `status` - Registry status
* `updated_at` - Registry last update date
* `url` - Access url of the registry
* `version` - Version of your registry
