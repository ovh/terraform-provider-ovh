resource "ovh_cloud_storage_file_share" "share" {
  service_name = "xxxxxxxxxx"
  name         = "my-share"
  size         = 100
  region       = "GRA1"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"
  network_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

resource "ovh_cloud_storage_file_share_snapshot" "snapshot" {
  service_name = ovh_cloud_storage_file_share.share.service_name
  region       = ovh_cloud_storage_file_share.share.region
  share_id     = ovh_cloud_storage_file_share.share.id
  name         = "my-snapshot"
  description  = "Daily snapshot of my-share"
}
