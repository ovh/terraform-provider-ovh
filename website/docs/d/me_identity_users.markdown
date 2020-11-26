---
layout: "ovh"
page_title: "OVH: me_identity_users"
sidebar_current: "docs-ovh-datasource-identity-users"
description: |-
  Get the list of the identity users for the account.
---

# ovh_me_identity_users

Use this data source to retrieve list of user logins of the account's identity users.

## Example Usage

```hcl
data "ovh_me_identity_users" "users" {}
```

## Argument Reference

This datasource takes no arguments.

## Attributes Reference

* `users` - The list of the user's logins of all the identity users.
