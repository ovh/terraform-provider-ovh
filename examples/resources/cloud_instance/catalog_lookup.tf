data "ovh_cloud_instance_flavors" "catalog" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
}

data "ovh_cloud_instance_images" "catalog" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
}

resource "ovh_cloud_ssh_key" "deploy" {
  service_name = "<Public cloud project id>"
  name         = "my-deploy-key"
  public_key   = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample user@host"
}

# Resolve flavor_id and image_id from the per-region catalog instead of
# hardcoding opaque UUIDs.
resource "ovh_cloud_instance" "from_catalog" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"

  name = "my-catalog-instance"

  flavor_id = one([
    for flavor in data.ovh_cloud_instance_flavors.catalog.flavors :
    flavor.id if flavor.name == "b3-8"
  ])

  image_id = one([
    for image in data.ovh_cloud_instance_images.catalog.images :
    image.id if image.name == "Debian 12"
  ])

  ssh_key_name = ovh_cloud_ssh_key.deploy.name

  networks = [
    { auto_assign_public_ip = true },
  ]
}
