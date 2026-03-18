resource "ovh_cloud_storage_block_volume" "example" {
  service_name = "xxxxxxxxxx"
  name         = "my-volume"
  size         = 10
  region       = "GRA1"
  volume_type = "STANDARD"
}
