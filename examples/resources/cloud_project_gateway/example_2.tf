resource "ovh_cloud_project_gateway" "imported_gateway" {
  service_name = ovh_cloud_project_network_private.mypriv.service_name
  name         = "<my-imported-gateway>"
  model        = "<my-model>"
  region       = "<my-region>"
  network_id   = "<my-imported-gateway-network-id>"
  subnet_id    = "<my-imported-gateway-subnet-id>"
  lifecycle {
    ignore_changes = [network_id, subnet_id]
  }
}

import {
  id = "<service-name>/<region>/<gateway-id>"
  to = ovh_cloud_project_gateway.imported_gateway
}
