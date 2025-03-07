data "ovh_me" "myaccount" {}

resource "ovh_okms_credential" "cred_no_csr" {
  okms_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "cred"
  identity_urns = ["urn:v1:eu:identity:account:${data.ovh_me.current_account.nichandle}"]
  description   = "Credential without CSR"
}

resource "ovh_okms_credential" "cred_from_csr" {
  okms_id       = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name          = "cred_csr"
  identity_urns = ["urn:v1:eu:identity:account:${data.ovh_me.current_account.nichandle}"]
  csr           = file("cred.csr")
  description   = "Credential from CSR"
}
