resource "ovh_domain_name" "domain" {
  domain_name = "example.com"

  target_spec = {
    dns_configuration = {
      name_servers = [
        {
          name_server = "dns101.ovh.net"
        },
        {
          name_server = "ns101.ovh.net"
        }
      ]
    }
  }
}
