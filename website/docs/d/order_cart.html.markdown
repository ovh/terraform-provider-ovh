---
subcategory : "Order"
---

# ovh_order_cart (Data Source)

Use this data source to create a temporary order cart to retrieve information order cart products.

## Example Usage

```hcl
data "ovh_order_cart" "mycart" {
 ovh_subsidiary = "fr"
 description    = "my cart"
}
```

## Argument Reference


* `ovh_subsidiary` - (Required) OVHcloud Subsidiary
* `description` - Description of your cart
* `assign` - Assign a shopping cart to an loggedin client. Values can be `true` or `false`. 
* `expire` - Expiration time (format: 2006-01-02T15:04:05+00:00)


## Attributes Reference

`id` is set to your cart ID
In addition, the following attributes are exported.

* `cart_id` - Cart identifier
* `read_only` - Indicates if the cart has already been validated
* `items` - Items of your cart
