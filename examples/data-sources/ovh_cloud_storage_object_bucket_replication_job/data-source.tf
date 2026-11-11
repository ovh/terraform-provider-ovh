data "ovh_cloud_storage_object_bucket_replication_job" "job" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-source-bucket"
  id           = "<job id>"
}

output "replication_job_progress" {
  value = "${data.ovh_cloud_storage_object_bucket_replication_job.job.completed}/${data.ovh_cloud_storage_object_bucket_replication_job.job.total}"
}
