resource "ovh_cloud_keymanager_secret" "certificate" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region       = "GRA"
  name         = "my-certificate"
  secret_type  = "CERTIFICATE"
}

resource "ovh_cloud_keymanager_secret" "private_key" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region       = "GRA"
  name         = "my-private-key"
  secret_type  = "PRIVATE"
}

resource "ovh_cloud_keymanager_container" "container" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region       = "GRA"
  name         = "my-certificate-container"
  type         = "CERTIFICATE"

  secret_refs {
    name      = "certificate"
    secret_id = ovh_cloud_keymanager_secret.certificate.id
  }

  secret_refs {
    name      = "private_key"
    secret_id = ovh_cloud_keymanager_secret.private_key.id
  }
}
