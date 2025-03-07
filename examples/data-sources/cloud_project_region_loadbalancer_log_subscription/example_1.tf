data "ovh_cloud_project_region_loadbalancer_log_subscription" "sub" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region_name = "gggg"
  loadbalancer_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  subscription_id = "zzzzzzzz-yyyy-xxxx-wwww-vvvvvvvvvvvv"
}
