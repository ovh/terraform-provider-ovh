---
subcategory : "Account Management"
---

# ovh_iam_policy (Data Source)

Use this data source to retrieve am IAM policy.

## Example Usage

```hcl
data "ovh_iam_policy" "my_policy" {
  id = "my_policy_id"
}
```

## Argument Reference

* `id` - UUID of the policy.

## Attributes Reference

* `name` - Name of the policy.
* `description` - Group description.
* `identities` - List of identities affected by the policy.
* `resources` - List of resources affected by the policy.
* `allow` - List of actions allowed by the policy.
* `except` - List of actions.
* `owner` - Owner of the policy.
* `created_at` - Creation date of this group.
* `updated_at` - Date of the last update of this group.
* `read_only` - Indicates that the policy is a default one.
