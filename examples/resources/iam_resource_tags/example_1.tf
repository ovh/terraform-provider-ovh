# Basic example showing how to tag a cloud project
data "ovh_cloud_project" "my_project" {
  service_name = "01234567890123456798012345678901"
}

resource "ovh_iam_resource_tags" "project_tags" {
  resource_urn = data.ovh_cloud_project.my_project.iam.urn

  tags = {
    environment = "test"
    team        = "perso"
    managed_by  = "opentofu"
  }
}
