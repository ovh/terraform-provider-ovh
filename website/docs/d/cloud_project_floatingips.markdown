---
subcategory : "Cloud Project"
---

# ovh_cloud_project_floatingips
Use this data source to get the floating ips of a public cloud project.

## Example Usage

To get information of an instance:

```hcl
data "ovh_cloud_project_floatingips" "ips" {
  service_name = "YYYY"
  region = "XXXX"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used
* `region` - (Required) Instance region

## Attributes Reference

The following attributes are exported:

* `associated_entity` - Associated entity with the floating ip
  * `ip` - IP of the port
  * `id` - ID of the port
  * `gateway_id` - ID of the gateway
  * `type` - Type of the port (dhcp┃instance┃loadbalancer┃routerInterface┃unknown)
* `id` - ID of the floating ip
* `ip` - Value of the floating ip
* `network_id` - ID of the network
* `region` - Floating ip region
* `status` - Status of the floating ip (active┃down┃error)