resource "ovh_dedicated_cloud_user" "john" {
  service_name = "pcc-1-2-3-4"
  name         = "john.doe"
}

resource "ovh_dedicated_cloud_user_object_right" "john_cluster_readonly" {
  service_name     = ovh_dedicated_cloud_user.john.service_name
  user_id          = ovh_dedicated_cloud_user.john.user_id
  type             = "cluster"
  vmware_object_id = "domain-c1847"
  right            = "readonly"
  propagate        = true
}
