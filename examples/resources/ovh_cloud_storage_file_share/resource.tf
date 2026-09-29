resource "ovh_cloud_storage_file_share" "share" {
  service_name     = "<public cloud project ID>"
  name             = "my-share"
  size             = 150
  region           = "GRA1"
  protocol         = "NFS"
  share_type       = "STANDARD_1AZ"
  share_network_id = "<share network id>"
  description      = "My NFS share"
}
