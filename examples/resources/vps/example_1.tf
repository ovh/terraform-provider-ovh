resource "ovh_vps" "web" {
  plan_code    = "vps-2025-model1"
  display_name = "web"

  location = {
    region = "eu-west-gra"
  }

  installation = {
    from = {
      image_name = "Debian 12"
    }
    ssh_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample user@laptop"
  }
}

output "vps_ipv4" {
  value = [for address in ovh_vps.web.current_state.ip_addresses : address.ip if address.version == 4]
}
