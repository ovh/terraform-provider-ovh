---
subcategory : "Account Management"
---

# ovh_iam_reference_actions (Data Source)

Use this data source to list the IAM action associated with a resource type.

## Example Usage

```hcl
data "ovh_iam_reference_actions" "vps_actions" {
    resource_type = "vps"
}
```

## Argument Reference

* `type` - Kind of resource we want the actions for

## Attributes Reference

* `actions` - List of actions
    * `action` - Name of the action
    * `categories` - List of the categories of the action
    * `description` - Description of the action
    * `resource_type` - Resource type the action is related to