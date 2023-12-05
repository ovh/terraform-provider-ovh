---
subcategory : "Domain names"
---

# ovh_domain_zone

Creates a domain zone.

## Important

-> __NOTE__ To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

## Example Usage

```hcl
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "zone" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "dns"
  plan_code      = "zone"
}

resource "ovh_domain_zone" "zone" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary

  plan {
    duration     = data.ovh_order_cart_product_plan.zone.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.zone.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.zone.selected_price.0.pricing_mode

    configuration {
      label = "zone"
      value = "myzone.mydomain.com"
    }

    configuration {
      label = "template"
      value = "minimized"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `ovh_subsidiary` - (Required) OVHcloud Subsidiary. Country of OVHcloud legal entity you'll be billed by. List of supported subsidiaries available on API at [/1.0/me.json under `models.nichandle.OvhSubsidiaryEnum`](https://eu.api.ovh.com/1.0/me.json)
* `plan` - (Required) Product Plan to order
  * `duration` - (Required) duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing model identifier
  * `catalog_name` - Catalog name
  * `configuration` - (Required) Representation of a configuration item for personalizing product. 2 configurations are required : one for `zone` and one for `template`
    * `label` - (Required) Identifier of the resource : `zone` or `template`
    * `value` - (Required) For `zone`, the value is the zone name `myzone.mydomain.com`. For `template`, the value can be `basic`, `minimized` or  `redirect` which is the same as `minimized` with additional entries for a redirect configuration.
* `plan_option` - (Optional) Product Plan to order
  * `duration` - (Required) duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing model identifier
  * `catalog_name` - Catalog name
  * `configuration` - (Optional) Representation of a configuration item for personalizing product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in API.OVH.COM


## Attributes Reference

Id is set to the order Id. In addition, the following attributes are exported:

* `dnssec_supported` - Is DNSSEC supported by this zone
* `has_dns_anycast` - hasDnsAnycast flag of the DNS zone
* `last_update` - Last update date of the DNS zone
* `name` - Zone name
* `name_servers` - Name servers that host the DNS zone

* `order` - Details about an Order
  * `date` - date
  * `order_id` - order id
  * `expiration_date` - expiration date
  * `details` - Information about a Bill entry
    * `description` - description
    * `order_detail_id` - order detail id
    * `domain` - expiration date
    * `quantity` - quantity

## Import
Zone can be imported using the `order_id` that can be retrieved in the [order page](https://www.ovh.com/manager/#/dedicated/billing/orders/orders) at the creation time of the zone. 
```bash
$ terraform import ovh_domain_zone.zone order_id
```