# Billed monthly. Monthly billing cannot be undone: once MONTHLY, the plan
# refuses HOURLY, and leaving billing out keeps the current period.
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
