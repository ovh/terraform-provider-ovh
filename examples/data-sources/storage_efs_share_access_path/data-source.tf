data "ovh_storage_efs" "efs" {
  service_name = "XXX"
}

data "ovh_storage_efs_share_access_path" "access_path" {
  service_name = data.ovh_storage_efs.efs.service_name
  share_id     = "XXX"
  id           = "XXX"
}