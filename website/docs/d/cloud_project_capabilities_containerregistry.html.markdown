---
subcategory : "Managed Private Registry (MPR)"
---

# ovh_cloud_project_capabilities_containerregistry (Data Source)

Use this data source to get the container registry capabilities of a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_capabilities_containerregistry" "capabilities" {
  service_name = "XXXXXX"
}
```

## Argument Reference


* `service_name` - The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 


## Attributes Reference

`id` is set to the md5 sum of the list of all capabilities plans ids. In addition,
the following attributes are exported:

* `result` - List of container registry capability for a single region
   * `region_name` - The region name
   * `plans` - Available plans in the region
     * `code` - Plan code from the catalog
     * `created_at` - Plan creation date
     * `features` - Features of the plan
       * `vulnerability` - Vulnerability scanning
     * `id` - Plan ID
     * `name` - Plan name
     * `registry_limits` - Container registry limits
       * `image_storage` - Docker image storage limits in bytes
       * `parallel_request` - Parallel requests on Docker image API (/v2 Docker registry API)
     * `updated_at` - Plan last update date
