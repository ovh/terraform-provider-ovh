resource "ovh_cloud_project_instance_snapshot" "snapshot" {
  service_name  = "<public cloud project ID>"
  instance_id   = "<instance ID>"
  name          = "SnapshotExample"
}