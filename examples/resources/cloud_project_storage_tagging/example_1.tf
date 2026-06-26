resource "ovh_cloud_project_storage_tagging" "tags" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA"
  name         = "my-existing-bucket"

  tags = {
    environment = "production"
    team        = "platform"
    managed-by  = "terraform"
  }
}
