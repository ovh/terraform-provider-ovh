---
subcategory : "Account Management"
---

# ovh_iam_permissions_group

Create am IAM permissions group.

## Example Usage

```hcl
# create a group allowing all actions in the category READ on VPSs
resource "ovh_iam_permissions_group" "read_vps" {
  name        = "read_vps"
  description = "Read access to vps"

  allow = [
    for act in data.ovh_iam_reference_actions.vps.actions : act.action if(contains(act.categories, "READ"))
  ]
}

data "ovh_iam_reference_actions" "vps" {
  type = "vps"
}
```

## Argument Reference

* `name` - Name of the permissions group.
* `description` - Group description.
* `allow` - Set of actions allowed by the permissions group.
* `except` - Set of actions that will be subtracted from the `allow` list.
* `deny` - Set of actions that will be denied no matter what permissions group exists.

## Attributes Reference

* `urn` - URN of the permissions group.
* `owner` - Owner of the permissions group.
* `created_at` - Creation date of this group.
* `updated_at` - Date of the last update of this group.
* `read_only` - Indicates that the permissions group is a default one.
