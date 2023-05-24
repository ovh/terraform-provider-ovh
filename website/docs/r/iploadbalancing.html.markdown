---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing

Orders an IP load balancing.

## Important

~> __WARNING__ This resource orders an OVHcloud product for a long period of time and may generate heavy costs!
Use with caution.

-> __NOTE__ To order a product through Terraform, your account needs to have a default payment method defined. This can be done in the [OVHcloud Control Panel](https://www.ovh.com/manager/#/dedicated/billing/payment/method) or via API with the [/me/payment/method](https://api.ovh.com/console/#/me/payment/method~GET) endpoint.


## Example Usage

```hcl
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "mycart"
}

data "ovh_order_cart_product_plan" "iplb" {
  cart_id        = data.ovh_order_cart.mycart.id
  price_capacity = "renew"
  product        = "ipLoadbalancing"
  plan_code      = "iplb-lb1"
}

data "ovh_order_cart_product_options_plan" "bhs" {
  cart_id           = data.ovh_order_cart_product_plan.iplb.cart_id
  price_capacity    = data.ovh_order_cart_product_plan.iplb.price_capacity
  product           = data.ovh_order_cart_product_plan.iplb.product
  plan_code         = data.ovh_order_cart_product_plan.iplb.plan_code
  options_plan_code = "iplb-zone-lb1-rbx"
}

resource "ovh_iploadbalancing" "iplb-lb1" {
  ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
  display_name   = "my ip loadbalancing"

  plan {
    duration     = data.ovh_order_cart_product_plan.iplb.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_plan.iplb.plan_code
    pricing_mode = data.ovh_order_cart_product_plan.iplb.selected_price.0.pricing_mode
  }

  plan_option {
    duration     = data.ovh_order_cart_product_options_plan.bhs.selected_price.0.duration
    plan_code    = data.ovh_order_cart_product_options_plan.bhs.plan_code
    pricing_mode = data.ovh_order_cart_product_options_plan.bhs.selected_price.0.pricing_mode
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - Set the name displayed in ManagerV6 for your iplb (max 50 chars)
* `ovh_subsidiary` - (Required) OVHcloud Subsidiary
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
* `ssl_configuration` - Modern oldest compatible clients : Firefox 27, Chrome 30, IE 11 on Windows 7, Edge, Opera 17, Safari 9, Android 5.0, and Java 8. Intermediate oldest compatible clients : Firefox 1, Chrome 1, IE 7, Opera 5, Safari 1, Windows XP IE8, Android 2.3, Java 7. Intermediate if null. one of "intermediate", "modern". 


## Attributes Reference

Id is set to the order Id. In addition, the following attributes are exported:
* `urn` - URN of the load balancer, used when writing IAM policies
* `ip_loadbalancing` - Your IP load balancing
* `ipv4` - The IPV4 associated to your IP load balancing
* `ipv6` - The IPV6 associated to your IP load balancing. DEPRECATED.
* `metrics_token` - The metrics token associated with your IP load balancing
* `offer` - The offer of your IP load balancing
* `order` - Details about an Order
  * `date` - date
  * `order_id` - order id
  * `expiration_date` - expiration date
  * `details` - Information about a Bill entry
    * `description` - description
    * `order_detail_id` - order detail id
    * `domain` - expiration date
    * `quantity` - quantity
* `orderable_zone` - Available additional zone for your Load Balancer
  * `name` - The zone three letter code
  * `plan_code` - The billing planCode for this zone
* `service_name` - The internal name of your IP load balancing
* `state` - Current state of your IP
* `vrack_eligibility` - Vrack eligibility
* `vrack_name` - Name of the vRack on which the current Load Balancer is attached to, as it is named on vRack product
* `zone` - Location where your service is
