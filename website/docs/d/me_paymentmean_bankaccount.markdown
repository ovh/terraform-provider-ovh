---
layout: "ovh"
page_title: "OVH: me_paymentmean_bankaccount"
sidebar_current: "docs-ovh-datasource-me-paymentmean-bankaccount"
description: |-
  Get information & status of an ovh bank account payment mean
---

# ovh_me_paymentmean_bankaccount

Use this data source to retrieve information about a bank account
payment mean associated with an OVH account.

## Example Usage

```hcl
data "ovh_me_paymentmean_bankaccount" "ba" {
   use_default = true
}
```

## Argument Reference


* `description_regexp` - (Optional) a regexp used to filter bank accounts 
on their `description` attributes.

* `use_default` - (Optional) Retrieve bank account marked as default payment mean.

* `use_oldest` - (Optional) Retrieve oldest bank account.
project.

* `state` - (Optional) Filter bank accounts on their `state` attribute.
Can be "blockedForIncidents", "valid", "pendingValidation"


## Attributes Reference

`id` is set to the ID of the bank account payment mean

* `description` - the description attribute of the bank account
* `default` - a boolean which tells if the retrieved bank account
is marked as the default payment mean
