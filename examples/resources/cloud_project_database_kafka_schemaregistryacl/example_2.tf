resource "ovh_cloud_project_database_kafka_schemaregistryacl" "schema_registry_acl" {
  # ...

  timeouts {
    create = "1h"
    delete = "45m"
  }
}
