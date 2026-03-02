---
subcategory : "Managed Databases"
---

~> **DEPRECATED:** Use `ovh_cloud_managed_database_mongodb_prometheus` instead. This data source will be removed in the next major version.

# ovh_cloud_project_database_mongodb_prometheus (Data Source)

Use this data source to get information about a prometheus of a MongoDB cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_database_mongodb_prometheus" "prometheus" {
  service_name  = "XXX"
  cluster_id    = "ZZZ"
}

output "name" {
  value = data.ovh_cloud_project_database_mongodb_prometheus.prometheus.username
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `id` - Cluster ID.
* `service_name` - See Argument Reference above.
* `username` - name of the prometheus user.
* `srv_domain` - Name of the srv domain endpoint.
