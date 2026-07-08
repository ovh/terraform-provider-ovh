---
subcategory : "Public IPs"
---

# ovh_cloud_project_floatingips

~> **NOTE** We recommend using the `ovh_cloud_floating_ips` data source instead. Floating IPs can now also be managed directly with the `ovh_cloud_floating_ip` resource.

Use this data source to get the floating IPs of a public cloud project.

## Example Usage

To get information of floating IPs:

```terraform
data "ovh_cloud_project_floatingips" "ips" {
  service_name = "YYYY"
  region = "XXXX"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project
* `region_name` - (Required) Public cloud region name

## Attributes Reference

The following attributes are exported:

* `associated_entity` - Associated entity with the floating IP
  * `ip` - IP of the port
  * `id` - ID of the port
  * `gateway_id` - ID of the gateway
  * `type` - Type of the port (dhcp┃instance┃loadbalancer┃routerInterface┃unknown)
* `id` - ID of the floating IP
* `ip` - Value of the floating IP
* `network_id` - ID of the network
* `region_name` - Floating IP region
* `status` - Status of the floating IP (active┃down┃error)
