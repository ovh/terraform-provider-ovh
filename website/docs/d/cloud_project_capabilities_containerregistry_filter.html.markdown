---
subcategory : "Managed Private Registry"
---

# ovh_cloud_project_capabilities_containerregistry_filter (Data Source)

Use this data source to filter the list of container registry capabilities associated with a public cloud project to match one and only one capability.

## Example Usage

```hcl
data "ovh_cloud_project_capabilities_containerregistry_filter" "capability" {
  service_name = "XXXXXX"
  region       = "GRA"
  plan_name    = "SMALL"
}
```

## Argument Reference


* `service_name` - The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 
* `region` - The region name
* `plan_name` - The plan name. It can be 'SMALL', 'MEDIUM' or 'LARGE'.

## Attributes Reference

The following attributes are exported:

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
