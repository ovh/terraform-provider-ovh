resource "ovh_cloud_project_workflow_backup" "my_backup" {
  service_name        = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region_name         = "GRA11"
  cron                = "50 4 * * *"
  instance_id         = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
  max_execution_count = "0"
  name                = "Backup workflow for instance"
  rotation            = "7"
}
