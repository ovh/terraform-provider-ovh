resource "ovh_cloud_managed_kubernetes" "cluster" {
  service_name = "<Public cloud project id>"
  region       = "GRA11"
  plan         = "FREE"
  name         = "my-cluster"
  version      = "1.34"

  cni     = { type = "CANAL" }
  network = { type = "PUBLIC" }
}
