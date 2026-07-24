resource "ovh_cloud_storage_file_share" "share" {
  service_name = "<Public cloud project id>"
  name         = "my-share"
  size         = 150
  region       = "GRA1"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
}

resource "ovh_cloud_file_storage_acl" "acl" {
  service_name = ovh_cloud_storage_file_share.share.service_name
  share_id     = ovh_cloud_storage_file_share.share.id
  access_to    = "10.0.0.0/24"
  access_level = "READ_WRITE"
}
