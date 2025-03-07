data "ovh_cloud_project_users" "project_users" {
  service_name = "XXX"
}

locals {
  # Get the user ID of a previously created user with the description "S3-User"
  users      = [for user in data.ovh_cloud_project_users.project_users.users : user.user_id if user.description == "S3-User"]
  s3_user_id = local.users[0]
}

data "ovh_cloud_project_user_s3_credentials" "my_s3_credentials" {
  service_name = data.ovh_cloud_project_users.project_users.service_name
  user_id      = local.s3_user_id
}

data "ovh_cloud_project_user_s3_credential" "my_s3_credential" {
  service_name  = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.service_name
  user_id       = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.user_id
  access_key_id = data.ovh_cloud_project_user_s3_credentials.my_s3_credentials.access_key_ids[0]
}

output "my_access_key_id" {
  value = data.ovh_cloud_project_user_s3_credential.my_s3_credential.access_key_id
}

output "my_secret_access_key" {
  value     = data.ovh_cloud_project_user_s3_credential.my_s3_credential.secret_access_key
  sensitive = true
}
