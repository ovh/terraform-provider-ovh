data "ovh_dedicated_cloud_user_object_right" "john_cluster_readonly" {
  service_name    = "pcc-1-2-3-4"
  user_id         = 42
  object_right_id = 123
}
