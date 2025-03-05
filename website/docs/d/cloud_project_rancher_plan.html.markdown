---
subcategory : "Managed Rancher Service (MRS)"
---

# ovh_cloud_project_rancher_plan

Use this datasource to retrieve information about the Managed Rancher plans available in the given public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_rancher_plan" "plans" {
  project_id = "XXXXXX"
}
```

## Schema

### Required

- `project_id` (String) Project ID

### Read-Only

- `plans` (Attributes Set) (see [below for nested schema](#nestedatt--plans))

<a id="nestedatt--plans"></a>
### Nested Schema for `plans`

Read-Only:

- `cause` (String) Cause for an unavailability
- `message` (String) Human-readable description of the unavailability cause
- `name` (String) Name of the plan
- `status` (String) Status of the plan
