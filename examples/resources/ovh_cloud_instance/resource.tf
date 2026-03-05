resource "ovh_cloud_instance" "example" {
  service_name = "xxxxxxxxxx"
  name         = "my-instance"
  flavor_id    = "068a57cf-b1b4-428f-9b17-4f32a526390c"
  image_id     = "8d75e170-1ef9-4e25-8fc8-d231929e56e8"
  region       = "GRA1"

  networks {
    id = "fbdf6240-8b56-4626-b57c-25e4af487606"
  }

  volume_ids = [
    "3809b388-066f-4c79-9d77-881a7cd19629"
  ]
}
