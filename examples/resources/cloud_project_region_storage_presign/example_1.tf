resource "ovh_cloud_project_region_storage_presign" "presigned_url" {
  service_name = "xxxxxxxxxxxxxxxxx"
  region_name  = "GRA"
  name         = "s3-bucket-name"
  expire       = 3600
  method       = "GET"
  object       = "an-object-in-the-bucket"
}

output "presigned_url" {
  value = ovh_cloud_project_region_storage_presign.presigned_url.url
}
