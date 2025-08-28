data "ovh_okms_secret" "v3" {
	okms_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path         = "app/api_credentials"
	version      = 3
	include_data = true
}
