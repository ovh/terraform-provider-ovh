resource "ovh_account_logs" "activity_logs" {
  log_type      = "activity"
  stream_id     = "xxx"
  kind          = "default"
}

resource "ovh_account_logs" "access_policy_logs" {
  log_type      = "access_policy"
  stream_id     = "xxx"
  kind          = "default"
}
