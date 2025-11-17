---
subcategory : "Key Management Service (KMS)"
---

# ovh_okms_secret (Data Source)

Retrieves metadata (and optionally the payload) of a secret stored in OVHcloud KMS.

> WARNING: If `include_data = true` the secret value is stored in cleartext (JSON) in the Terraform state file. Marked **Sensitive** only hides it from CLI output. If you use this option it is recommended to protect your state with encryption and access controls.

## Example Usage

Get the latest secret version (metadata only):

```terraform
data "ovh_okms_secret" "latest" {
	okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path    = "app/api_credentials"
}
```

Get the latest secret version including its data:

```terraform
data "ovh_okms_secret" "latest_with_data" {
	okms_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path         = "app/api_credentials"
	include_data = true
}

locals {
	secret_obj = jsondecode(data.ovh_okms_secret.latest_with_data.data)
}

output "api_key" {
	value     = local.secret_obj.api_key
	sensitive = true
}
```

Get a specific version including its payload:

```terraform
data "ovh_okms_secret" "v3" {
	okms_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path         = "app/api_credentials"
	version      = 3
	include_data = true
}
```

## Argument Reference

The following arguments are supported:

### Required

- `okms_id` (String) OKMS service ID that owns the secret.
- `path` (String) Secret path (identifier within the OKMS instance).

### Optional

- `version` (Number) Specific version to retrieve. If omitted, the latest (current) version is selected.
- `include_data` (Boolean) If true, retrieves the secret payload (`data` attribute). Defaults to false. When false only metadata is returned.

## Attributes Reference (Read-Only)

In addition to the arguments above, the following attributes are exported:

- `version` (Number) The resolved version number (requested or current latest).
- `data` (String, Sensitive) Raw JSON secret payload (present only if `include_data` is true).
- `metadata` (Block) Secret metadata:
  - `cas_required` (Boolean)
  - `created_at` (String)
  - `updated_at` (String)
  - `current_version` (Number)
  - `oldest_version` (Number)
  - `max_versions` (Number)
  - `deactivate_version_after` (String)
  - `custom_metadata` (Map of String)
- `iam` (Block) IAM resource metadata:
  - `display_name` (String)
  - `id` (String)
  - `tags` (Map of String)
  - `urn` (String)

## Behavior & Notes

- The `data` attribute retains the raw JSON returned by the API. Use `jsondecode()` to work with individual keys.
- Changing only `include_data` (true -> false) will cause the `data` attribute to become null in subsequent refreshes (state no longer holds the payload).
