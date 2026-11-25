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
