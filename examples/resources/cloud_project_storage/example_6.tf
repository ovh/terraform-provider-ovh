resource "ovh_cloud_project_storage" "storage_with_replication" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA"
  name         = "my-replicated-storage"

  versioning = {
    status = "enabled"
  }

  replication = {
    rules = [
      {
        id                        = "replicate-all"
        priority                  = 1
        status                    = "enabled"
        delete_marker_replication = "disabled"
        destination = {
          name   = "my-replicated-storage-replica"
          region = "GRA"
        }
      }
    ]
  }
}
