resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "<public cloud project ID>"
  name         = "my-container"
  region       = "GRA"
  type         = "PUBLIC"
  cors         = ["https://www.example.com"]
}
