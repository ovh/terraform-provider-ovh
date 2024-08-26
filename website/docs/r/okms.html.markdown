---
subcategory : "KMS"
---

# ovh_okms (Resource)

Creates an OVHcloud Key Management Service (okms).

## Important

-> __NOTE__ To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.


## Example Usage

```hcl
resource "ovh_okms" "newkms" {
  ovh_subsidiary = "FR"
  region         = "EU_WEST_RBX"
  display_name   = "terraformed KMS"
}

```

## Argument Reference

### Required

- `ovh_subsidiary` (String) OVH subsidiary
- `region` (String) KMS region

### Optional

- `display_name` (String) Set the name displayed in Manager for this KMS

## Attributes Reference

- `iam` (Attributes) IAM resource metadata (see [below for nested schema](#nestedatt--iam))
- `id` (String) OKMS ID
- `kmip_endpoint` (String) KMS kmip API endpoint
- `public_ca` (String) KMS public CA (Certificate Authority)
- `rest_endpoint` (String) KMS rest API endpoint
- `swagger_endpoint` (String) KMS rest API swagger UI

<a id="nestedatt--iam"></a>
### Nested Schema for `iam`

Read-Only:

- `display_name` (String) Resource display name
- `id` (String) Unique identifier of the resource
- `tags` (Map of String) Resource tags. Tags that were internally computed are prefixed with ovh:
- `urn` (String) Unique resource name used in policies
