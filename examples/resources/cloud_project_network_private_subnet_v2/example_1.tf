resource "ovh_cloud_project_network_private_subnet_v2" "subnet" {
  service_name      = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  network_id        = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name              = "my_private_subnet"
  region            = "XXX1"
  dns_nameservers   = ["1.1.1.1"]
  cidr              = "192.168.168.0/24"
  dhcp              = true
  enable_gateway_ip = true
  use_default_public_dns_resolver = false
}
