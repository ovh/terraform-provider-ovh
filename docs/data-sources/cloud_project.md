---
subcategory : "Cloud Project"
---

# ovh_cloud_project

Get the details of a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project" "project" {
  service_name = "XXX"
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

## Attributes Reference

The following attributes are exported:

- `service_name` - ID of the public cloud project
- `project_name` - Project name
- `project_id` - Project ID
- `status` - Current status
- `unleash` - Project unleashed
- `plan_code` - Order plan code
- `order_id` - Project order ID
- `manual_quota` - Manual quota prevent automatic quota upgrade
- `expiration` - Expiration date of your project. After this date, your project will be deleted
- `description` - Description of your project
- `creation_date` - Project creation date
- `access` - Project access
- `iam` - IAM resource information
  - `urn` - URN of the private database, used when writing IAM policies
  - `display_name` - Resource display name
  - `id` - Unique identifier of the resource in the IAM
  - `tags` - Resource tags. Tags that were internally computed are prefixed with `ovh:`
