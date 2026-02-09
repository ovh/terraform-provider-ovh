data "ovh_me" "my_account" {}

resource "ovh_storage_efs" "efs" {
  name = "MyEFS"

  ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary

  plan = [
    {
      plan_code    = "enterprise-file-storage-premium-1tb"
      duration     = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "region"
          value = "eu-west-gra"
        },
        {
          label = "network"
          value = "vrack"
        }
      ]
    }
  ]
}
