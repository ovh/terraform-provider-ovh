resource "ovh_cloud_managed_analytics_kafka_schemaregistryacl" "schema_registry_acl" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}
