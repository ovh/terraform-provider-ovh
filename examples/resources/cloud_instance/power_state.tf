# SHUTOFF: the instance stays provisioned and keeps its disks, ports and IPs,
# but the VM is powered off.
resource "ovh_cloud_instance" "stopped" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-stopped-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  power_state  = "SHUTOFF"

  networks = [
    { auto_assign_public_ip = true },
  ]
}

# SHELVED: the instance is powered off and its resources are released on the
# hypervisor. Switching power_state back to ACTIVE unshelves it in place.
resource "ovh_cloud_instance" "shelved" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-shelved-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  power_state  = "SHELVED"

  networks = [
    { auto_assign_public_ip = true },
  ]
}
