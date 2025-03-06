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
  description    = data.ovh_order_cart.my_cart.description
  name           = data.ovh_order_cart.my_cart.description
  ovh_subsidiary = data.ovh_order_cart.my_cart.ovh_subsidiary

  plan {
    duration     = data.ovh_order_cart_product_plan.vrack.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.vrack.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.vrack.selected_price.0.pricing_mode
  }
}

data "ovh_order_cart_product_plan" "ipblock" {
  cart_id        = data.ovh_order_cart.my_cart.id
  price_capacity = "renew"
  product        = "ip"
  plan_code      = "ip-v4-s30-ripe"
}

resource "ovh_ip_service" "ipblock" {
  ovh_subsidiary = data.ovh_order_cart.my_cart.ovh_subsidiary
  description    = data.ovh_order_cart.my_cart.description

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

resource "ovh_vrack_ip" "vrack_block" {
  service_name = ovh_vrack.vrack.service_name
  block        = ovh_ip_service.ipblock.ip
}
