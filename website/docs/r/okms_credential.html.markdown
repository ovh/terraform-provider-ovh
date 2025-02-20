---
subcategory : "Key Management Service (KMS)"
---

# ovh_okms_credential (Resource)

Creates a credential for an OVHcloud KMS.

## Example Usage

```hcl
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

## Attributes Reference

- `certificate_pem` (String) Certificate PEM of the credential.
- `created_at` (String) Creation time of the credential
- `expired_at` (String) Expiration time of the credential
- `from_csr` (Boolean) Whether the credential was generated from a CSR
- `id` (String) ID of the credential
- `private_key_pem` (String, Sensitive) Private Key PEM of the credential if no CSR is provided
- `status` (String) Status of the credential
