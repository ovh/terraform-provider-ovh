---
subcategory : "Account Management"
---

# ovh_me (Data Source)

Use this data source to get information about the current OVHcloud account.

## Example Usage

```hcl
data "ovh_me" "myaccount" {}
```

## Argument Reference

There are no arguments to this datasource.

## Attributes Reference

The following attributes are exported:

* `address`: Postal address of the account
* `area`: Area of the account
* `birth_city`: City of birth
* `birth_day`: Birth date
* `city`: City of the account
* `company_national_identification_number`: This is the national identification number of the company that possess this account
* `corporation_type`: Type of corporation
* `country`: Country of the account
* `currency`:
  * `code`: Currency code used by this account (e.g EUR, USD, ...)
  * `symbol`: Currency symbol used by this account (e.g €, $, ...)
* `customer_code`: The customer code of this account (a numerical value used for identification when contacting support via phone call)
* `email`: Email address
* `fax`: Fax number
* `firstname`: First name
* `italian_sdi`: Italian SDI
* `language`: Preferred language for this account
* `legalform`: Legal form of the account
* `name`: Name of the account holder
* `national_identification_number`: National Identification Number of this account
* `nichandle`: Nic handle / customer identifier
* `organisation`: Name of the organisation for this account
* `ovh_company`: OVHcloud subsidiary
* `ovh_subsidiary`: OVHcloud subsidiary
* `phone`: Phone number
* `phone_country`: Country code of the phone number
* `sex`: Gender of the account holder
* `spare_email`: Backup email address
* `state`: State of the postal address
* `vat`: VAT number
* `zip`: Zipcode of the address
