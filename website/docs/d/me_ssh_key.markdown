---
layout: "ovh"
page_title: "OVH: me_ssh_key"
sidebar_current: "docs-ovh-datasource-ssh-key-x"
description: |-
  Get information & status of an SSH key.
---

# ovh_me_ssh_key

Use this data source to retrieve information about an SSH key.

## Example Usage

```hcl
data "ovh_me_ssh_key" "mykey" {
   key_name = "mykey"
}
```

## Argument Reference

* `key_name` - (Required) The name of the SSH key.

## Attributes Reference

* `key_name` - See Argument Reference above.
* `key` - The content of the public key.
E.g.: "ssh-ed25519 AAAAC3..."
* `default` - True when this public SSH key is used for rescue mode and reinstallations.
