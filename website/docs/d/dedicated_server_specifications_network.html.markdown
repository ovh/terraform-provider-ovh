---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_specifications_network (Data Source)

Use this data source to get the network information about a dedicated server associated with your OVHcloud Account.

## Example Usage

```hcl
data "ovh_dedicated_server_specifications_network" "spec" {
  service_name = "myserver"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your dedicated server.

## Attributes Reference

The following attributes are exported:

* `bandwidth` - Bandwidth details
  * `internet_to_ovh` - Bandwidth limitation Internet to OVH
    * `unit`
    * `value`
  * `ovh_to_internet` - Bandwidth limitation OVH to Internet
    * `unit`
    * `value`
  * `ovh_to_ovh` - Bandwidth limitation OVH to OVH
    * `unit`
    * `value`
  * `type` - Bandwidth offer type
* `connection_val` - Network connection flow rate
  * `unit`
  * `value`
* `ola` - OLA details
  * `available` - Is the OLA feature available
  * `available_modes` - Supported modes
    * `default` - Whether it is the default configuration of the server
    * `interfaces` - Interface layout
      * `aggregation` - Interface aggregation status
      * `count` - Interface count
      * `type` - OVH Link Aggregation interface type (public┃vrack)
    * `name` - Mode name
  * `supported_modes` - Supported modes (DEPRECATED)
* `routing` - Routing details
  * `ipv4` - Ipv4 routing details
    * `gateway` - Server gateway
    * `ip` - Server main IP
    * `network` - Server network
  * `ipv6` - Ipv6 routing details
    * `gateway` - Server gateway
    * `ip` - Server main IP
    * `network` - Server network
* `switching` - Switching details
  * `name` - Switch name
* `traffic` - Traffic details
  * `input_quota_size` - Monthly input traffic quota allowed
    * `unit`
    * `value`
  * `input_quota_used` - Monthly input traffic consumed this month
    * `unit`
    * `value`
  * `is_throttled` - Whether bandwidth is throttleted for being over quota
  * `output_quota_size` - Monthly output traffic quota allowed
    * `unit`
    * `value`
  * `output_quota_used` - Monthly output traffic consumed this month
    * `unit`
    * `value`
  * `reset_quota_date` - Next reset quota date for traffic counter
* `vmac` - VMAC information for this dedicated server
  * `supported` - Whether server is compatible vmac
* `vrack` - vRack details
  * `bandwidth` - vrack bandwidth limitation
    * `unit`
    * `value`
  * `type` - Bandwidth offer type (included┃standard)