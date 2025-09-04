data "ovh_me" "myaccount" {}

resource "ovh_okms_credential" "cred_no_csr" {
  okms_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "cred"
  identity_urns = ["urn:v1:eu:identity:account:${data.ovh_me.current_account.nichandle}"]
  description   = "Credential without CSR"
  certificate_type = "ECDSA"
}

# Write the generated certificate and private key to local files
provider "local" {}

resource "local_file" "credential_certificate" {
  content         = ovh_okms_credential.cred_no_csr.certificate_pem
  filename        = "${path.module}/certificate.pem"
  file_permission = "0600"
}

resource "local_sensitive_file" "credential_private_key" {
  content         = ovh_okms_credential.cred_no_csr.private_key_pem
  filename        = "${path.module}/private.key"
  file_permission = "0600"
}

resource "ovh_okms_credential" "cred_from_csr" {
  okms_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "cred_csr"
  identity_urns = ["urn:v1:eu:identity:account:${data.ovh_me.current_account.nichandle}"]
  csr           = file("cred.csr")
  description   = "Credential from CSR"
}
