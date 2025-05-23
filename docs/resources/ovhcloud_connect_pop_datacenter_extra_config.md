---
subcategory: "Ovhcloud Connect (OCC)"
---

# ovh_ovhcloud_connect_pop_datacenter_extra_config (Resource)

Creates an extra datacenter configuration for an Ovhcloud Connect product.

Please take a look at the list of available `types` in the `Required` section in order to know the list of available type configurations.

## Example Usage

```terraform
data "ovh_ovhcloud_connect_config_pops" "pop_cfgs" {
  service_name = "XXX"
}

data "ovh_ovhcloud_connect_config_pop_datacenters" "datacenter_cfgs" {
  service_name = data.ovh_ovhcloud_connect_config_pops.pop_cfgs.service_name
  config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
}

resource "ovh_ovhcloud_connect_pop_datacenter_extra_config" "extra" {
    service_name = data.ovh_ovhcloud_connect_config_pops.pop_cfgs.service_name
    config_pop_id = tolist(data.ovh_ovhcloud_connect_config_pops.pop_cfgs.pop_configs)[0].id
    config_datacenter_id = tolist(data.ovh_ovhcloud_connect_config_pop_datacenters.datacenter_cfgs.datacenter_configs)[0].id
    type = "network"
    next_hop = "P.P.P.P"
    subnet = "I.I.I.I/M"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config_datacenter_id` (Number) ID of the datacenter configuration
- `config_pop_id` (Number) ID of the pop configuration
- `service_name` (String) Service name
- `type` (String) Type of the configuration. Availaible types:
    * `bgp`
    * `network`

### Optional

- `bgp_neighbor_area` (Number) BGP AS number
- `bgp_neighbor_ip` (String) Router IP for BGP
- `next_hop` (String) Static route next hop
- `subnet` (String) Static route ip

### Read-Only

- `id` (Number) ID of the extra configuration
- `status` (String) Status of the pop configuration
