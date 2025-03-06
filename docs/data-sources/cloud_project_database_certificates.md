---
subcategory : "Managed Databases"
---

# ovh_cloud_project_database_certificates (Data Source)

Use this data source to get information about certificates of a cluster associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_database_certificates" "certificates" {
  service_name  = "XXX"
  engine        = "YYY"
  cluster_id    = "ZZZ"
}

output "certificates_ca" {
  value = data.ovh_cloud_project_database_certificates.certificates.ca
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `engine` - (Required) The engine of the database cluster you want database information. To get a full list of available engine visit: [public documentation](https://docs.ovh.com/gb/en/publiccloud/databases). Available engines:
  * `cassandra`
  * `kafka`
  * `mysql`
  * `postgresql`

* `cluster_id` - (Required) Cluster ID

## Attributes Reference

The following attributes are exported:

`id` is set to the md5 sum of the CA. In addition, the following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `ca` - CA certificate used for the service.
