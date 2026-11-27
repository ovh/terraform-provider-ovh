---
page_title: "Choosing the plan, duration and pricing mode of an ordered product"
---

# Choosing the plan, duration and pricing mode of an ordered product

Resources that place an order, such as `ovh_dedicated_server`, `ovh_cloud_project` or `ovh_vrack`, take a `plan` block and, for some products, `plan_option` blocks. Each block describes one offer to order:

* `plan_code`: identifier of the offer, e.g. `vps-le-2-2-40`.
* `duration`: duration of the subscription, in [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601#Durations) format. It is the period billed at each renewal, e.g. `P1M` (one month), `P3M` (three months), `P6M` (six months), `P1Y` (one year) or `P2Y` (two years).
* `pricing_mode`: pricing mode of the offer, e.g. `default`.

`duration` and `pricing_mode` must be taken from the same price of the offer. The available values depend on the product, the offer and your OVHcloud subsidiary.

## List the available values

The `prices` of the [`ovh_order_cart_product_plan`](../data-sources/order_cart_product_plan) data source list the prices of an offer. The prices that can be used for a subscription are the ones with the `renew` capacity. The other prices, such as `installation` fees, have a `P0D` duration.

```terraform
data "ovh_me" "my_account" {}

data "ovh_order_cart" "my_cart" {
  ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "vps" {
  cart_id        = data.ovh_order_cart.my_cart.id
  price_capacity = "renew"
  product        = "vps"
  plan_code      = "vps-le-2-2-40"
}

output "vps_prices" {
  value = [
    for price in data.ovh_order_cart_product_plan.vps.prices : {
      duration     = price.duration
      pricing_mode = price.pricing_mode
    }
    if contains(price.capacities, "renew")
  ]
}
```

For this VPS offer, the output is:

```
vps_prices = [
  {
    "duration" = "P1M"
    "pricing_mode" = "default"
  },
  {
    "duration" = "P1M"
    "pricing_mode" = "degressivity12"
  },
  {
    "duration" = "P1M"
    "pricing_mode" = "degressivity24"
  },
  {
    "duration" = "P1Y"
    "pricing_mode" = "upfront12"
  },
  {
    "duration" = "P2Y"
    "pricing_mode" = "upfront24"
  },
]
```

Here, `default` is billed every month. `upfront12` and `upfront24` are paid upfront for one or two years. `degressivity12` and `degressivity24` are billed every month, with a commitment of 12 or 24 months.

The same values can be found with the API: [`GET /order/cart/{cartId}/{product}`](https://eu.api.ovh.com/console/?section=%2Forder&branch=v1#get-/order/cart/-cartId-/vps) lists the offers of a product with their `prices`.

~> **NOTE** The `selected_price` attribute of `ovh_order_cart_product_plan` holds a single price with the requested `price_capacity`. When an offer has several `renew` prices, like most VPS and dedicated servers, it is not necessarily the monthly one: choose the price from `prices` instead.

## Use a price in the order

The following example orders a dedicated server billed every month. Its options use the same price as its plan:

```terraform
data "ovh_me" "my_account" {}

data "ovh_order_cart" "my_cart" {
  ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "server" {
  cart_id        = data.ovh_order_cart.my_cart.id
  price_capacity = "renew"
  product        = "baremetalServers"
  plan_code      = "24rise01"
}

locals {
  server_price = one([
    for price in data.ovh_order_cart_product_plan.server.prices : price
    if contains(price.capacities, "renew") && price.duration == "P1M" && price.pricing_mode == "default"
  ])
}

resource "ovh_dedicated_server" "server" {
  ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary
  os             = "debian12_64"

  plan = [
    {
      duration     = local.server_price.duration
      plan_code    = data.ovh_order_cart_product_plan.server.plan_code
      pricing_mode = local.server_price.pricing_mode

      configuration = [
        {
          label = "dedicated_datacenter"
          value = "bhs"
        },
        {
          label = "dedicated_os"
          value = "none_64.en"
        },
        {
          label = "region"
          value = "canada"
        }
      ]
    }
  ]

  plan_option = [
    for option in ["ram-32g-rise13", "bandwidth-500-included-rise", "softraid-2x512nvme-rise", "vrack-bandwidth-100-rise-included"] : {
      duration     = local.server_price.duration
      plan_code    = option
      pricing_mode = local.server_price.pricing_mode
      quantity     = 1
    }
  ]
}
```

Changing the `plan` or `plan_option` blocks after the order replaces the resource, which places a new order. To learn which rights the provider needs to place orders, see [Consumer key access rules for resources that place orders](order_consumer_key_access_rules).
