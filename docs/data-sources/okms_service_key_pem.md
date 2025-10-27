---
subcategory: "Key Management Service (KMS)"
---

# ovh_okms_service_key_pem (Data Source)

Use this data source to retrieve information about a KMS service key, in the PEM format.

## Example Usage

```terraform
data "ovh_okms_service_key_pem" "key_info" {
  okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

### Required

- `id` (String) ID of the service key
- `okms_id` (String) ID of the KMS

### Read-Only

- `created_at` (String) Creation time of the key
- `curve` (String) Curve type for Elliptic Curve (EC) keys
- `iam` (Attributes) IAM resource metadata (see [below for nested schema](#nestedatt--iam))
- `keys_pem` (Attributes List) The keys in PEM format (see [below for nested schema](#nestedatt--keys_pem))
- `name` (String) Key name
- `operations` (List of String) The operations for which the key is intended to be used
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

<a id="nestedatt--keys_pem"></a>

### Nested Schema for `keys_pem`

Read-Only:

- `pem` (String) The key in base64 encoded PEM format
