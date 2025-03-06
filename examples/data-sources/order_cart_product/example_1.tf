data "ovh_me" "my_account" {}

data "ovh_order_cart" "my_cart" {
  ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary
}

data "ovh_order_cart_product" "plans" {
  cart_id = data.ovh_order_cart.my_cart.id
  product = "..."
}
