---
subcategory : "Account Management"
---

# ovh_iam_policy

Creates an IAM policy.

## Example Usage

```hcl
data "ovh_me" "account" {}

resource "ovh_me_identity_group" "my_group" {
  name        = "my_group"
  description = "my_group created in Terraform"
}

resource "ovh_iam_policy" "manager" {
  name        = "allow_ovh_manager"
  description = "Users are allowed to use the OVH manager"
  identities  = [ovh_me_identity_group.my_group.urn]
  resources   = [data.ovh_me.account.urn]
  # these are all the actions 
  allow = [
    "account:apiovh:me/get",
    "account:apiovh:me/supportLevel/get",
    "account:apiovh:me/certificates/get",
    "account:apiovh:me/tag/get",
    "account:apiovh:services/get",
    "account:apiovh:*",
  ]
}
```

## Argument Reference

* `name` - Name of the policy, must be unique
* `description` - Description of the policy
* `identities` - List of identities affected by the policy
* `resources` - List of resources affected by the policy
* `allow` - List of actions allowed on resources by identities
* `except` - List of overrides of action that must not be allowed even if they are caught by allow. Only makes sens if allow contains wildcards.

## Attributes Reference

* `owner` - Owner of the policy.
* `created_at` - Creation date of this group.
* `updated_at` - Date of the last update of this group.
* `read_only` - Indicates that the policy is a default one.
