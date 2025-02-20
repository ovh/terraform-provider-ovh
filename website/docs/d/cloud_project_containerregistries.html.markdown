---
subcategory : "Managed Private Registry (MPR)"
---

# ovh_cloud_project_containerregistries (Data Source)

Use this data source to get the container registries of a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_containerregistries" "registries" {
  service_name = "XXXXXX"
}
```

## Argument Reference


* `service_name` - The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 


## Attributes Reference

`id` is set to the md5 sum of the list of all registries ids. In addition,
the following attributes are exported:

* `result` - The list of container registries associated with the project.
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
