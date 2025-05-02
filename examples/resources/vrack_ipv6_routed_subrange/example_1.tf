resource "ovh_vrack_ipv6_routed_subrange" "vrack_routed_subrange" {
  service_name = "<vRack service name>"
  block = "<ipv6 block>"
  routed_subrange = "<routed subrange>"
  nexthop = "<nexthop>"
}
