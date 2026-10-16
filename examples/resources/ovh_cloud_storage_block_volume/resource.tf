resource "ovh_cloud_storage_block_volume" "example" {
  service_name = "xxxxxxxxxx"
  name         = "my-volume"
  size         = 10
  region       = "GRA11"
  volume_type  = "CLASSIC"
}
