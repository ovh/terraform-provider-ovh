# Every shape a `networks[]` entry can take. Entries keep the order they are
# written in: the API returns them sorted by network id and the provider
# re-orders them back to the configuration.
resource "ovh_cloud_instance" "networks" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-multi-homed-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"

  networks = [
    # 1. Public interface with a public IP assigned by the platform.
    #    At most one entry may set auto_assign_public_ip.
    { auto_assign_public_ip = true },

    # 2. A public IP the project already owns: an additional IP, or an Ext-Net
    #    IP of the project in the instance's region. Several such entries are
    #    allowed, and they may coexist with auto_assign_public_ip.
    { ip = "203.0.113.10" },

    # 3. Private interface with an address picked by IPAM.
    {
      network_id = "<private network id>"
      subnet_id  = "<subnet id>"
    },

    # 4. Private interface with an explicit address: the port's fixed IP is
    #    pinned to `ip` when it falls inside the subnet CIDR, otherwise the
    #    existing floating IP with that address is associated.
    {
      network_id = "<private network id>"
      subnet_id  = "<subnet id>"
      ip         = "10.0.0.42"
    },
  ]
}
