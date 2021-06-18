---
layout: "ovh"
page_title: "OVH: vrack_ip"
sidebar_current: "docs-ovh-resource-vrack-ip"
description: |-
  Attach an Ip block to a VRack.
---

# ovh_vrack_ip

Attach a Ip block to a VRack.

## Example Usage

```hcl
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "my cart"
}

data "ovh_order_cart_product_plan" "vrack" {
 cart_id        = data.ovh_order_cart.mycart.id
 price_capacity = "renew"
 product        = "vrack"
 plan_code      = "vrack"
}

resource "ovh_vrack" "vrack" {
 description    = data.ovh_order_cart.mycart.description
 name           = data.ovh_order_cart.mycart.description
 ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
 payment_mean   = "fidelity"

 plan {
   duration     = data.ovh_order_cart_product_plan.vrack.selected_price.0.duration
   plan_code    = data.ovh_order_cart_product_plan.vrack.plan_code
   pricing_mode = data.ovh_order_cart_product_plan.vrack.selected_price.0.pricing_mode
 }
}

data "ovh_order_cart_product_plan" "ipblock" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "ip"
  plan_code      = "ip-v4-s30-ripe"
}

resource "ovh_ip_service" "ipblock" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  payment_mean   = "ovh-account"
  description    = data.ovh_order_cart.mycart.description

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

resource "ovh_vrack_ip" "vrackblock" {
  service_name = ovh_vrack.vrack.service_name
  block        = ovh_ip_service.ipblock.ip
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your vrack
* `block` - (Required) Your IP block.
    
## Attributes Reference

The following attributes are exported:

* `gateway` - Your gateway
* `ip` - Your IP block
* `zone` - Where you want your block announced on the network
