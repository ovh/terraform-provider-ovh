---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_opensearch_user (Data Source)

Use this data source to get information about a user of a opensearch cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics_opensearch_user" "os_user" {
  service_name  = "XXX"
  cluster_id    = "YYY"
  name          = "ZZZ"
}

output "os_user_acls" {
  value = data.ovh_cloud_managed_analytics_opensearch_user.os_user.acls
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

* `name` - (Required) Name of the user.

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `acls` - Acls of the user.
  * `pattern` - Pattern of the ACL.
  * `permission` - Permission of the ACL.
* `service_name` - Current status of the user.
* `status` - Current status of the user.
* `name` - Name of the user.
