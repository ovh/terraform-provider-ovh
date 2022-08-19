---
layout: "ovh"
page_title: "OVH: dedicated_server_boots"
sidebar_current: "docs-ovh-datasource-dedicated-server-boots"
description: |-
  Get the list of compatible netboots for a dedicated server associated with your OVH Account.
---

# ovh_dedicated_server_boots (Data Source)

Use this data source to get the list of compatible netboots for a dedicated server associated with your OVH Account.

## Example Usage

```hcl
data "ovh_dedicated_server_boots" "netboots" {
  service_name = "myserver"
  boot_type    = "harddisk"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your dedicated server.

* `boot_type` - (Optional) Filter the value of bootType property (harddisk, rescue, ipxeCustomerScript, internal, network)

## Attributes Reference

The following attributes are exported:

* `result` - The list of dedicated server netboots.
