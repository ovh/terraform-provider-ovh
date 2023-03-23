---
layout: "ovh"
page_title: "OVH: ovh_ip_service"
sidebar_current: "docs-ovh-resource-ip-service-x"
description: |-
  Orders an ip service.
---

# ovh_ip_service

Orders an ip service.


## Important

This resource orders an OVHcloud product for a long period of time and may generate heavy costs !
Use with caution.

-> __NOTE__ To order a product with terraform, your account needs to have a default payment method defined. This can be done in the [console](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) routes


## Example Usage

```hcl
data "ovh_order_cart" "mycart" {
  ovh_subsidiary = "fr"
  description    = "order ip block"
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
 }
}
```

## Argument Reference

The following arguments are supported:

* `description` - Custom description on your ip.
* `ovh_subsidiary` - (Required) OVHcloud Subsidiary
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
