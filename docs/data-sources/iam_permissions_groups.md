---
subcategory : "Account Management (IAM)"
---

# ovh_iam_permissions_group (Data Source)

Use this data source to retrieve an IAM permissions group.

## Example Usage

```terraform
data "ovh_iam_permissions_group" "website" {
  urn = "urn:v1:eu:permissionsGroup:ovh:controlPanelAccess"
}
```

## Argument Reference

* `urn` - URN of the permissions group.

## Attributes Reference

* `name` - Name of the permissions group.
* `description` - Group description.
* `allow` - Set of actions allowed by the permissions group.
* `except` - Set of actions that will be subtracted from the `allow` list.
* `deny` - Set of actions that will always be denied even if it is explicitly allowed by a policy.
* `owner` - Owner of the permissions group.
* `created_at` - Creation date of this group.
* `updated_at` - Date of the last update of this group.
* `read_only` - Indicates that this is a default permissions group, managed by OVHcloud.
