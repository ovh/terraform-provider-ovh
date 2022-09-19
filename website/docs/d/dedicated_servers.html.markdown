---
layout: "ovh"
page_title: "OVH: dedicated_servers"
sidebar_current: "docs-ovh-datasource-dedicated-servers"
description: |-
  Get the list of dedicated servers associated with your OVHcloud Account.
---

# ovh_dedicated_servers (Data Source)

Use this data source to get the list of dedicated servers associated with your OVHcloud Account.

## Example Usage

```hcl
data "ovh_dedicated_servers" "servers" {}
```

## Argument Reference

This datasource takes no argument.

## Attributes Reference

The following attributes are exported:

* `result` - The list of dedicated servers IDs associated with your OVHcloud Account.
