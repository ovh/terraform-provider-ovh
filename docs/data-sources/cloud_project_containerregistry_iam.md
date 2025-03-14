---
subcategory : "Managed Private Registry (MPR)"
---

# ovh_cloud_project_containerregistry_iam (Data Source)

Use this data source to enable OVHcloud IAM on a Managed Private Registry.

## Example Usage

```hcl
data "ovh_cloud_project_containerregistry_iam" "my_iam" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

output "iam-enabled" {
  value = data.ovh_cloud_project_containerregistry_iam.my_iam.iam_enabled
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `registry_id` - The id of the Managed Private Registry.

## Attributes Reference

The following attributes are exported:

* `service_name` - The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `registry_id` - The ID of the Managed Private Registry.
* `iam_enabled` - OVHcloud IAM feature status.
