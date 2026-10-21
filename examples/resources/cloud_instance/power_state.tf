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

# SUSPENDED: the memory state is saved to disk and the VM stops, but its
# resources stay allocated on the hypervisor. Switching power_state back to
# ACTIVE resumes it. PAUSED works the same way but keeps the memory state in
# the host's RAM.
resource "ovh_cloud_instance" "suspended" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-suspended-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  power_state  = "SUSPENDED"

  networks = [
    { auto_assign_public_ip = true },
  ]
}
