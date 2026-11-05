data "ovh_cloud_managed_kubernetes" "cluster" {
  service_name = "<Public cloud project id>"
  id           = "<cluster id>"
}

output "api_server_url" {
  value = data.ovh_cloud_managed_kubernetes.cluster.current_state.api_server.access.public_endpoint.url
}
