resource "ovh_cloud_key_manager_secret" "certificate" {
  service_name = "Public cloud project ID"
  region       = "GRA"
  name         = "my-certificate"
  secret_type  = "CERTIFICATE"
}

resource "ovh_cloud_key_manager_secret" "private_key" {
  service_name = "Public cloud project ID"
  region       = "GRA"
  name         = "my-private-key"
  secret_type  = "PRIVATE"
}

resource "ovh_cloud_key_manager_container" "container" {
  service_name = "Public cloud project ID"
  region       = "GRA"
  name         = "my-certificate-container"
  type         = "CERTIFICATE"

  secret_refs = [
    {
      name      = "certificate"
      secret_id = ovh_cloud_key_manager_secret.certificate.id
    },
    {
      name      = "private_key"
      secret_id = ovh_cloud_key_manager_secret.private_key.id
    }
  ]
}
