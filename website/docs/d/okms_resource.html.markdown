---
subcategory: "Key Management Service (KMS)"
---

# ovh_okms_resource (Data Source)

Use this data source to retrieve information about a KMS associated with this account

## Example Usage

```hcl
data "ovh_okms_resource" "kms" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

```

## Argument Reference

- `id` (String) Should be set to the ID of your KMS

## Attributes Reference

- `iam` (Attributes) IAM resource metadata (see [below for nested schema](#nestedatt--iam))
- `kmip_endpoint` (String) KMS kmip API endpoint
- `public_ca` (String) KMS public CA (Certificate Authority)
- `region` (String) Region
- `rest_endpoint` (String) KMS rest API endpoint
- `swagger_endpoint` (String) KMS rest API swagger UI

<a id="nestedatt--iam"></a>
### Nested Schema for `iam`

Read-Only:

- `display_name` (String) Resource display name
- `id` (String) Unique identifier of the resource
- `tags` (Map of String) Resource tags. Tags that were internally computed are prefixed with ovh:
- `urn` (String) Unique resource name used in policies
