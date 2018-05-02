---
layout: "ovh"
page_title: "OVH: ovh_vps"
sidebar_current: "docs-ovh-resource-vps"
description: |-
  Provides a OVH vps resource 
---

# ovh\_vps

Provides a OVH vps resource

## Example Usage

```hcl
# Add a vps resource
resoure ovh_vps test {
    type = vps_ssd_model1_2018v3
    displayname = "test"
    keymap = "us"
    slamonitoring: false
    netbootmode: "rescue"
}
```

## Argument Reference

The following attributes are supported:

* `type` - (Required) A valid ovh vps type.
* `displayname` - (Required) The displayed name in the ovh web admin
* `keymap` - The keymap for the ip kvm, valid values "", "fr", "us", default to ""
* `slamonitoring` - A boolean to activate OVH sla monitoring, default to `true`
* `netbootmode` - The source of the boot kernel, either `local` or `rescue`. Default to `local`
* `state` - The state of the vps ("running", "stopped", "terminating") 

## Attribute Reference

* `cluster` - The ovh cluster the vps is in
* `datacenter` - The datacenter in which the vps is located
  * `datacenter.longname` - The fullname of the datacenter (ex: "Strasbourg SBG1")
  * `datacenter.name` - The short name of the datacenter (ex: "sbg1)
* `displayname` - The displayed name in the ovh web admin
* `id` - The vps name (ex: "vps-123456.ovh.net")
* `ips` - The list of IPs addresses attached to the vps
* `keymap` - The keymap for the ip kvm, valid values "", "fr", "us"
* `memory` - The amount of memory in MB of the vps. 
* `model` - A dict describing the type of vps.
* `model.name` - The model name (ex: model1)
* `model.offer` - The model human description (ex: "VPS 2016 SSD 1")
* `model.version` - The model version (ex: "2017v2")
* `netbootmode` - The source of the boot kernel
* `offertype` - The type of offer (ssd, cloud, classic)
* `slamonitoring` - A boolean to indicate if OVH sla monitoring is active.
* `state` -  The state of the vps
* `type` - The type of server
* `vcore` - The number of vcore of the vps
* `zone` - The OVH zone where the vps is
