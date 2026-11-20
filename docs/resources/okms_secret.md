---
subcategory : "Key Management Service (KMS)"
---

# ovh_okms_secret (Resource)

Manages a secret stored in OVHcloud KMS.

> WARNING: `version.data` is marked **Sensitive** but still ends up in the state file. To keep the secret data out of the state, use the write-only `version.data_wo` instead (Terraform 1.11 or later).

## Example Usage

Create a secret whose value is a JSON object. Use `jsonencode()` to produce a deterministic JSON string (ordering/whitespace) to minimize diffs.

```terraform
resource "ovh_okms_secret" "example" {
	okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path    = "app/api_credentials"

	metadata = {
		max_versions             = 10         # keep last 10 versions
		cas_required             = true       # enforce optimistic concurrency control (server will require current secret version on the cas attribute to allow update)
		deactivate_version_after = "0s"       # keep versions active indefinitely (example)
		custom_metadata = {
			environment = "prod"
			owner       = "payments-team"
		}
	}

	# Initial version (will create version 1)
	version = {
		data = jsonencode({
			api_key    = var.api_key
			api_secret = var.api_secret
		})
	}
}

# Reading a field from the secret version data
locals {
	secret_json = jsondecode(ovh_okms_secret.example.version.data)
}

output "api_key" {
	value     = local.secret_json.api_key
	sensitive = true
}
```

# Reading a field from the secret version data:

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

Updating the secret (new version) while enforcing optimistic concurrency control using CAS:

```terraform
resource "ovh_okms_secret" "example" {
	okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path    = "app/api_credentials"

	# Ensure no concurrent update happened: set cas to the current version
	# (metadata.current_version is populated after first apply)
	cas = ovh_okms_secret.example.metadata.current_version

	metadata = {
		cas_required = true
	}

	version = {
		data = jsonencode({
			api_key    = var.api_key
			api_secret = var.new_api_secret  # changed value -> creates new version
		})
	}
}
```

Keeping the secret data out of the Terraform state with the write-only `version.data_wo` (Terraform 1.11 or later):

```terraform
# The secret data is write-only: it is sent to the KMS but never stored in the
# Terraform state or plan. Increment data_wo_version to send a new version.
ephemeral "random_password" "api_secret" {
	length = 32
}

resource "ovh_okms_secret" "write_only" {
	okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path    = "app/api_credentials"

	version = {
		data_wo = jsonencode({
			api_secret = ephemeral.random_password.api_secret.result
		})
		data_wo_version = 1
	}
}

# A secret created by Terraform and filled outside of it: its data is never read
# back into the state, nor overwritten while data_wo_version is unchanged.
resource "ovh_okms_secret" "filled_outside" {
	okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path    = "app/env"

	version = {
		data_wo         = jsonencode({})
		data_wo_version = 1
	}
}
```

## Argument Reference

The following arguments are supported:

### Required

- `okms_id` (String) ID of the OKMS service to create the secret in.
- `path` (String) Secret path (acts as the secret identifier within the OKMS instance). Immutable after creation.
- `version` (Block) Definition of the version to create/update. See Version Block below. (On updates providing a new `version.data` creates a new version.)

### Optional

- `cas` (Number) Check‑and‑set parameter used only on update (if `cas_required` metadata is set to true) to enforce optimistic concurrency control: its value must equal the current secret version (`metadata.current_version`) for the update to succeed. Ignored on create.
- `metadata` (Block) Secret metadata configuration (subset of fields are user-settable). See Metadata Block below.

### Metadata Block

User configurable attributes inside `metadata`:

- `cas_required` (Boolean) If true, the server will enforce optimistic concurrency control by requiring the `cas` parameter to match the current version number on every write (update) request.
- `custom_metadata` (Map of String) Arbitrary key/value metadata.
- `deactivate_version_after` (String) Duration (e.g. `"24h"`) after which a version is deactivated. `"0s"` (default) means never automatically deactivate.
- `max_versions` (Number) Number of versions to retain (default 10). Older versions beyond this limit are pruned.

Computed (read‑only) metadata attributes:

- `created_at` (String) Creation timestamp of the secret.
- `updated_at` (String) Last update timestamp.
- `current_version` (Number) Current (latest) version number.
- `oldest_version` (Number) Oldest retained version number.

### Version Block

Exactly one of `data` or `data_wo` must be set:

- `data` (String, Sensitive) Secret payload, stored in the Terraform state. Commonly set with `jsonencode(...)` so that Terraform comparisons are stable. Any valid JSON (object, array, string, number, bool) is accepted. Changing `data` creates a new secret version.
- `data_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Secret payload, never stored in the Terraform state or plan. Requires Terraform 1.11 or later, and `data_wo_version`. Same format as `data`.
- `data_wo_version` (Number) Version of `data_wo`, required with it. As `data_wo` is not stored, changing it alone is not detected: increment `data_wo_version` to create a new secret version from `data_wo`.

Computed (read‑only) attributes:

- `id` (Number) Version number.
- `created_at` (String) Version creation timestamp.
- `deactivated_at` (String) Deactivation timestamp if the version was deactivated.
- `state` (String) Version state (e.g. `ACTIVE`).

## Attributes Reference (Read-Only)

In addition to arguments above, the following attributes are exported:

- `iam` (Block) IAM metadata: `display_name`, `id`, `tags`, `urn`.
- `metadata.*` computed fields as listed above.
- `version.*` computed fields as listed above.

## Behavior & Notes

- Updating with a new `version.data` performs an API PUT that creates a new version; the previous version remains (subject to `max_versions`).
- With `version.data_wo`, the provider never reads the secret data back: there is no drift detection on the data, and a version created outside of Terraform is neither stored in the state nor overwritten. A new version is created only when `data_wo_version` changes.
- Moving an existing secret from `data` to `data_wo` (or setting `data_wo_version` on it for the first time) adopts the current version as is: nothing is written to the KMS, and the data is removed from the state. Remove any `lifecycle { ignore_changes = [version] }` at the same time, otherwise the change is not applied.
- If `cas_required` is true, all write operations must include a correct `cas` query parameter (the expected current version number). Set `cas = ovh_okms_secret.example.metadata.current_version` to enforce optimistic concurrency. A mismatch causes the API to reject the update (preventing overwriting unseen changes).
- `cas` is ignored on create (no existing version).
