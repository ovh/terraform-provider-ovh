resource "ovh_cloud_project_region_loadbalancer_log_subscription" "subscription" {
  service_name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  region_name = "yyyy"
  loadbalancer_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  kind = "haproxy"
  stream_id = "ffffffff-gggg-hhhh-iiii-jjjjjjjjjjjj"
}
