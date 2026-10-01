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
  size         = 5368709120
}
```

## Schema

### Required

- `domain` (String) Name of the email domain.
- `account_name` (String) Name of the email account (without the domain part).
- `password` (String, Sensitive) Password of the email account.

### Optional

- `description` (String) Description of the email account.
- `size` (Number) Size of the email account in bytes.

### Read-Only

- `id` (String) Unique identifier for the resource (`domain/account_name`).
- `email` (String) Full email address of the account.
- `is_blocked` (Boolean) Whether the account is blocked.
