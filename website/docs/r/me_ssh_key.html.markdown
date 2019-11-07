---
layout: "ovh"
page_title: "OVH: ovh_me_ssh_key"
sidebar_current: "docs-ovh-resource-me-ssh-key"
description: |-
  Creates an SSH Key.
---

# ovh_me_ssh_key

Creates an SSH Key.

## Example Usage

```hcl
resource "ovh_me_ssh_key" "mykey" {
  key_name = "mykey"
  key      = "ssh-ed25519 AAAAC3..."
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Required) The friendly name of this SSH key.

* `key` - (Required) The content of the public key in the form "ssh-algo content", e.g. "ssh-ed25519 AAAAC3...".

* `default` - True when this public SSH key is used for rescue mode and reinstallations.

## Attributes Reference

The following attributes are exported:

* `key_name` - See Argument Reference above.
* `key` - See Argument Reference above.
* `default` - See Argument Reference above.
