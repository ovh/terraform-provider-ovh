resource "ovh_cloud_instance_group" "group" {
  service_name = <Public cloud project id>
  region       = "GRA11"
  name         = "my-instance-group"
  policy       = "ANTI_AFFINITY"
}
