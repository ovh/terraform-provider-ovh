resource "ovh_cloud_security_group" "web" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-web-security-group"
  description  = "Allow SSH and HTTPS"

  rule = [
    {
      direction        = "INGRESS"
      ethernet_type    = "IPV4"
      protocol         = "TCP"
      port_range_min   = 22
      port_range_max   = 22
      remote_ip_prefix = "0.0.0.0/0"
      description      = "SSH"
    },
    {
      direction        = "INGRESS"
      ethernet_type    = "IPV4"
      protocol         = "TCP"
      port_range_min   = 443
      port_range_max   = 443
      remote_ip_prefix = "0.0.0.0/0"
      description      = "HTTPS"
    },
  ]
}

# Exactly these groups are applied to every interface of the instance.
resource "ovh_cloud_instance" "filtered" {
  service_name       = "<Public cloud project id>"
  region             = "GRA11"
  name               = "my-filtered-instance"
  flavor_id          = "<flavor id>"
  image_id           = "<image id>"
  security_group_ids = [ovh_cloud_security_group.web.id]

  networks = [
    { auto_assign_public_ip = true },
  ]
}

# An explicit empty list applies no security group at all: the instance accepts
# no inbound traffic. Omit the attribute instead to let the platform apply the
# project's `default` security group.
resource "ovh_cloud_instance" "no_filtering" {
  service_name       = "<Public cloud project id>"
  region             = "GRA11"
  name               = "my-unfiltered-instance"
  flavor_id          = "<flavor id>"
  image_id           = "<image id>"
  security_group_ids = []

  networks = [
    { auto_assign_public_ip = true },
  ]
}
