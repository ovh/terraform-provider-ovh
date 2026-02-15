---
subcategory : "VPS"
---

# ovh_vps (Data Source)

Use this data source to retrieve information about a vps associated with your OVHcloud Account.

## Example Usage

```terraform
data "ovh_vps" "server" {
  service_name = "XXXXXX"
}
```

## Argument Reference

* `service_name` - (Required) The service_name of your dedicated server.

## Attributes Reference

`id` is set with the service_name of the vps name (ex: "vps-123456.ovh.net")

In addition, the following attributes are exported:

* `urn` - The URN of the vps
* `cluster` - The OVHcloud cluster the vps is in
* `datacenter` - The datacenter in which the vps is located
  * `datacenter.longname` - The fullname of the datacenter (ex: "Strasbourg SBG1")
  * `datacenter.name` - The short name of the datacenter (ex: "sbg1)
* `displayname` - The displayed name in the OVHcloud web admin
* `ips` - The list of IPs addresses attached to the vps
* `keymap` - The keymap for the ip kvm, valid values "", "fr", "us"
* `memory` - The amount of memory in MB of the vps.
* `model` - A dict describing the type of vps.
* `model.name` - The model name (ex: model1)
* `model.offer` - The model human description (ex: "VPS 2016 SSD 1")
* `model.version` - The model version (ex: "2017v2")
* `netbootmode` - The source of the boot kernel
* `offertype` - The type of offer (ssd, cloud, classic)
* `slamonitoring` - A boolean to indicate if OVHcloud SLA monitoring is active.
* `state` - The state of the vps
* `type` - The type of server
* `vcore` - The number of vcore of the vps
* `zone` - The OVHcloud zone where the vps is
