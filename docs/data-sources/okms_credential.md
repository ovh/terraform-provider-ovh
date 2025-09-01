---
subcategory: "Key Management Service (KMS)"
---

# ovh_okms_credential (Data Source)

Use this data source to retrieve data associated with a KMS credential, such as the PEM encoded certificate.

## Example Usage

```terraform
data "ovh_okms_resource" "kms" {
  okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `id` (String) ID of the credential
- `okms_id` (String) ID of the KMS

## Attributes Reference

- `certificate_type` (String) Type of the certificate (ECDSA or RSA)
- `certificate_pem` (String) PEM encoded certificate of the credential
- `created_at` (String) Creation time of the credential
- `description` (String) Description of the credential
- `expired_at` (String) Expiration time of the credential
- `from_csr` (Boolean) Is the credential generated from CSR
- `identity_urns` (List of String) List of identity URNs associated with the credential
- `name` (String) Name of the credential
- `status` (String) Status of the credential
