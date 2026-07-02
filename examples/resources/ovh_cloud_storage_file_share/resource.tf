resource "ovh_cloud_storage_file_share" "share" {
  service_name = "<public cloud project ID>"
  name         = "my-share"
  size         = 150
  region       = "GRA1"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
  description  = "My NFS share"

  access_rules = [
    {
      access_to    = "10.0.0.0/24"
      access_level = "READ_WRITE"
    }
  ]
}
