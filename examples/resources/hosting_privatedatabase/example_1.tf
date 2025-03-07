data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "database" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "privateSQL"
  plan_code      = "private-sql-512-instance"
}

resource "ovh_hosting_privatedatabase" "database" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  display_name   = "Postgresql-12"

  plan {
    duration     = data.ovh_order_cart_product_plan.database.prices[3].duration
    plan_code    = data.ovh_order_cart_product_plan.database.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.database.selected_price[0].pricing_mode

    configuration {
      label = "dc"
      value = "gra3"
    }

    configuration {
      label = "engine"
      value = "postgresql_12"
    }
  }
}

output "privatedatabase_service_name" {
  value = ovh_hosting_privatedatabase.database.service_name
}
