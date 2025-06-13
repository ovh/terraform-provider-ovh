data "ovh_vps" "my_vps" {
  service_name = "vps-xxxxxxxxxx.vps.ovh.net"
}

resource "ovh_cloud_project_ssh_key" "key" {
  service_name = "xxxxxxxxxx"
  public_key   = "ssh-rsa AAAAB3N..."
  name       = "my_ssh_key"
}

resource "ovh_vps_reinstall" "vps_reinstall" {
  service_name = data.ovh_vps.my_vps.service_name
  # Debian 12
  image_id       = "45b2f222-ab10-44ed-863f-720942762b6f"
  public_ssh_key = ovh_cloud_project_ssh_key.key.public_key
}