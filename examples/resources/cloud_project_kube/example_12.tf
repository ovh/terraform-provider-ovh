# The nodes only exist once the node pool is created: read them from another
# configuration, or apply twice
data "ovh_cloud_project_kube_nodes" "nodes" {
  service_name = ovh_cloud_project_kube.my_cluster.service_name
  kube_id      = ovh_cloud_project_kube.my_cluster.id
}

data "ovh_cloud_instance" "node" {
  for_each     = { for node in data.ovh_cloud_project_kube_nodes.nodes.nodes : node.name => node.instance_id }
  service_name = ovh_cloud_project_kube.my_cluster.service_name
  id           = each.value
}

# Public IPv4 of each node, from its public (Ext-Net) interface
output "nodes_egress_ips" {
  value = flatten([
    for instance in data.ovh_cloud_instance.node : [
      for network in instance.current_state.networks : [
        for address in network.addresses : address.ip if address.version == 4
      ] if network.id == null || network.id == ""
    ]
  ])
}
