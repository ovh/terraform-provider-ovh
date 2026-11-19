---
subcategory : "Managed Databases"
---

# ovh_cloud_managed_analytics_kafka_user_access (Data Source)

Use this data source to get the client certificate and key of a user of a managed analytics service (kafka), using the Cloud Databases API v2.

## Example Usage

```terraform
data "ovh_cloud_managed_analytics_kafka_user_access" "app" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  service_id   = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
  user_id      = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
}

output "app_key" {
  value     = data.ovh_cloud_managed_analytics_kafka_user_access.app.key
  sensitive = true
}
```

The key is stored in the Terraform state.

## Argument Reference

* `service_name` - (Optional) ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `service_id` - (Required) ID of the service.
* `user_id` - (Required) ID of the user.

## Attributes Reference

The following attributes are exported:

* `engine` - Engine of the service: kafka.
* `cert` - Client certificate of the user, in PEM format, for the SSL authentication to Kafka.
* `key` - (Sensitive) Private key of the client certificate, in PEM format.
