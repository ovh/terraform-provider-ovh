---
subcategory : "Managed Rancher Service (MRS)"
---

# ovh_cloud_project_rancher_version

Use this datasource to retrieve information about the Managed Rancher available versions in the given public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_rancher_version" "versions" {
  project_id = "XXXXXX"
}
```

## Schema

### Required

- `project_id` (String) Project ID

### Read-Only

- `versions` (Attributes Set) (see [below for nested schema](#nestedatt--versions))

<a id="nestedatt--versions"></a>
### Nested Schema for `versions`

Read-Only:

- `cause` (String) Cause for an unavailability
- `changelog_url` (String) Changelog URL of the version
- `message` (String) Human-readable description of the unavailability cause
- `name` (String) Name of the version
- `status` (String) Status of the version
