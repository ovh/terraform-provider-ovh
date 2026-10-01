# A bootable volume, created from a Glance image.
resource "ovh_cloud_storage_block_volume" "boot" {
  service_name = "<Public cloud project id>"
  name         = "my-boot-volume"
  size         = 50
  region       = "GRA11"
  volume_type  = "HIGH_SPEED_GEN2"

  create_from = {
    image_id = "<image id>"
  }
}

# image_id is omitted: the instance boots from the volume instead. The
# instance's current_state.image stays null for a boot-from-volume instance.
resource "ovh_cloud_instance" "boot_from_volume" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-boot-from-volume-instance"
  flavor_id    = "<flavor id>"
  volume_ids   = [ovh_cloud_storage_block_volume.boot.id]

  networks = [
    { auto_assign_public_ip = true },
  ]
}
