data "ovh_storage_efs" "efs" {
  service_name = "XXX"
}

resource "ovh_storage_efs_share" "share" {
  service_name = data.ovh_storage_efs.efs.service_name
  name         = "share"
  description  = "My share"
  protocol     = "NFS"
  size         = 100
}