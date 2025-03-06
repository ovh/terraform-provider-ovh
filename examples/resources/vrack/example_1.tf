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
