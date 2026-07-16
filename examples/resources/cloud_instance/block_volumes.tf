resource "ovh_cloud_storage_block_volume" "data" {
  service_name = "<Public cloud project id>"
  name         = "my-data-volume"
  size         = 100
  region       = "GRA11"
  volume_type  = "CLASSIC"
}

resource "ovh_cloud_storage_block_volume" "logs" {
  service_name = "<Public cloud project id>"
  name         = "my-logs-volume"
  size         = 20
  region       = "GRA11"
  volume_type  = "CLASSIC"
}

# volume_ids is mutable: adding or removing an id attaches or detaches the
# volume in place, without recreating the instance.
resource "ovh_cloud_instance" "with_volumes" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance-with-volumes"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"

  volume_ids = [
    ovh_cloud_storage_block_volume.data.id,
    ovh_cloud_storage_block_volume.logs.id,
  ]

  networks = [
    { auto_assign_public_ip = true },
  ]
}
