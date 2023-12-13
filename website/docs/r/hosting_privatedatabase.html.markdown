---
subcategory : "Web Cloud Private SQL"
---

# ovh_hosting_privatedatabase

Creates an OVHcloud managed private cloud database.

## Important

-> __NOTE__ To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

~> __WARNING__ `BANK_ACCOUNT` is not supported anymore, please update your default payment method to `SEPA_DIRECT_DEBIT`

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `description` - Custom description on your privatedatabase order.
* `ovh_subsidiary` - (Required) OVHcloud Subsidiary. Country of OVHcloud legal entity you'll be billed by. List of supported subsidiaries available on API at [/1.0/me.json under `models.nichandle.OvhSubsidiaryEnum`](https://eu.api.ovh.com/1.0/me.json)
* `plan` - (Required) Product Plan to order
  * `duration` - (Required) duration.
  * `plan_code` - (Required) Plan code.
  * `pricing_mode` - (Required) Pricing model identifier
  * `catalog_name` - Catalog name
  * `configuration` - (Optional) Representation of a configuration item for personalizing product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in API.OVH.COM

Plan order valid values can be found on OVHcloud [APIv6](https://api.ovh.com/console/#/hosting/privateDatabase/availableOrderCapacities~GET)

## Attributes Reference

The following attributes are exported:

* `urn` - URN of the private database, used when writing IAM policies
* `cpu` - Number of CPU on your private database
* `datacenter` - Datacenter where this private database is located
* `display_name` - Name displayed in customer panel for your private database
* `hostname` - Private database hostname
* `hostname_ftp` - Private database FTP hostname
* `id` - Private database id
* `infrastructure` - Infrastructure where service was stored
* `offer` - Type of the private database offer
* `order` - Details about your Order
  * `date` - date
  * `order_id` - order id
  * `expiration_date` - expiration date
  * `details` - Information about a Bill entry
    * `description` - description
    * `order_detail_id` - order detail id
    * `domain` - expiration date
    * `quantity` - quantity
* `quantity` - quantity
* `ovh_subsidiary` - OVHcloud Subsidiary
* `plan` - Product Plan
  * `catalog_name` - Catalog name
  * `configuration` - Representation of a configuration item for personalizing product
  * `duration` - Service duration
  * `plan_code` - Plan code
  * `pricing_mode` - Pricing model identifier
* `plan_option`: Product Plan to order
* `port`: Private database service port
* `port_ftp`: Private database FTP port
* `quota_size`: Space allowed (in MB) on your private database
* `quota_used`: Sapce used (in MB) on your private database
* `ram`: Amount of ram (in MB) on your private database
* `server`: Private database server name
* `service_name`: Service name
* `state`: Private database state
* `type`: Private database type
* `version`: Private database available versions
* `version_label`: Private database version label
* `version_number`: Private database version number

## Import

OVHcloud Webhosting database can be imported using the `service_name`, E.g.,

```
$ terraform import ovh_hosting_privatedatabase.database service_name
```
