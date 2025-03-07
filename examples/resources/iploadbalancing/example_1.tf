data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
 ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "iplb" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "ipLoadbalancing"
  plan_code      = "iplb-lb1"
}

data "ovh_order_cart_product_options_plan" "bhs" {
  cart_id           = data.ovh_order_cart_product_plan.iplb.cart_id
  price_capacity    = data.ovh_order_cart_product_plan.iplb.price_capacity
  product           = data.ovh_order_cart_product_plan.iplb.product
  plan_code         = data.ovh_order_cart_product_plan.iplb.plan_code
  options_plan_code = "iplb-zone-lb1-rbx"
}

resource "ovh_iploadbalancing" "iplb-lb1" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  display_name   = "my ip loadbalancing"

  plan {
    duration     = data.ovh_order_cart_product_plan.iplb.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.iplb.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.iplb.selected_price.0.pricing_mode
  }

  plan_option {
    duration     = data.ovh_order_cart_product_options_plan.bhs.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_options_plan.bhs.plan_code
    pricing_mode = data.ovh_order_cart_product_options_plan.bhs.selected_price.0.pricing_mode
  }
}
