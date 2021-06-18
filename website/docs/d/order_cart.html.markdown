---
layout: "ovh"
page_title: "OVH: order_cart"
sidebar_current: "docs-ovh-datasource-order-cart"
description: |-
  Create a temporary order cart to retrieve information order cart products.
---

# ovh_order_cart

Use this data source to create a temporary order cart to retrieve information order cart products.

## Example Usage

```hcl
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "..."
}
```

## Argument Reference


* `ovh_subsidiary` - (Required) Ovh Subsidiary
* `description` - Description of your cart
* `expire` - Expiration time (format: 2006-01-02T15:04:05+00:00)


## Attributes Reference

`id` is set to your cart ID
In addition, the following attributes are exported.

* `cart_id` - Cart identifier
* `read_only` - Indicates if the cart has already been validated
* `items` - Items of your cart
