resource "ovh_cloud_project_storage" "storage_with_tags" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA"
  name         = "my-tagged-storage"

  tags = {
    environment = "production"
    team        = "platform"
  }
}
