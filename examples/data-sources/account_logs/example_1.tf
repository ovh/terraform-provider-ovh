data "ovh_account_logs" "audit_logs" {
  log_type        = "audit"
  subscription_id = "xxx"
}
