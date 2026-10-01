data "ovh_dedicated_cloud" "pcc" {
  service_name = "pcc-1-2-3-4"
}

resource "ovh_dedicated_cloud_user" "john" {
  service_name = data.ovh_dedicated_cloud.pcc.service_name
  name         = "john.doe"
  email        = "john.doe@example.com"
  first_name   = "John"
  last_name    = "Doe"
}
