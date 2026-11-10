---
subcategory : "Email Domain"
---

# ovh_email_domain_account (Resource)

Manages an email account on an OVHcloud email domain (MX Plan / Web Hosting).

## Example Usage

```hcl
resource "ovh_email_domain_account" "my_account" {
  domain       = "example.com"
  account_name = "contact"
  password     = "AStr0ngP@ssw0rd!"
  description  = "Contact email account"
  size         = 5000000000
}
```

## Schema

### Required

- `domain` (String) Name of the email domain.
- `account_name` (String) Name of the email account (without the domain part).

### Optional

- `password` (String, Sensitive) Password of the email account. Required to create an account. The API never returns it, so an imported account has none in state: leave it out to keep the current password, or set it to change it.
- `description` (String) Description of the email account.
- `size` (Number) Size of the email account in bytes.

### Read-Only

- `id` (String) Unique identifier for the resource (`domain/account_name`).
- `email` (String) Full email address of the account.
- `is_blocked` (Boolean) Whether the account is blocked.

## Import

An email account can be imported using its `domain` and `account_name`, separated by a slash:

```hcl
import {
  to = ovh_email_domain_account.my_account
  id = "example.com/contact"
}
```

```bash
terraform import ovh_email_domain_account.my_account example.com/contact
```

The API never returns the password, so it is not imported. Leave `password` out of the configuration to manage the
account without changing it: setting it plans an update that changes the password of the account.
