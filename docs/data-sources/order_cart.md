---
subcategory : "Order"
---

# ovh_order_cart (Data Source)

Use this data source to create a temporary order cart to retrieve information order cart products.

## Example Usage

```terraform
data "ovh_me" "my_account" {}

data "ovh_order_cart" "my_cart" {
 ovh_subsidiary = data.ovh_me.my_account.ovh_subsidiary
}
```

## Argument Reference

* `ovh_subsidiary` - (Required) OVHcloud Subsidiary. Country of OVHcloud legal entity you'll be billed by. List of supported subsidiaries available on API at [/1.0/me.json under `models.nichandle.OvhSubsidiaryEnum`](https://eu.api.ovh.com/1.0/me.json)
* `description` - Description of your cart
* `assign` - Assign a shopping cart to a logged in client. Values can be `true` or `false`.
* `expire` - Expiration time (format: 2006-01-02T15:04:05+00:00)

## Attributes Reference

`id` is set to your cart ID In addition, the following attributes are exported.

* `cart_id` - Cart identifier
* `read_only` - Indicates if the cart has already been validated
* `items` - Items of your cart
