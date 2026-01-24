resource "ovh_cloud_project_storage" "storage_with_lock" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA"
  name         = "my-locked-storage"

  object_lock = {
    status = "enabled"
    rule = {
      mode   = "governance"
      period = "P30D" # 30 days retention
    }
  }
}
