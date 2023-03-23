---
layout: "ovh"
page_title: "OVH: ovh_domain_zone"
sidebar_current: "docs-ovh-resource-domain-zone-x"
description: |-
  Creates a domain zone.
---

# ovh_domain_zone

Creates a domain zone.

## Important

-> __NOTE__ To order a product with terraform, your account needs to have a default payment method defined. This can be done in the [console](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) routes

## Example Usage

```hcl
data "ovh_order_cart" "mycart" {
  ovh_subsidiary = "fr"
}

data "ovh_order_cart_product_plan" "zone" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "dns"
  plan_code      = "zone"
}

resource "ovh_domain_zone" "zone" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary

  plan {
    duration     = data.ovh_order_cart_product_plan.zone.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.zone.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.zone.selected_price.0.pricing_mode

    configuration {
      label = "zone"
      value = "myzone.mydomain.com"
    }

    configuration {
      label = "template"
      value = "minimized"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

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

* `dnssec_supported` - Is DNSSEC supported by this zone
* `has_dns_anycast` - hasDnsAnycast flag of the DNS zone
* `last_update` - Last update date of the DNS zone
* `name` - Zone name
* `name_servers` - Name servers that host the DNS zone

* `order` - Details about an Order
  * `date` - date
  * `order_id` - order id
  * `expiration_date` - expiration date
  * `details` - Information about a Bill entry
    * `description` - description
    * `order_detail_id` - order detail id
    * `domain` - expiration date
    * `quantity` - quantity
