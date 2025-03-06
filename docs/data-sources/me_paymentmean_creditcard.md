---
subcategory : "Account Management (IAM)"
---

# ovh_me_paymentmean_creditcard (Data Source)

Use this data source to retrieve information about a credit card payment mean associated with an OVHcloud account.

## Example Usage

```terraform
data "ovh_me_paymentmean_creditcard" "cc" {
  use_default = true
}
```

## Argument Reference

* `description_regexp` - (Optional) a regexp used to filter credit cards on their `description` attributes.

* `use_default` - (Optional) Retrieve credit card marked as default payment mean.

* `use_last_to_expire` - (Optional) Retrieve the credit card that will be the last to expire according to its expiration date.

* `states` - (Optional) Filter credit cards on their `state` attribute. Can be "expired", "valid", "tooManyFailures"

## Attributes Reference

`id` is set to the ID of the credit card payment mean

* `description` - the description attribute of the credit card
* `state` - the state attribute of the credit card
* `default` - a boolean which tells if the retrieved credit card is marked as the default payment mean
