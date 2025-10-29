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
