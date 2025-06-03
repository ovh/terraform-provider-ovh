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

resource "ovh_storage_efs_share_acl" "acl" {
  service_name = data.ovh_storage_efs.efs.service_name
  share_id     = ovh_storage_efs_share.share.id
  access_level = "ro"
  access_to    = "10.0.0.1"
}