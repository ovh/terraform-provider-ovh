---
layout: "ovh"
page_title: "OVH: cloud_project_failover_ip_attach"
sidebar_current: "docs-ovh-resource-cloud-project-failover-ip-attach-x"
description: |-
  Attaches failover ip addresses to cloud instances
---

# ovh_cloud_project_failover_ip_attach

Attaches a failover IP address to a compute instance

## Example Usage

```hcl
resource "ovh_cloud_project_failover_ip_attach" "myfailoverip" {
 service_name = "XXXXXX"
 ip = "XXXXXX"
 routed_to = "XXXXXX"
}
```

## Argument Reference


* `service_name` - The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `ip` - The failover ip address to attach
* `routed_to` - The GUID of an instance to which the failover IP address is be attached

## Attributes Reference

* `block` - The IP block
* `continentCode` - The Ip continent
* `geoloc` - The Ip location
* `id` - The Ip id
* `ip` - The Ip Address
* `progress` - Current operation progress in percent
* `routedTo` - Instance where ip is routed to
* `status` - Ip status, can be `ok` or `operationPending`
* `subType` - IP sub type, can be `cloud` or `ovh`