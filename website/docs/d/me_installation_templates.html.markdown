---
layout: "ovh"
page_title: "OVH: me_installation_templates"
sidebar_current: "docs-ovh-datasource-me-installation-templates"
description: |-
  Get the list of custom installation templates available for dedicated servers.
---

# ovh_me_installation_templates (Data Source)

Use this data source to get the list of custom installation templates available for dedicated servers.

## Example Usage

```hcl
data "ovh_me_installation_templates" "templates" {}
```

## Argument Reference

This datasource takes no argument.

## Attributes Reference

The following attributes are exported:

* `result` - The list of custom installation templates IDs available for dedicated servers.
