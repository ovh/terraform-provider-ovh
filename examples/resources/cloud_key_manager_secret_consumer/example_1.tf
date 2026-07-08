resource "ovh_cloud_key_manager_secret_consumer" "consumer" {
  service_name  = "Public cloud project ID"
  secret_id     = "00000000-0000-0000-0000-000000000000"
  service       = "LOADBALANCER"
  resource_type = "LOADBALANCER"
  resource_id   = "11111111-1111-1111-1111-111111111111"
}
