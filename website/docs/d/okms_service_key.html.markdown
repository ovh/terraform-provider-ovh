---
subcategory: "KMS"
---

# ovh_okms_service_key (Data Source)

Use this data source to retrieve information about a KMS service key.

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
- `curve` (String) Curve type for Elliptic Curve (EC) keys
- `iam` (Attributes) IAM resource metadata (see [below for nested schema](#nestedatt--iam))
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
