resource "ovh_cloud_storage_file_share" "share" {
  service_name     = "<Public cloud project id>"
  name             = "my-share"
  size             = 150
  region           = "GRA11"
  protocol         = "NFS"
  share_type       = "STANDARD_1AZ"
  share_network_id = "<share network id>"
}

# Attaching a share creates a Manila access rule granting the instance's fixed
# IPv4 address on the share's network. Omit access_level to let the API apply
# its default (READ_WRITE).
resource "ovh_cloud_instance" "with_shares" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance-with-shares"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"

  shares = [
    {
      id = ovh_cloud_storage_file_share.share.id
    },
    {
      id           = "<read only share id>"
      access_level = "READ_ONLY"
    },
  ]

  networks = [
    {
      network_id = "<private network id>"
      subnet_id  = "<subnet id>"
    },
  ]
}
