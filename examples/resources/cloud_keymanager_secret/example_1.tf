resource "ovh_cloud_keymanager_secret" "secret" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region       = "GRA"
  name         = "my-secret"
  secret_type  = "OPAQUE"

  payload              = base64encode("my-secret-value")
  payload_content_type = "APPLICATION_OCTET_STREAM"

  metadata = {
    environment = "production"
  }
}
