---
subcategory : "vRack"
---

# ovh_vrack

Orders a vrack.

## Important

-> **NOTE** To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

~> **WARNING** `BANK_ACCOUNT` is not supported anymore, please update your default payment_method to `SEPA_DIRECT_DEBIT`

## Example Usage

```terraform
data "ovh_me" "my_account" {}

data "ovh_order_cart" "my_cart" {
  ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "vrack" {
  cart_id        = data.ovh_order_cart.my_cart.id
  price_capacity = "renew"
  product        = "vrack"
  plan_code      = "vrack"
}

resource "ovh_vrack" "vrack" {
  ovh_subsidiary = data.ovh_order_cart.my_cart.ovh_subsidiary
  name           = "my-vrack"
  description    = "my vrack"

  plan {
    duration     = data.ovh_order_cart_product_plan.vrack.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.vrack.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.vrack.selected_price.0.pricing_mode
  }
}
```

## Argument Reference

The following arguments are supported:
* `description` - yourvrackdescription
* `name` - yourvrackname
* `ovh_subsidiary` - (Required) OVHcloud Subsidiary. Country of OVHcloud legal entity you'll be billed by. List of supported subsidiaries available on API at [/1.0/me.json under `models.nichandle.OvhSubsidiaryEnum`](https://eu.api.ovh.com/1.0/me.json)
* `plan` - (Required) Product Plan to order
  * `duration` - (Required) duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing model identifier
  * `catalog_name` - Catalog name
  * `configuration` - (Optional) Representation of a configuration item for personalizing product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in API.OVH.COM
* `plan_option` - (Optional) Product Plan to order
  * `duration` - (Required) duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing model identifier
  * `catalog_name` - Catalog name
  * `configuration` - (Optional) Representation of a configuration item for personalizing product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in API.OVH.COM

## Attributes Reference

Id is set to the order Id. In addition, the following attributes are exported:

* `urn` - The URN of the vrack, used with IAM permissions
* `order` - Details about an Order
  * `date` - date
  * `order_id` - order id
  * `expiration_date` - expiration date
  * `details` - Information about a Bill entry
    * `description` - description
    * `order_detail_id` - order detail id
    * `domain` - expiration date
    * `quantity` - quantity
* `service_name` - The internal name of your vrack

## Timeouts

```terraform
resource "ovh_vrack" "vrack" {
  # ...

  timeouts {
    create = "1h"
  }
}
```

* `create` - (Default 30m)

## Import

A vRack can be imported using the `service_name`. Using the following configuration:

```terraform
import {
  to = ovh_vrack.vrack
  id = "<service name>"
}
```

You can then run:

```bash
$ terraform plan -generate-config-out=vrack.tf
$ terraform apply
```

The file `vrack.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above. See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
