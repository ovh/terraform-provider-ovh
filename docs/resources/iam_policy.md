---
subcategory : "Account Management (IAM)"
---

# ovh_iam_policy

Creates an IAM policy.

## Example Usage

```terraform
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

resource "ovh_iam_policy" "ip_prod_access" {
  name        = "ip_prod_access"
  description = "Allow access only from a specific IP to resources tagged prod"
  identities  = [ovh_me_identity_group.my_group.urn]
  resources   = ["urn:v1:eu:resource:vps:*"]

  allow = [
    "vps:apiovh:*",
  ]

  conditions {
    operator = "MATCH"
    values = {
      "resource.Tag(environment)" = "prod"
      "request.IP"                = "192.72.0.1"
    }
  }
}

resource "ovh_iam_policy" "workdays_expiring" {
  name        = "workdays_expiring"
  description = "Allow access only on workdays, expires end of 2026"
  identities  = [ovh_me_identity_group.my_group.urn]
  resources   = ["urn:v1:eu:resource:vps:*"]

  allow = [
    "vps:apiovh:*",
  ]

  conditions {
    operator = "MATCH"
    values = {
      "date(Europe/Paris).WeekDay.In" = "monday,tuesday,wednesday,thursday,friday"
    }
  }

  expired_at = "2026-12-31T23:59:59Z"
}
```

## Argument Reference

* `name` - Name of the policy, must be unique
* `description` - Description of the policy
* `identities` - List of identities affected by the policy
* `resources` - List of resources affected by the policy
* `allow` - List of actions allowed on resources by identities
* `except` - List of overrides of action that must not be allowed even if they are caught by allow. Only makes sens if allow contains wildcards.
* `deny` - List of actions that will always be denied even if also allowed by this policy or another one.
* `permissions_groups` - Set of permissions groups included in the policy. At evaluation, these permissions groups are each evaluated independently (notably, excepts actions only affect actions in the same permission group).
* `expired_at` - (Optional) Expiration date of the policy in RFC3339 format (e.g., `2025-12-31T23:59:59Z`). After this date, the policy will no longer be applied.
* `conditions` - (Optional) Conditions restrict permissions based on resource tags, date/time, or request attributes. See Conditions below.

### Conditions

The `conditions` block supports:

* `operator` - (Required) Operator to combine conditions. Valid values are `AND`, `OR`, `NOT`, or `MATCH`.
* `condition` - (Optional) List of condition blocks. Each condition supports:
  * `operator` - (Required) Operator for this condition (typically `MATCH`).
  * `values` - (Optional) Map of key-value pairs to match. Keys can reference:
    * Resource tags: `resource.Tag(tag_name)` (e.g., `resource.Tag(environment)`)
    * Date/time: `date(timezone).WeekDay`, `date(timezone).WeekDay.In` (e.g., `date(Europe/Paris).WeekDay`)
    * Request attributes: `request.IP`

## Attributes Reference

* `owner` - Owner of the policy.
* `created_at` - Creation date of this group.
* `updated_at` - Date of the last update of this group.
* `read_only` - Indicates that the policy is a default one.
