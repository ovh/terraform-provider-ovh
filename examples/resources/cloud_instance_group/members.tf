# ANTI_AFFINITY spreads the members across distinct hypervisors: use it for the
# replicas of a highly-available service, so a single host failure cannot take
# them all down.
resource "ovh_cloud_instance_group" "spread" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-anti-affinity-group"
  policy       = "ANTI_AFFINITY"
}

# AFFINITY packs the members onto the same hypervisor: use it for tightly
# coupled instances that benefit from staying on one host.
resource "ovh_cloud_instance_group" "packed" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-affinity-group"
  policy       = "AFFINITY"
}

# Membership is never managed by this resource: it is set once, at instance
# creation, through the instance's group_id.
resource "ovh_cloud_instance" "replica" {
  count = 3

  service_name = "<Public cloud project id>"
  region       = ovh_cloud_instance_group.spread.region
  name         = "my-replica-${count.index}"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  group_id     = ovh_cloud_instance_group.spread.id

  networks = [
    { auto_assign_public_ip = true },
  ]
}
