# user_data carries a cloud-init payload, base64-encoded (standard encoding,
# 65535 bytes max). It is write-only: the API accepts it but never returns it, so
# Terraform state holds the only copy and a change made outside Terraform stays
# invisible.
#
# WARNING: changing or clearing user_data reinstalls the instance — the API
# rebuilds it on its current image and the root disk is wiped. The instance id,
# ports and IPs survive (it is an in-place update, not a replacement), but
# anything written to the root disk since boot is lost.

# From a file next to the configuration.
resource "ovh_cloud_instance" "from_file" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance-with-cloud-init"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  user_data    = base64encode(file("${path.module}/cloud-init.yaml"))

  networks = [
    { auto_assign_public_ip = true },
  ]
}

# Inline, with a heredoc.
resource "ovh_cloud_instance" "inline" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-instance-with-inline-cloud-init"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"

  user_data = base64encode(<<-EOT
    #cloud-config
    packages:
      - nginx
    runcmd:
      - [systemctl, enable, --now, nginx]
  EOT
  )

  networks = [
    { auto_assign_public_ip = true },
  ]
}
