---
subcategory: "Cloud Project"
---

# ovh_cloud_project

Orders a public cloud project.

## Important

-> **NOTE** To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.

~> **WARNING** `BANK_ACCOUNT` is not supported anymore, please update your default payment method to `SEPA_DIRECT_DEBIT`

## Example Usage

```hcl
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
  ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "cloud" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "cloud"
  plan_code      = "project.2018"
  # plan_code    = "project" # when running in the US

}

resource "ovh_cloud_project" "my_cloud_project" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  description    = "my cloud project"

  plan {
    duration     = data.ovh_order_cart_product_plan.cloud.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.cloud.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.cloud.selected_price.0.pricing_mode
  }
}
```

-> **WARNING** Currently, the OVHcloud Terraform provider does not support deletion of a public cloud project in the US. Removal is possible by manually deleting the project and then manually removing the public cloud project from terraform state.

## HDS Certification

If you want to add the certification HDS option at project creation (you must have a business support level on your account), you can add hds datasource and the plan_option configuration on the `ovh_cloud_project`.

```hcl

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

```

## Argument Reference

The following arguments are supported:

- `urn` - The URN of the cloud project
- `description` - A description associated with the user.
- `ovh_subsidiary` - (Required) OVHcloud Subsidiary. Country of OVHcloud legal entity you'll be billed by. List of supported subsidiaries available on API at [/1.0/me.json under `models.nichandle.OvhSubsidiaryEnum`](https://eu.api.ovh.com/1.0/me.json)
- `plan` - (Required) Product Plan to order
  - `duration` - (Required) duration
  - `plan_code` - (Required) Plan code. This value must be adapted depending on your `OVH_ENDPOINT` value. It's `project.2018` for `ovh-{eu,ca}` and `project` when using `ovh-us`.
  - `pricing_mode` - (Required) Pricing model identifier
  - `catalog_name` - Catalog name
  - `configuration` - (Optional) Representation of a configuration item for personalizing product
    - `label` - (Required) Identifier of the resource
    - `value` - (Required) Path to the resource in API.OVH.COM
- `plan_option` - (Optional) Product Plan to order
  - `duration` - (Required) duration
  - `plan_code` - (Required) Plan code
  - `pricing_mode` - (Required) Pricing model identifier
  - `catalog_name` - Catalog name
  - `configuration` - (Optional) Representation of a configuration item for personalizing product
    - `label` - (Required) Identifier of the resource
    - `value` - (Required) Path to the resource in API.OVH.COM

## Attributes Reference

`id` is set to the order Id. In addition, the following attributes are exported:

- `access` - project access right for the identity that trigger the terraform script.
- `description` - Project description
- `order` - Details about the order that was used to create the public cloud project
  - `date` - date
  - `order_id` - order id, the same as the `id`
  - `expiration_date` - expiration date
  - `details` - Information about a Bill entry
    - `description` - description
    - `order_detail_id` - order detail id
    - `domain` - expiration date
    - `quantity` - quantity
- `project_name` - openstack project name
- `project_id` - openstack project id
- `status` - project status

## Import

Cloud project can be imported using the `order_id` that can be retrieved in the [order page](https://www.ovh.com/manager/#/dedicated/billing/orders/orders) at the creation time of the Public Cloud project.

```bash
$ terraform import ovh_cloud_project.my_cloud_project order_id
```
