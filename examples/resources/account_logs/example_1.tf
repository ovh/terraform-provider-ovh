resource "ovh_account_logs" "audit_logs" {
  log_type      = "audit"
  stream_id     = "xxx"
  kind          = "default"
}
