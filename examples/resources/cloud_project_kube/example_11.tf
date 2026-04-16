resource "ovh_cloud_project_kube" "my_multizone_cluster" {
  service_name  = ovh_cloud_project_network_private.network.service_name
  name          = "custom-cilium-configuration-mks"
  region        = "RBX-A"
  plan          = "standard"

  ip_allocation_policy {
    pods_ipv4_cidr = "10.7.0.0/16"
    services_ipv4_cidr = "10.6.0.0/16"
  }

  customization_cilium {
    cluster_id = 1
    cluster_mesh {
      enabled = true
      api_server {
        service_type = "LoadBalancer"
      }
    }
    hubble {
      enabled = true # default true
      relay {
        enabled = true # default true
      }
      ui {
        enabled = true # default true
        backend_resources { # no limits or requests by default
          limits {
            cpu = "500m"
            memory = "500Mi"
          }
          requests {
            cpu = "50m"
            memory = "50Mi"
          }
        }
        frontend_resources {  # no limits or requests by default
          limits {
            cpu = "500m"
            memory = "500Mi"
          }
          requests {
            cpu = "50m"
            memory = "50Mi"
          }
        }
      }
    }

  private_network_id = tolist(ovh_cloud_project_network_private.network.regions_attributes[*].openstackid)[0]
  nodes_subnet_id    = ovh_cloud_project_network_private_subnet.subnet.id

}