---
layout: "ovh"
page_title: "OVH: ovh_me_ipxe_script"
sidebar_current: "docs-ovh-resource-me-ipxe-script"
description: |-
  Creates an IPXE Script.
---

# ovh_me_ipxe_script

Creates an IPXE Script.

## Example Usage

```hcl
resource "ovh_me_ipxe_script" "script" {
  name   = "myscript"
  script = file("${path.module}/boot.ipxe")
}
```

## Argument Reference

The following arguments are supported:

* `description` - For documentation purpose only. This attribute is not passed to the OVHcloud API as it cannot be retrieved back. Instead a fake description ('$name auto description') is passed at creation time.

* `name` - (Required) The name of the IPXE Script.

* `script` - (Required) The content of the script.

## Attributes Reference

The following attributes are exported:

* `description` - See Argument Reference above.
* `name` - See Argument Reference above.
* `script` - See Argument Reference above.
