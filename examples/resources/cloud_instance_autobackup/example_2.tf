resource "ovh_cloud_instance_autobackup" "backup_cross_region" {
  project_id  = "<public cloud project ID>"
  name        = "my-cross-region-backup"
  image_name  = "my-instance-backup"
  cron        = "0 3 * * 0"
  rotation    = 4
  instance_id = "<instance ID>"
  region      = "GRA11"

  distant = {
    region     = "SBG5"
    image_name = "my-instance-backup-remote"
  }
}
