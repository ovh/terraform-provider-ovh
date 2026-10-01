resource "ovh_cloud_key_manager_secret" "secret" {
  service_name = "Public cloud project ID"
  region       = "GRA"
  name         = "my-secret"
  secret_type  = "OPAQUE"

  payload              = base64encode("my-secret-value")
  payload_content_type = "APPLICATION_OCTET_STREAM"

  metadata = {
    environment = "production"
  }
}
