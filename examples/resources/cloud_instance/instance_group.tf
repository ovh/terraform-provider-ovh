resource "ovh_cloud_instance_group" "spread" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-anti-affinity-group"
  policy       = "ANTI_AFFINITY"
}

# Group membership is only settable here, through the instance's group_id. It
# is immutable: moving an instance to another group, or out of its group,
# recreates the instance.
resource "ovh_cloud_instance" "grouped" {
  count = 2

  service_name = "<Public cloud project id>"
  region       = ovh_cloud_instance_group.spread.region
  name         = "my-grouped-instance-${count.index}"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  group_id     = ovh_cloud_instance_group.spread.id

  networks = [
    { auto_assign_public_ip = true },
  ]
}
