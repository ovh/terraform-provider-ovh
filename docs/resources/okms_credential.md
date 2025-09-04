---
subcategory : "Key Management Service (KMS)"
---

# ovh_okms_credential (Resource)

Creates a credential for an OVHcloud KMS.

## Example Usage

```terraform
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
  certificate_type = "ECDSA"
}
```

## Argument Reference

### Required

- `identity_urns` (List of String) List of identity URNs associated with the credential (max 25)
- `name` (String) Name of the credential (max 50 characters)
- `okms_id` (String) ID of the KMS

### Optional

- `csr` (String) Certificate Signing Request. The CSR should be encoded in the PEM format. If this argument is not set, the server will generate a CSR for this credential, and the corresponding private key will be returned in the `private_key_pem` attribute.
- `description` (String) Description of the credential (max 200 characters)
- `validity` (Number) Validity in days (default: 365 days, max: 365 days)
- `certificate_type` (String) Type of the certificate key algorithm. Allowed values: `ECDSA`, `RSA`. Default to `ECDSA`. Changing forces a new credential.

## Attributes Reference

- `certificate_type` (String) Type of the certificate key algorithm (`ECDSA` or `RSA`).
- `certificate_pem` (String) Certificate PEM of the credential.
- `created_at` (String) Creation time of the credential
- `expired_at` (String) Expiration time of the credential
- `from_csr` (Boolean) Whether the credential was generated from a CSR
- `id` (String) ID of the credential
- `private_key_pem` (String, Sensitive) Private Key PEM of the credential if no CSR is provided
- `status` (String) Status of the credential
