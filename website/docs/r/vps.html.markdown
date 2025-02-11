---
subcategory : "VPS"
---

# ovh_vps

Creates an OVHcloud Virtual Private Server (VPS).

## Important

-> __NOTE__ To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

~> __WARNING__ `BANK_ACCOUNT` is not supported anymore, please update your default payment method to `SEPA_DIRECT_DEBIT`

## Example Usage

```hcl
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
}

output "vps_display_name" {
  value = ovh_vps.my_vps.display_name
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - Custom display name
* `netboot_mode` - VPS netboot mode (local┃rescue)
* `ovh_subsidiary` - (Required) OVHcloud Subsidiary. Country of OVHcloud legal entity you'll be billed by. List of supported subsidiaries available on API at [/1.0/me.json](https://eu.api.ovh.com/console-preview/?section=%2Fme&branch=v1#get-/me)
* `plan` - (Required) Product Plan to order
  * `duration` - (Required) duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing model identifier
  * `configuration` - (Optional) Representation of a configuration item for personalizing product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in api.ovh.com
* `plan_option` - (Optional) Product Plan to order
  * `duration` - (Required) duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing model identifier
  * `configuration` - (Optional) Representation of a configuration item for personalizing product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in api.ovh.com

## Attributes Reference

The following attributes are exported:

* `iam` - IAM resource information
  * `urn` - URN of the private database, used when writing IAM policies
  * `display_name` - Resource display name
  * `id` - Unique identifier of the resource in the IAM
  * `tags` - Resource tags. Tags that were internally computed are prefixed with `ovh:`
* `cluster` - VPS cluster
* `display_name` - Custom display name
* `keymap` - KVM keyboard layout on VPS Cloud
* `memory_limit` - RAM of this VPS
* `model` - Structure describing characteristics of a VPS model
  * `available_options` - All options the VPS can have (additionalDisk┃automatedBackup┃cpanel┃ftpbackup┃plesk┃snapshot┃veeam┃windows)
  * `datacenter` - Datacenters where this model is available
  * `disk` - Disk capacity of this VPS
  * `maximum_additionnal_ip` - Maximum number of additional IPs
  * `memory` - RAM of the VPS
  * `name` - Plan code of the VPS
  * `offer` - Description of this VPS offer
  * `vcore` - Number of vcores
  * `version` - All versions that VPS can have (2013v1┃2014v1┃2015v1┃2017v1┃2017v2┃2017v3┃2018v1┃2018v2┃2019v1)
* `monitoring_ip_blocks` - IP blocks for OVH monitoring servers
* `name` - Name of the VPS
* `netboot_mode` - VPS netboot mode (local┃rescue)
* `offer_type` - All offers a VPS can have (beta-classic┃classic┃cloud┃cloudram┃game-classic┃lowlat┃ssd)
* `sla_monitoring`
* `state` - State of the VPS (backuping┃installing┃maintenance┃rebooting┃rescued┃running┃stopped┃stopping┃upgrading)
* `vcore` - Number of vcores
* `zone` - OpenStask region where the VPS is located

## Import

OVHcloud VPS database can be imported using the `service_name`, E.g.,

```hcl
import {
  to = ovh_vps.myvps
  id = "<your existing VPS service_name>"
}
```

You can then run:

```sh
terraform plan -generate-config-out=./vps.tf
```

The file `vps.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above.
See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.