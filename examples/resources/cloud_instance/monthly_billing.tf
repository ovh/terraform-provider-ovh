# Billed monthly. Monthly billing cannot be undone: once MONTHLY, the API
# refuses HOURLY at apply. Leaving billing out keeps the current period on an
# in-place update, but a replacement creates a new instance billed HOURLY.
resource "ovh_cloud_instance" "monthly" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  name         = "my-monthly-instance"
  flavor_id    = "<flavor id>"
  image_id     = "<image id>"

  billing = {
    period = "MONTHLY"
  }

  networks = [
    { auto_assign_public_ip = true },
  ]
}
