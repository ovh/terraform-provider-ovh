---
layout: "ovh"
page_title: "OVH: cloud_project_user_s3_credential"
sidebar_current: "docs-ovh-datasource-cloud-project-user-s3-credential"
description: |-
  Get the Secret Access Key of an Access Key ID of a public cloud project's user.
---

# ovh_cloud_project_user_s3_credential (Data Source)

Use this data source to retrieve the Secret Access Key of an Access Key ID associated with a public cloud project's user.

## Example Usage

```hcl
data "ovh_cloud_project_user_s3_credentials" "my_s3_credentials" {
   service_name = "XXXXXX"
   user_id      = "1234"
}

data "ovh_cloud_project_user_s3_credential" "my_s3_credential" {
   service_name  = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.service_name
   user_id       = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.user_id
   access_key_id = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.access_key_ids[0]
}

output "my_access_key_id" {
   value = ovh_cloud_project_user_s3_credential.my_s3_credential.access_key_id
}

output "my_secret_access_key" {
   value     = ovh_cloud_project_user_s3_credential.my_s3_credential.secret_access_key
   sensitive = true
}
```

## Argument Reference

- `service_name` - (Required) The ID of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `user_id` - (Required) The ID of a public cloud project's user.

- `access_key_id` - the Access Key ID

## Attributes Reference

- `secret_access_key` - (Sensitive) the Secret Access Key
