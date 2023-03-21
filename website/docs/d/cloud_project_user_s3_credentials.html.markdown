---
subcategory : "Account Management"
---

# ovh_cloud_project_user_s3_credentials (Data Source)

Use this data source to retrieve the list of all the S3 access_key_id associated with a public cloud project's user.

## Example Usage

```hcl
data "ovh_cloud_project_user_s3_credentials" "my_s3_credentials" {
  service_name = "XXXXXX"
  user_id      = "1234"
}

output "access_key_ids" {
  value = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.access_key_ids
}
```

## Argument Reference

- `service_name` - (Required) The ID of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `user_id` - (Required) The ID of a public cloud project's user.

## Attributes Reference

- `access_key_ids` - The list of the Access Key ID associated with this user.
