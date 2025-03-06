data "ovh_order_cart_product_options_plan" "hds" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "cloud"
  plan_code      = "project.2018"
  options_plan_code  = "certification.hds.2018"
  # plan_code    = "project" # when running in the US
  # options_plan_code  = "certification.hds" # when running in the US
}

resource "ovh_cloud_project" "my_cloud_project" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  description    = "my cloud project"

  plan {
    duration     = data.ovh_order_cart_product_plan.cloud.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.cloud.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.cloud.selected_price.0.pricing_mode
  }

  plan_option {
    duration     = data.ovh_order_cart_product_options_plan.hds.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_options_plan.hds.options_plan_code
    pricing_mode = data.ovh_order_cart_product_options_plan.hds.selected_price.0.pricing_mode
  }
}
