data "ovh_me" "my_account" {}

data "ovh_order_cart" "my_cart" {
  ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary
}

data "ovh_order_cart_product_options_plan" "plan" {
  cart_id           = data.ovh_order_cart.my_cart.id
  price_capacity    = "renew"
  product           = "cloud"
  plan_code         = "project"
  options_plan_code = "vrack"
}
