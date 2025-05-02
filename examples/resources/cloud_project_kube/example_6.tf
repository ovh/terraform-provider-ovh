resource "ovh_cloud_project_kube" "my_cluster" {
  service_name    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name            = "my_kube_cluster"
  region          = "GRA11"
  kube_proxy_mode = "ipvs" # or "iptables"	
	
  customization_kube_proxy {
    iptables {
      min_sync_period = "PT0S"
      sync_period = "PT0S"
    }
        
    ipvs {
      min_sync_period = "PT0S"
      sync_period = "PT0S"
      scheduler = "rr"
      tcp_timeout = "PT0S"
      tcp_fin_timeout = "PT0S"
      udp_timeout = "PT0S"
    }
  }
}
