data "ovh_cloud_project_user_s3_credentials" "my_s3_credentials" {
  service_name = "XXXXXX"
  user_id      = "1234"
}

output "access_key_ids" {
  value = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.access_key_ids
}
