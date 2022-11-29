---
layout: "ovh"
page_title: "OVH: cloud_project_region_storage_presign"
sidebar_current: "docs-ovh-resource-cloud-project-region-storage-presign"
description: |-
  Generates a temporary presigned S3 URLs to download or upload an object.
---

# ovh_cloud_project_region_storage_presign

Generates a temporary presigned S3 URLs to download or upload an object.

-> __NOTE__ This resource is only compatible with the `High Performance - S3` solution for object storage.

## Example Usage

```hcl
resource "ovh_cloud_project_region_storage_presign" "presigned_url" {
  service_name = "xxxxxxxxxxxxxxxxx"
  region_name  = "GRA"
  name         = "s3-bucket-name"
  expire       = 3600
  method       = "GET"
  object       = "an-object-in-the-bucket"
}

output "presigned_url" {
  value = ovh_cloud_project_region_storage_presign.presigned_url.url
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `region_name` - (Required) The region in which your storage is located.
  Ex.: "GRA".
* `name` - (Required) The name of your S3 storage container/bucket.
* `expire` - (Required) Define, in seconds, for how long your URL will be valid.
* `method` - (Required) The method you want to use to interact with your object. Can be either 'GET' or 'PUT'.
* `object` - (Required) The name of the object in your S3 bucket.


## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `region_name` - See Argument Reference above.
* `name` - See Argument Reference above.
* `expire` - See Argument Reference above.
* `method` - See Argument Reference above.
* `object` - See Argument Reference above.
* `url` - Computed URL result.
