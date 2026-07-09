resource "ovh_cloud_storage_file_share" "share" {
  service_name = "<public cloud project ID>"
  name         = "my-share"
  size         = 150
  region       = "GRA1"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
}

resource "ovh_cloud_storage_file_share_snapshot" "snapshot" {
  service_name = ovh_cloud_storage_file_share.share.service_name
  share_id     = ovh_cloud_storage_file_share.share.id
  name         = "my-snapshot"
  description  = "Daily snapshot of my-share"
}
