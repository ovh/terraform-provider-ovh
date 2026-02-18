data "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
  service_name   = "<public cloud project ID>"
  region_name    = "GRA"
  container_name = "my-bucket"
}
