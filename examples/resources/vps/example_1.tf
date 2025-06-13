data "ovh_me" "my_account" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "vps" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "vps"
  plan_code      = "vps-le-2-2-40"
}

resource "ovh_vps" "my_vps" {
  display_name = "dev_vps"

  image_id = "45b2f222-ab10-44ed-863f-720942762b6f"

  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  plan = [
    {
      duration     = "P1M"
      plan_code    = data.ovh_order_cart_product_plan.vps.plan_code
      pricing_mode = "default"

      configuration = [
        {
          label = "vps_datacenter"
          value = "WAW"
        },
        {
          label = "vps_os"
          value = "Debian 10"
        }
      ]
    }
  ]

  public_ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDSD76EaLUzJjf70W8W2uU9FzEyl68di67Bd20qtYfBLJpFJuX/RJC9StI1y1RnXXqC1Lf/Yo+yJzvNx0iqLxCX1G7g0XYex74HkgC6a2QeNhp9M56ANZtA3TKKAbkZ1xobfhOPWpq3lEFp7dgJctcILBPL3l6OjKf6NIxHo5yF67Vy4D0nWl5utumNdWhhlX7MtVQooszLyIwPlNO+DzD3ZnJFCt2Z1jdRkhm/Oobtx17CZ+5SN23tgHXS6pLOgM6w30M11zkI510z95IAIHhRT7MbiXICkvG/0qHuSftz1j/CcHFbttNB27dH86vByumfSEgRKaoRkCqrn64IWrSsFr3Smsf7gZWLBlYLliGPyn8Tsr9bT5pRul6yTvVbfZ31RREBr1I0Lp4q++d+fIpa3LtMGRaMb9huJYy8cwW/Vfzbxsqfz9xzjIOFNcYl7J9l4cvz3hgSlai2Jgngw5ShNVlxcIKUdiynZWm09nQudlYNHgor9ID+JACzCfPkUZ8"
}

output "vps_display_name" {
  value = ovh_vps.my_vps.display_name
}
