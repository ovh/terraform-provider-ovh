resource "ovh_dbaas_logs_encryption_key" "key" {
  service_name = "ldp-xx-xxxxx"
  title        = "my PGP key"
  content      = file("my-pgp-public-key.asc")
  fingerprint  = "ABCDEF1234567890ABCDEF1234567890ABCDEF12"
}

resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name            = "ldp-xx-xxxxx"
  title                   = "my stream"
  description             = "my encrypted graylog stream"
  cold_storage_enabled    = true
  cold_storage_target     = "PCA"
  cold_storage_retention  = 1
  encryption_keys_ids     = [ovh_dbaas_logs_encryption_key.key.id]
}
