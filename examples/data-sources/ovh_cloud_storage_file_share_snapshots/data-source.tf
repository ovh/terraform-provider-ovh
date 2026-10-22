data "ovh_cloud_storage_file_share_snapshots" "snapshots" {
  service_name = "<public cloud project ID>"
}

data "ovh_cloud_storage_file_share_snapshots" "share_snapshots" {
  service_name = "<public cloud project ID>"
  share_id     = "<file share ID>"
}
