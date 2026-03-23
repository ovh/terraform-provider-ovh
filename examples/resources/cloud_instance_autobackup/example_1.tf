resource "ovh_cloud_instance_autobackup" "backup" {
  project_id  = "<public cloud project ID>"
  name        = "my-daily-backup"
  image_name  = "my-instance-backup"
  cron        = "0 2 * * *"
  rotation    = 7
  instance_id = "<instance ID>"
  region      = "GRA11"
}

output "next_execution" {
  value = ovh_cloud_instance_autobackup.backup.current_state.next_execution_time
}
