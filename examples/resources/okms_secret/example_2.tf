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
