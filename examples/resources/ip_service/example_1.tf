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
