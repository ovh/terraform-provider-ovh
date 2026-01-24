---
subcategory : "Account Management (IAM)"
---

# ovh_me_identity_user_token

Creates a token for an identity user.

## Example Usage

```terraform
resource "ovh_me_identity_user" "user" {
  description = "User description"
  email       = "user.email@example.com"
  group       = "DEFAULT"
  login       = "user_login"
  password    = "SecretPassword123"
}

resource "ovh_me_identity_user_token" "token" {
  user_login  = ovh_me_identity_user.user.login
  name        = "token_name"
  description = "Token description"
  expires_at  = "2030-01-01T00:00:00Z"
}
```

## Argument Reference

* `user_login` - (Required) User's login suffix.
* `name` - (Required) Token name.
* `description` - (Required) Token description.
* `expires_at` - (Optional) Token expiration date.
* `expires_in` - (Optional) Token validity duration in seconds.

## Attributes Reference

* `token` - The token value.
* `creation` - Creation date of this token.
* `last_used` - Last use of this token.
