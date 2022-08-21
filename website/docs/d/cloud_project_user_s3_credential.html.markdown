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
data "ovh_cloud_project_user_s3_credentials" "s3_credentials" {
   service_name = "XXXXXX"
   user_id      = "1234"
}

data "ovh_cloud_project_user_s3_credential" "s3_cred_1" {
   service_name  = "XXXXXX"
   user_id       = "1234"
   access_key_id = data.ovh_cloud_project_user_s3_credentials.s3_credentials.access_key_ids[0]
}
```

## Argument Reference

- `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

- `user_id` - (Required) The id of a public cloud project's user.

- `access_key_id` - the Access Key ID

## Attributes Reference

- `secret_access_key` - (Sensitive) the Secret Access Key
