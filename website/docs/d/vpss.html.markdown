---
layout: "ovh"
page_title: "OVH: vpss"
sidebar_current: "docs-ovh-datasource-vpss"
description: |-
  Get the list of VPS associated with your OVH Account.
---

# vpss (Data Source)

Use this data source to get the list of VPS associated with your OVH Account.

## Example Usage

```hcl
data "ovh_vpss" "servers" {}
```

## Argument Reference

This datasource takes no argument.

## Attributes Reference

The following attributes are exported:

* `result` - The list of VPS IDs associated with your OVH Account.
