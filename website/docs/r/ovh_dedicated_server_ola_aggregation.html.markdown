---
layout: "ovh"
page_title: "OVH: ovh_dedicated_server_ola_aggregation"
sidebar_current: "docs-ovh-resource-ovh-dedicated-server-ola-aggregation"
description: |-
  Group interfaces into an aggregation
---

# ovh_dedicated_server_ola_aggregation

Group interfaces into an aggregation.

## Example Usage

```hcl
data "ovh_dedicated_server" "server" {
  service_name = "nsXXXXXX.ip-XX-XXX-XX.eu"
}

resource "ovh_dedicated_server_ola_aggregation" "server"{
  service_name = "nsXXXXXX.ip-XX-XXX-XX.eu"
  name = "bond0"
  virtual_network_interfaces = data.ovh_dedicated_server.server.vnis.*.uuid
}

output "aggregation" {
  value = ovh_dedicated_server_ola_aggregation.server
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service_name of your dedicated server.
* `name` - (Required) The resulting VirtualNetworkInterface name.
* `virtual_network_interfaces` - (Required) The list of interfaces to aggregate.

## Attributes Reference

The following attributes are exported:

* `id` - The task id
* `comment` - Details of this task. (should be `Install asked`)
* `done_date` - Completion date in RFC3339 format.
* `function` - Function name (should be `hardInstall`).
* `last_update` - Last update in RFC3339 format.
* `start_date` - Task creation date in RFC3339 format.
* `status` - Task status (should be `done`)
* `enabled_vrack_aggregation_vni` - vrack_aggregation VNI uuid.
* `enabled_public_vni` - public VNI uuid.
