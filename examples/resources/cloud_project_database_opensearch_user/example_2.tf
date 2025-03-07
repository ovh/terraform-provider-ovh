data "ovh_cloud_project_database" "opensearch" {
  service_name  = "XXX"
  engine        = "opensearch"
  id            = "ZZZ"
}

# Change password_reset with the datetime each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_opensearch_user" "userDatetime" {
  service_name    = data.ovh_cloud_project_database.opensearch.service_name
  cluster_id      = data.ovh_cloud_project_database.opensearch.id
  acls {
    pattern    = "logs_*"
    permission = "read"
  }
  acls {
    pattern    = "data_*"
    permission = "deny"
  }
  name            = "alice"
  password_reset  = "2024-01-02T11:00:00Z"
}

variable "something" {
  type = string
}

# Set password_reset to be based on the update of another variable to reset the password
resource "ovh_cloud_project_database_opensearch_user" "userMd5" {
  service_name    = data.ovh_cloud_project_database.opensearch.service_name
  cluster_id      = data.ovh_cloud_project_database.opensearch.id
  acls {
    pattern    = "logs_*"
    permission = "read"
  }
  acls {
    pattern    = "data_*"
    permission = "deny"
  }
  name            = "bob"
  password_reset  = md5(var.something)
}

# Change password_reset each time you want to reset the password to trigger an update
resource "ovh_cloud_project_database_opensearch_user" "user" {
  service_name    = data.ovh_cloud_project_database.opensearch.service_name
  cluster_id      = data.ovh_cloud_project_database.opensearch.id
  acls {
    pattern    = "logs_*"
    permission = "read"
  }
  acls {
    pattern    = "data_*"
    permission = "deny"
  }
  name            = "johndoe"
  password_reset  = "reset1"
}

output "user_password" {
  value     = ovh_cloud_project_database_opensearch_user.user.password
  sensitive = true
}
