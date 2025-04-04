resource "ovh_vrack_ipv6" "vrack_block_v6" {
  service_name = "<vRack service name>"
  block        = "<ipv6 block>"
  bridged_subrange {
    slaac = "<enabled|disabled>"
  }
}

resource "ovh_vrack_ipv6" "vrack_block_v6" {
  service_name = "<vRack service name>"
  block        = "<ipv6 block>"
}