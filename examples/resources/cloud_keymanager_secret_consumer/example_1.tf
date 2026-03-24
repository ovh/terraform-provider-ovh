resource "ovh_cloud_keymanager_secret_consumer" "consumer" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  secret_id     = "00000000-0000-0000-0000-000000000000"
  service       = "LOADBALANCER"
  resource_type = "LOADBALANCER"
  resource_id   = "11111111-1111-1111-1111-111111111111"
}
