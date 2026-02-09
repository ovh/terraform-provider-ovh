---
subcategory: "Enterprise File Storage"
---

# ovh_storage_efs (Resource)

Order and manage an Enterprise File Storage service.

## Important

-> **NOTE** To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

~> **WARNING** `BANK_ACCOUNT` is not supported anymore, please update your default payment_method to `SEPA_DIRECT_DEBIT`

## Example Usage

```terraform
data "ovh_me" "my_account" {}

resource "ovh_storage_efs" "efs" {
  name = "MyEFS"

  ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary

  plan = [
    {
      plan_code    = "enterprise-file-storage-premium-1tb"
      duration     = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region"
          value = "eu-west-gra"
        },
        {
          label = "network"
          value = "vrack"
        }
      ]
    }
  ]
}
```

## Schema

### Optional

- `name` (String) Custom service display name
- `ovh_subsidiary` (String) OVH subsidiaries
- `plan` (Attributes List) (see [below for nested schema](#nestedatt--plan))
- `plan_option` (Attributes List) (see [below for nested schema](#nestedatt--plan_option))

### Read-Only

- `created_at` (String) Service creation date
- `iam` (Attributes) IAM resource metadata (see [below for nested schema](#nestedatt--iam))
- `id` (String) Service ID
- `order` (Attributes) Details about an Order (see [below for nested schema](#nestedatt--order))
- `performance_level` (String) Service performance level
- `product` (String) Product name
- `quota` (Number) Service quota
- `region` (String) Service region
- `service_name` (String) Service name (same as `id`)
- `status` (String) Service status

<a id="nestedatt--plan"></a>
### Nested Schema for `plan`

Required:

- `duration` (String) Duration selected for the purchase of the product
- `plan_code` (String) Identifier of the option offer
- `pricing_mode` (String) Pricing mode selected for the purchase of the product

Optional:

- `configuration` (Attributes List) (see [below for nested schema](#nestedatt--plan--configuration))
- `item_id` (Number) Cart item to be linked
- `quantity` (Number) Quantity of product desired

<a id="nestedatt--plan--configuration"></a>
### Nested Schema for `plan.configuration`

Required:

- `label` (String) Label for your configuration item
- `value` (String) Value or resource URL on API.OVH.COM of your configuration item



<a id="nestedatt--plan_option"></a>
### Nested Schema for `plan_option`

Required:

- `duration` (String) Duration selected for the purchase of the product
- `plan_code` (String) Identifier of the option offer
- `pricing_mode` (String) Pricing mode selected for the purchase of the product
- `quantity` (Number) Quantity of product desired

Optional:

- `configuration` (Attributes List) (see [below for nested schema](#nestedatt--plan_option--configuration))

<a id="nestedatt--plan_option--configuration"></a>
### Nested Schema for `plan_option.configuration`

Required:

- `label` (String) Label for your configuration item
- `value` (String) Value or resource URL on API.OVH.COM of your configuration item



<a id="nestedatt--iam"></a>
### Nested Schema for `iam`

Read-Only:

- `display_name` (String) Resource display name
- `id` (String) Unique identifier of the resource
- `tags` (Map of String) Resource tags. Tags that were internally computed are prefixed with ovh:
- `urn` (String) Unique resource name used in policies


<a id="nestedatt--order"></a>
### Nested Schema for `order`

Read-Only:

- `date` (String)
- `details` (Attributes List) (see [below for nested schema](#nestedatt--order--details))
- `expiration_date` (String)
- `order_id` (Number)

<a id="nestedatt--order--details"></a>
### Nested Schema for `order.details`

Read-Only:

- `description` (String)
- `detail_type` (String) Product type of item in order
- `domain` (String)
- `order_detail_id` (Number)
- `quantity` (String)

## Timeouts

```terraform
resource "ovh_storage_efs" "efs" {
  # ...

  timeouts {
    create = "1h"
  }
}
```

* `create` - (Default 30m)

## Import

An Enterprise File Storage service can be imported using its `id`. Using the following configuration:

```terraform
import {
  to = ovh_storage_efs.efs
  id = "xxx-xxx-xxx-xxx-xxx"
}
```

You can then run:

```bash
$ terraform plan -generate-config-out=efs.tf
$ terraform apply
```

The file `efs.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above. See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
