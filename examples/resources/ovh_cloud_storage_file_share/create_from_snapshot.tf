resource "ovh_cloud_storage_file_share_snapshot" "snapshot" {
  service_name = ovh_cloud_storage_file_share.share.service_name
  share_id     = ovh_cloud_storage_file_share.share.id
  name         = "my-snapshot"
}

resource "ovh_cloud_storage_file_share" "share_from_snapshot" {
  service_name = ovh_cloud_storage_file_share.share.service_name
  name         = "my-share-from-snapshot"
  region       = "GRA1"
  protocol     = "NFS"

  create_from = {
    snapshot_id = ovh_cloud_storage_file_share_snapshot.snapshot.id
  }
}
