resource "ovh_cloud_project_alerting" "my_alert" {
  service_name = "XXX"
  delay = 3600
  email = "aaa.bbb@domain.com"
  monthly_threshold = 1000
}
