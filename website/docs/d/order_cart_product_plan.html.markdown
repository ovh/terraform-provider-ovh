---
subcategory : "Order"
---

# ovh_order_cart_product_plan (Data Source)

Use this data source to retrieve information of order cart product plan.

## Example Usage

```hcl
data "ovh_order_cart" "mycart" {
  ovh_subsidiary = "fr"
  description    = "my cart"
}

data "ovh_order_cart_product_plan" "plan" {
  cart_id           = data.ovh_order_cart.mycart.id
  price_capacity    = "renew"
  product           = "cloud"
  plan_code         = "project"
}
```

## Argument Reference

* `cart_id` - (Required) Cart identifier
* `catalog_name` - Catalog name
* `plan_code` - (Required) Product offer identifier
* `price_capacity` - (Required) Capacity of the pricing (type of pricing)
* `product` - (Required) Product

## Attributes Reference

The following attributes are exported.

* `selected_price` - Selected Price according to capacity
  * `capacities` - Capacities of the pricing (type of pricing)
  * `description` - Description of the pricing
  * `duration` - Duration for ordering the product
  * `interval` - Interval of renewal
  * `maximum_quantity` - Maximum quantity that can be ordered
  * `maximum_repeat` - Maximum repeat for renewal
  * `minimum_quantity` - Minimum quantity that can be ordered
  * `minimum_repeat` - Minimum repeat for renewal
  * `price` - Price of the product (Price with its currency and textual representation)
    * `currency_code` - Currency code
    * `text` - Textual representation
    * `value` - The effective price
  * `price_in_ucents` - Price of the product in micro-centims
  * `pricing_mode` - Pricing model identifier
  * `pricing_type` - Pricing type
* `plan_code` - Product offer identifier
* `product_name` - Name of the product
* `product_type` - Product type
* `prices` - Prices of the product offer
  * `capacities` - Capacities of the pricing (type of pricing)
  * `description` - Description of the pricing
  * `duration` - Duration for ordering the product
  * `interval` - Interval of renewal
  * `maximum_quantity` - Maximum quantity that can be ordered
  * `maximum_repeat` - Maximum repeat for renewal
  * `minimum_quantity` - Minimum quantity that can be ordered
  * `minimum_repeat` - Minimum repeat for renewal
  * `price` - Price of the product (Price with its currency and textual representation)
    * `currency_code` - Currency code
    * `text` - Textual representation
    * `value` - The effective price
  * `price_in_ucents` - Price of the product in micro-centims
  * `pricing_mode` - Pricing model identifier
  * `pricing_type` - Pricing type
