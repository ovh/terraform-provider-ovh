---
subcategory : "KMS"
---

# ovh_okms_service_key (Resource)

Creates a Service Key in an OVHcloud KMS.

## Example Usage

```hcl
resource "ovh_okms_service_key" "key_symetric" {
  okms_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "key_oct"
  type       = "oct"
  size       = 256
  operations = ["encrypt", "decrypt"]
}

resource "ovh_okms_service_key" "key_rsa" {
  okms_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "key_rsa"
  type       = "RSA"
  size       = 2048
  operations = ["sign", "verify"]
}

resource "ovh_okms_service_key" "key_ecdsa" {
  okms_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name       = "key_ecdsa"
  type       = "EC"
  curve      = "P-256"
  operations = ["sign", "verify"]
}
```

## Argument Reference

### Required

- `name` (String) Key name
- `okms_id` (String) Okms ID
- `operations` (List of String) The operations for which the key is intended to be used
- `type` (String) Type of key to create

### Optional

- `context` (String) Context of the key
- `curve` (String) Curve type, for Elliptic Curve (EC) keys (Either P-256, P-384 or P-521)
- `size` (Number) Size of the key to be created, for symmetric and RSA keys (One of 128, 192 or 256 for symmetric keys, or one of 2048, 3072 or 4096 for RSA keys)

### Read-Only

- `created_at` (String) Creation time of the key
- `deactivation_reason` (String) Key deactivation reason
- `id` (String) Key ID
- `state` (String) State of the key
