resource "ovh_okms_secret" "example" {
	okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path    = "app/api_credentials"

	# Check‑and‑set parameter used only on update (if `cas_required` metadata is set to true) 
	# to enforce optimistic concurrency control: its value must equal the current secret version (`metadata.current_version`) 
	# for the update to succeed. Ignored on create.
	cas = 1

	metadata = {
		max_versions             = 10         # keep last 10 versions
		cas_required             = true       # enforce optimistic concurrency control (server will require current secret version on the cas attribute to allow update)
		deactivate_version_after = "0s"       # keep versions active indefinitely (example)
		custom_metadata = {
			environment = "prod"
			appname     = "helloworld"
		}
	}

	# Initial version (will create version 1)
	version = {
		data = jsonencode({
			api_key    = "mykey"
			api_secret = "mysecret"
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
