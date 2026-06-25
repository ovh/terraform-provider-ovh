data "ovh_cloud_key_manager_container_consumer" "consumer" {
  service_name = "Public cloud project ID"
  container_id = "00000000-0000-0000-0000-000000000000"
  consumer_id  = "Q09NUFVURTpJTlNUQU5DRToxMTExMTExMS0xMTExLTExMTEtMTExMS0xMTExMTExMTExMTE"
}

output "consumer_service" {
  value = data.ovh_cloud_key_manager_container_consumer.consumer.service
}
