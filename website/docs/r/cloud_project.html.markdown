---
layout: "ovh"
page_title: "OVH: ovh_cloud_project"
sidebar_current: "docs-ovh-resource-cloud-project-x"
description: |-
  Orders a public cloud project.
---

# ovh_cloud_project

Orders a public cloud project.

## Important

This resource is in beta state. Use with caution.

## Example Usage

```hcl
data "ovh_order_cart" "mycart" {
  ovh_subsidiary = "fr"
  description    = "my cloud order cart"
}

data "ovh_order_cart_product_plan" "cloud" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "cloud"
  plan_code      = "project.2018"
}

resource "ovh_cloud_project" "my_cloud_project" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  description    = "my cloud project"
  payment_mean   = "fidelity"

  plan {
    duration     = data.ovh_order_cart_product_plan.cloud.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.cloud.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.cloud.selected_price.0.pricing_mode
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - A description associated with the user.
* `ovh_subsidiary` - (Required) OVHcloud Subsidiary
* `payment_mean` - (Required) OVHcloud payment mode (One of "default-payment-mean", "fidelity", "ovh-account")
* `plan` - (Required) Product Plan to order
  * `duration` - (Required) duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing model identifier
  * `catalog_name` - Catalog name
  * `configuration` - (Optional) Representation of a configuration item for personalizing product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in API.OVH.COM
* `plan_option` - (Optional) Product Plan to order
  * `duration` - (Required) duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing model identifier
  * `catalog_name` - Catalog name
  * `configuration` - (Optional) Representation of a configuration item for personalizing product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in API.OVH.COM


## Attributes Reference

`id` is set to the order Id. In addition, the following attributes are exported:

* `access` - project access right for the identity that trigger the terraform script.
* `description` - Project description
* `order` - Details about the order that was used to create the public cloud project
  * `date` - date
  * `order_id` - order id, the same as the `id`
  * `expiration_date` - expiration date
  * `details` - Information about a Bill entry
    * `description` - description
    * `order_detail_id` - order detail id
    * `domain` - expiration date
    * `quantity` - quantity
* `project_name` - openstack project name
* `project_id` - openstack project id
* `status` - project status

## Import
Cloud project can be imported using the `order_id` that can be retrieved in the [order page](https://www.ovh.com/manager/#/dedicated/billing/orders/orders) at the creation time of the Public Cloud project. 
```bash
$ terraform import ovh_cloud_project.my_cloud_project order_id
```
