resource "ovh_domain_name_servers" "name_servers" {
  domain = "mydomain.ovh"

  servers {
    host = "dns105.ovh.net"
    ip = "213.251.188.144"
  }

  servers {
    host = "ns105.ovh.net"
  }
}
