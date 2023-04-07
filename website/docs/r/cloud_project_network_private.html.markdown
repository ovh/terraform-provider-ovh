---
subcategory : "Public Cloud Network"
---

# ovh_cloud_project_network_private

Creates a private network in a public cloud project.

## Example Usage

```hcl
resource "ovh_cloud_project_network_private" "net" {
  service_name = "XXXXXX"
  name         = "admin_network"
  regions      = ["GRA1", "BHS1"]
}
```

## Argument Reference

The following arguments are supported:


* `service_name` - (Required) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 

* `name` - (Required) The name of the network.

* `vlan_id` - a vlan id to associate with the network.
   Changing this value recreates the resource. Defaults to 0.

* `regions` - an array of valid OVHcloud public cloud region ID in which the network
   will be available. Ex.: "GRA1". Defaults to all public cloud regions.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the network
* `service_name` - See Argument Reference above.
* `name` - See Argument Reference above.
* `vlan_id` - See Argument Reference above.
* `regions` - See Argument Reference above.
* `regions_attributes` - A map representing information about the region.
* `regions_attributes/region` - The id of the region.
* `regions_attributes/status` - The status of the network in the region.
* `regions_attributes/openstackid` - The private network id in the region.
* `regions_status` - (Deprecated) A map representing the status of the network per region.
* `regions_status/region` - (Deprecated) The id of the region.
* `regions_status/status` - (Deprecated) The status of the network in the region.
* `status` - the status of the network. should be normally set to 'ACTIVE'.
* `type` - the type of the network. Either 'private' or 'public'. 

## Import

Private network in a public cloud project can be imported using the `service_name` and the `network_id`, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_network_private.mynet ookie9mee8Shaeghaeleeju7Xeghohv6e/pn-12345678
```
