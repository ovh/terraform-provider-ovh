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