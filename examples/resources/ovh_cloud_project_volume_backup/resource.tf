resource "ovh_cloud_project_volume_backup" "backup" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA9"
  volume_id    = "<volume ID>"
  name         = "ExampleBackup"
}