---
subcategory : "Cloud Project"
---

# ovh_cloud_projects

Get the details of your public cloud projects.

## Example Usage

```terraform
data "ovh_cloud_projects" "projects" {}
```

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
