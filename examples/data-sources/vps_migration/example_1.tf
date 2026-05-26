data "ovh_vps_migration" "migration" {
  service_name = "vpsXXXXX.ovh.net"
}

output "migration_date" {
  value = data.ovh_vps_migration.migration.date
}
