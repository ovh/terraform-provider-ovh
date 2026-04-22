data "ovh_cloud_storage_object_bucket_object_versions" "example" {
  service_name = "<public cloud project ID>"
  bucket_name  = "my-bucket"
  object_key   = "path/to/my-object.txt"
  limit        = 100
}
