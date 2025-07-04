---
subcategory : "Object Storage"
---

# ovh_cloud_project_user_s3_credential

Creates an S3 Credential for a user in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_user" "user" {
  service_name = "XXX"
  description  = "my user for acceptance tests"
  role_names   = [
    "objectstore_operator"
  ]
}

resource "ovh_cloud_project_user_s3_credential" "my_s3_credentials" {
  service_name = ovh_cloud_project_user.user.service_name
  user_id      = ovh_cloud_project_user.user.id
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `user_id` - (Required) The ID of a public cloud project's user.

## Attributes Reference

The following attributes are exported:

- `service_name` - See Argument Reference above.
- `user_id` - See Argument Reference above.
- `access_key_id` - the Access Key ID
- `secret_access_key` - (Sensitive) the Secret Access Key

## Import

OVHcloud User S3 Credentials can be imported using the `service_name`, `user_id` and `access_key_id` of the credential, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_user_s3_credential.s3_credential service_name/user_id/access_key_id
```
