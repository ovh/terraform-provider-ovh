resource "ovh_cloud_instance" "instance" {
  service_name = <Public cloud project id>
  region       = "GRA11"
  name         = "my-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"
  ssh_key_name = "my-ssh-key"

  networks = [
    { public = true },
  ]
}
