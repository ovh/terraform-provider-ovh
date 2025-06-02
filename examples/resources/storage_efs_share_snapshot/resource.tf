data "ovh_storage_efs" "efs" {
  service_name = "XXX"
}

# This resource will destroy (at least) 10 seconds after ovh_storage_efs_share_snapshot resource
resource "ovh_storage_efs_share" "share" {
  service_name = data.ovh_storage_efs.efs.service_name
  name         = "share"
  description  = "My share"
  protocol     = "NFS"
  size         = 100
}

resource "ovh_storage_efs_share_snapshot" "snapshot" {
  depends_on = [time_sleep.wait_10_seconds]

  service_name = data.ovh_storage_efs.efs.service_name
  share_id     = ovh_storage_efs_share.share.id
  name         = "snapshot"
  description  = "My snapshot"
}

# This resource adds a delay between ovh_storage_efs_share_snapshot and ovh_storage_efs_share resources destruction
resource "time_sleep" "wait_10_seconds" {
  depends_on = [ovh_storage_efs_share.share]

  destroy_duration = "10s"
}