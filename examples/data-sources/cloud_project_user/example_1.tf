data "ovh_cloud_project_users" "project_users" {
  service_name = "XXX"
}

locals {
  # Get the user ID of a previously created user with the description "S3-User"
  users      = [for user in data.ovh_cloud_project_users.project_users.users : user.user_id if user.description == "S3-User"]
  s3_user_id = local.users[0]
}

data "ovh_cloud_project_user" "my_user" {
  service_name = data.ovh_cloud_project_users.project_users.service_name
  user_id      = local.s3_user_id
}
