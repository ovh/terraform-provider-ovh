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

# image_id is omitted and boot_from_volume_id names the volume: the instance
# boots from it instead. The boot volume can only be chosen at creation, cannot
# be detached, and is kept when the instance is deleted. The instance's
# current_state.image stays null for a boot-from-volume instance.
#
# To give each volume a role, use volumes instead:
#   volumes = [{ id = ovh_cloud_storage_block_volume.boot.id, role = "BOOT" }]
resource "ovh_cloud_instance" "boot_from_volume" {
  service_name        = "<Public cloud project id>"
  region              = "GRA11"
  name                = "my-boot-from-volume-instance"
  flavor_id           = "<flavor id>"
  boot_from_volume_id = ovh_cloud_storage_block_volume.boot.id

  networks = [
    { auto_assign_public_ip = true },
  ]
}
