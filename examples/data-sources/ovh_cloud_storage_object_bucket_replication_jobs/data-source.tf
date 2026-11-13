data "ovh_cloud_storage_object_bucket_replication_jobs" "jobs" {
  service_name = "<Public cloud project id>"
  bucket_id    = "GRA_my-source-bucket"
}

output "replication_job_statuses" {
  value = { for job in data.ovh_cloud_storage_object_bucket_replication_jobs.jobs.jobs : job.id => job.status }
}
