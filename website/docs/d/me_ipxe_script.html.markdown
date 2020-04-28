---
layout: "ovh"
page_title: "OVH: me_ipxe_script"
sidebar_current: "docs-ovh-datasource-ipxe-script-x"
description: |-
  Get information & status of an IPXE Script.
---

# ovh_me_ipxe_script

Use this data source to retrieve information about an IPXE Script.

## Example Usage

```hcl
data "ovh_me_ipxe_script" "script" {
   name = "myscript"
}
```

## Argument Reference

* `name` - (Required) The name of the IPXE Script.

## Attributes Reference

* `name` - See Argument Reference above.
* `script` - The content of the script.
