---
subcategory: "Key Management Service (KMS)"
---

# ovh_okms_service_key_jwk (Data Source)

Use this data source to retrieve information about a KMS service key, in the JWK format.

## Example Usage

```hcl
data "ovh_okms_service_key" "key_info" {
  okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `id` (String) ID of the service key
- `okms_id` (String) ID of the KMS

### Read-Only

- `created_at` (String) Creation time of the key
- `iam` (Attributes) IAM resource metadata (see [below for nested schema](#nestedatt--iam))
- `keys` (Attributes List) The key in JWK format (see [below for nested schema](#nestedatt--keys))
- `name` (String) Key name
- `size` (Number) Size of the key
- `state` (String) State of the key
- `type` (String) Key type

<a id="nestedatt--iam"></a>
### Nested Schema for `iam`

Read-Only:

- `display_name` (String) Resource display name
- `id` (String) Unique identifier of the resource
- `tags` (Map of String) Resource tags. Tags that were internally computed are prefixed with ovh:
- `urn` (String) Unique resource name used in policies


<a id="nestedatt--keys"></a>
### Nested Schema for `keys`

Read-Only:

- `alg` (String) The algorithm intended to be used with the key
- `crv` (String) The cryptographic curve used with the key
- `e` (String) The exponent value for the RSA public key
- `key_ops` (List of String) The operation for which the key is intended to be used
- `kid` (String) key ID parameter used to match a specific key
- `kty` (String) Key type parameter identifies the cryptographic algorithm family used with the key, such as RSA or EC
- `n` (String) The modulus value for the RSA public key
- `use` (String) The intended use of the public key
- `x` (String) The x coordinate for the Elliptic Curve point
- `y` (String) The y coordinate for the Elliptic Curve point
