---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_database_certificates (Data Source)

Use this data source to get the CA certificate of a managed database service (mysql, postgresql), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_database_certificates" "this" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  engine       = "postgresql"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
}
```

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `engine` - (Required) Engine of the service: mysql, postgresql.
* `service_id` - (Required) ID of the service.

## Attributes Reference

The following attributes are exported:

* `ca` - CA certificate, in PEM format, that verifies the connections to the service.
