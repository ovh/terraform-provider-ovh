# A service account only used through KMS certificates, its client secret is never stored
resource "ovh_me_api_oauth2_client" "kms_service_account" {
  name                  = "kms service account"
  description           = "Service account accessing the KMS with a certificate"
  flow                  = "CLIENT_CREDENTIALS"
  discard_client_secret = true
}

resource "ovh_okms_credential" "service_account_cred" {
  okms_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "service_account_cred"
  identity_urns = [ovh_me_api_oauth2_client.kms_service_account.identity]
  description   = "Credential for the KMS service account"
}

resource "local_sensitive_file" "service_account_private_key" {
  content         = ovh_okms_credential.service_account_cred.private_key_pem
  filename        = "${path.module}/service_account.key"
  file_permission = "0600"
}
