---
layout: "ovh"
page_title: "OVH: dedicated_installation_templates"
sidebar_current: "docs-ovh-datasource-dedicated-installation-templates"
description: |-
  Get the list of installation templates available for dedicated servers.
---

# ovh_dedicated_installation_templates

Use this data source to get the list of installation templates available for dedicated servers.

## Example Usage

```hcl
data "ovh_dedicated_installation_templates" "templates" {}
```

## Argument Reference

This datasource takes no argument.

## Attributes Reference

The following attributes are exported:

* `result` - The list of installation templates IDs available for dedicated servers.
