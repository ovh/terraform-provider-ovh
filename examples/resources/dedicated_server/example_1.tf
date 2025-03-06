data "ovh_me" "account" {}

resource "ovh_dedicated_server" "server" {
  ovh_subsidiary = data.ovh_me.account.ovh_subsidiary
  display_name = "My server display name"
  os = "debian12_64"
  plan = [
    {
      plan_code = "24rise01"
      duration = "P1M"
      pricing_mode = "default"

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
    {
      duration = "P1M"
      plan_code = "ram-32g-rise13"
      pricing_mode = "default"
      quantity = 1
    },
    {
      duration = "P1M"
      plan_code = "bandwidth-500-included-rise"
      pricing_mode = "default"
      quantity = 1
    },
    {
      duration = "P1M"
      plan_code = "softraid-2x512nvme-rise"
      pricing_mode = "default"
      quantity = 1
    },
    {
      duration = "P1M"
      plan_code = "vrack-bandwidth-100-rise-included"
      pricing_mode = "default"
      quantity = 1
    }
  ]
}
