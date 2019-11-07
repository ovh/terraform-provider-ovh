---
layout: "ovh"
page_title: "OVH: me_ssh_keys"
sidebar_current: "docs-ovh-datasource-ssh-keys"
description: |-
  Get the list of the SSH keys of the account.
---

# ovh_me_ssh_keys

Use this data source to retrieve list of names of the account's SSH keys.

## Example Usage

```hcl
data "ovh_me_ssh_keys" "mykeys" {}
```

## Argument Reference

This datasource takes no arguments.

## Attributes Reference

* `names` - The list of the names of all the SSH keys.
