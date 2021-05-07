---
layout: "ovh"
page_title: "OVH: cloud_project_containerregistry"
sidebar_current: "docs-ovh-resource-cloud-project-containerregistry-x"
description: |-
  Creates a container registry associated with a public cloud project.
---

# ovh_cloud_project_containerregistry

Creates a container registry associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_capabilities_containerregistry_filter" "regcap" {
  service_name = "XXXXXX"
  plan_name    = "SMALL"
  region       = "GRA"
}

resource "ovh_cloud_project_containerregistry" "reg" {
  service_name = data.ovh_cloud_project_capabilities_containerregistry_filter.regcap.service_name
  plan_id      = data.ovh_cloud_project_capabilities_containerregistry_filter.regcap.id
  region       = data.ovh_cloud_project_capabilities_containerregistry_filter.regcap.region
  name         = "mydockerregistry"
}
```

## Argument Reference


* `service_name` - The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 

* `name` - Registry name
* `region` - Region of the registry
* `plan_id` - Plan ID of the registry


## Attributes Reference

The following attributes are exported:

* `created_at` - Registry creation date
* `id` - Registry ID
* `name` - Registry name
* `plan` -  Plan of the registry
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
* `project_id` - Project ID of your registry
* `region` - Region of the registry
* `size` - Current size of the registry (bytes)
* `status` - Registry status
* `updated_at` - Registry last update date
* `url` - Access url of the registry
* `version` - Version of your registry
