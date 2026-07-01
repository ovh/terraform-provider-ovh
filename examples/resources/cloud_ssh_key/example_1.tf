resource "ovh_cloud_ssh_key" "my_key" {
  service_name = "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  name         = "my-deploy-key"
  public_key   = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample user@host"
}
