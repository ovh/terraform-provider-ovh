---
subcategory : "Additional IP"
---

# ovh_ip_service

Orders an ip service.

## Important

This resource orders an OVHcloud product for a long period of time and may generate heavy costs ! Use with caution.

-> **NOTE** To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

~> **WARNING** `BANK_ACCOUNT` is not supported anymore, please update your default payment method to `SEPA_DIRECT_DEBIT`

## Example Usage

```terraform
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = "fr"
}

data "ovh_order_cart_product_plan" "ipblock" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "ip"
  plan_code      = "ip-v4-s30-ripe"
}

resource "ovh_ip_service" "ipblock" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  description    = "my ip block"

  plan {
    duration     = data.ovh_order_cart_product_plan.ipblock.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.ipblock.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.ipblock.selected_price.0.pricing_mode

    configuration {
      label = "country"
      value = "FR"
    }

    configuration {
      label = "region"
      value = "europe"
    }

    configuration {
      label = "destination"
      value = "parking"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - Custom description on your ip.
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
  * `configuration` - (Optional) Representation of a configuration item for personalizing product. The list of available configurations can be retrieved using call [GET /order/cart/{cartId}/item/{itemId}/requiredConfiguration](https://eu.api.ovh.com/console/?section=%2Forder&branch=v1#get-/order/cart/-cartId-/item/-itemId-/requiredConfiguration)
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in API.OVH.COM

## Attributes Reference

Id is set to the order Id. In addition, the following attributes are exported:

* `can_be_terminated` - can be terminated
* `country` - country
* `ip` - ip block
* `order` - Details about an Order
  * `date` - date
  * `order_id` - order id
  * `expiration_date` - expiration date
  * `details` - Information about a Bill entry
    * `description` - description
    * `order_detail_id` - order detail id
    * `domain` - expiration date
    * `quantity` - quantity
* `organisation_id` - IP block organisation Id
* `routed_to` - Routage information
  * `service_name` - Service where ip is routed to
* `service_name`: service name
* `type` - Possible values for ip type

## Timeouts

```terraform
resource "ovh_ip_service" "ipblock" {
  # ...

  timeouts {
    create = "1h"
  }
}
```

* `create` - (Default 30m)

## Import

The resource can be imported using its `service_name`, E.g.,

```terraform
import {
  to = ovh_ip_service.ipblock
  id = "ip-xx.xx.xx.xx"
}
```

```bash
$ terraform plan -generate-config-out=ipblock.tf
$ terraform apply
```

The file `ipblock.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above. See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
