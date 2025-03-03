---
subcategory : "Managed Rancher Service"
---

# ovh_cloud_project_rancher

Retrieve information about a Managed Rancher Service in the given public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_rancher" "rancher" {
  project_id = "<public cloud project ID>"
  id         = "<Rancher service ID>"
}
```

## Schema

### Required

- `id` (String) Unique identifier
- `project_id` (String) Project ID

### Read-Only

- `created_at` (String) Date of the managed Rancher service creation
- `current_state` (Attributes) Current configuration applied to the managed Rancher service (see [below for nested schema](#nestedatt--current_state))
- `current_tasks` (Attributes List) Asynchronous operations ongoing on the managed Rancher service (see [below for nested schema](#nestedatt--current_tasks))
- `resource_status` (String) Reflects the readiness of the managed Rancher service. A new target specification request will be accepted only in `READY` status
- `target_spec` (Attributes) Last target specification of the managed Rancher service (see [below for nested schema](#nestedatt--target_spec))
- `updated_at` (String) Date of the last managed Rancher service update

<a id="nestedatt--current_state"></a>
### Nested Schema for `current_state`

Read-Only:

- `bootstrap_password` (String, Sensitive) Bootstrap password of the managed Rancher service, returned only on creation
- `ip_restrictions` (Attributes List) List of allowed CIDR blocks for a managed Rancher service's IP restrictions. When empty, any IP is allowed (see [below for nested schema](#nestedatt--current_state--ip_restrictions))
- `name` (String) Name of the managed Rancher service
- `networking` (Attributes) Networking properties of a managed Rancher service (see [below for nested schema](#nestedatt--current_state--networking))
- `plan` (String) Plan of the managed Rancher service
- `region` (String) Region of the managed Rancher service
- `url` (String) URL of the managed Rancher service
- `usage` (Attributes) Latest metrics regarding the usage of the managed Rancher service (see [below for nested schema](#nestedatt--current_state--usage))
- `version` (String) Version of the managed Rancher service

<a id="nestedatt--current_state--ip_restrictions"></a>
### Nested Schema for `current_state.ip_restrictions`

Read-Only:

- `cidr_block` (String) Allowed CIDR block (/subnet is optional, if unspecified then /32 will be used)
- `description` (String) Description of the allowed CIDR block


<a id="nestedatt--current_state--networking"></a>
### Nested Schema for `current_state.networking`

Read-Only:

- `egress_cidr_blocks` (List of String) Specifies the CIDR ranges for egress IP addresses used by Rancher. Ensure these ranges are allowed in any IP restrictions for services that Rancher will access.


<a id="nestedatt--current_state--usage"></a>
### Nested Schema for `current_state.usage`

Read-Only:

- `datetime` (String) Date of the sample
- `orchestrated_vcpus` (Number) Total number of vCPUs orchestrated by the managed Rancher service through the downstream clusters



<a id="nestedatt--current_tasks"></a>
### Nested Schema for `current_tasks`

Read-Only:

- `id` (String) Identifier of the current task
- `link` (String) Link to the task details
- `status` (String) Current global status of the current task
- `type` (String) Type of the current task


<a id="nestedatt--target_spec"></a>
### Nested Schema for `target_spec`

Read-Only:

- `ip_restrictions` (Attributes List) List of allowed CIDR blocks for a managed Rancher service's IP restrictions. When empty, any IP is allowed (see [below for nested schema](#nestedatt--target_spec--ip_restrictions))
- `name` (String) Name of the managed Rancher service
- `plan` (String) Plan of the managed Rancher service. Available plans for an existing managed Rancher can be retrieved using GET /rancher/rancherID/capabilities/plan
- `version` (String) Version of the managed Rancher service. Available versions for an existing managed Rancher can be retrieved using GET /rancher/rancherID/capabilities/version

<a id="nestedatt--target_spec--ip_restrictions"></a>
### Nested Schema for `target_spec.ip_restrictions`

Read-Only:

- `cidr_block` (String) Allowed CIDR block (/subnet is optional, if unspecified then /32 will be used)
- `description` (String) Description of the allowed CIDR block
