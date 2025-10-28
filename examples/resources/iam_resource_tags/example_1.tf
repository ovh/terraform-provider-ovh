data "ovh_cloud_project" "my_project" {
  service_name = "<public cloud project ID>"
}

resource "ovh_iam_resource_tags" "project_tags" {
  urn = data.ovh_cloud_project.my_project.iam.urn

  tags = {
    environment = "production"
    team        = "platform"
    managed_by  = "terraform"
  }
}
