data "ovh_storage_efs" "efs" {
  service_name = "XXX"
}

data "ovh_storage_efs_share_access_paths" "access_paths" {
  service_name = data.ovh_storage_efs.efs.service_name
  share_id     = "XXX"
}