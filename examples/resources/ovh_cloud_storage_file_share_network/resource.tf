resource "ovh_cloud_storage_file_share_network" "network" {
  service_name = "<public cloud project ID>"
  name         = "my-share-network"
  description  = "Share network for my NFS shares"
  network_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  region       = "GRA1"
}
