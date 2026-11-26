resource "ovh_cloud_managed_kubernetes" "private" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  plan         = "STANDARD"
  name         = "my-private-cluster"

  cni = {
    type = "CILIUM"
    cilium = {
      hubble = { enabled = true }
    }
  }

  network = {
    type = "VRACK"
    vrack = {
      private_network_id = "<private network id>"
      subnet_id          = "<subnet id>"
      default_egress     = { type = "MANAGED_GATEWAY" }
    }
  }

  api_server = {
    access = {
      authorized_ips = {
        rules = [{ cidr = "203.0.113.0/24" }]
      }
    }
  }
}
