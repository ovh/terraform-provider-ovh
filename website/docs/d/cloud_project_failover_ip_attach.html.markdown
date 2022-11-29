---
layout: "ovh"
page_title: "OVH: cloud_project_failover_ip_attach"
sidebar_current: "docs-ovh-datasource-cloud-project-failover-ip-attach-x"
description: |-
  Reads failover ip addresses for a service
---

# ovh_cloud_project_failover_ip_attach (Data Source)

Use this data source to get the details of a failover ip address of a service in a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_failover_ip_attach" "myfailoverip" {
  service_name = "XXXXXX"
  ip           = "XXXXXX"
}
```

## Argument Reference


* `service_name` - The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 
* `ip` - The failover ip address to query

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
