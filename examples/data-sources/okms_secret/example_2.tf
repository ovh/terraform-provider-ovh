data "ovh_okms_secret" "latest" {
	okms_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	path    = "app/api_credentials"
}
