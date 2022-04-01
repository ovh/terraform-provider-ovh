---
layout: "ovh"
page_title: "OVH: ovh_hosting_privatedatabase_whitelist"
sidebar_current: "docs-ovh-resource-hosting-privatedatabase-whitelist"
description: |-
  Create a new IP whitelist on your private cloud database instance.
---

# ovh_hosting_privatedatabase_whitelist

Create a new IP whitelist on your private cloud database instance.

## Example Usage

```hcl
resource "ovh_hosting_privatedatabase_whitelist" "ip" {
  service_name = "XXXXXX"
  ip           = "1.2.3.4"
  name         = "A name for your IP address"
  service      = true
  sftp         = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - The internal name of your private database.
* `ip` - (Required) The whitelisted IP in your instance.
* `name` - (Required) Custom name for your Whitelisted IP.
* `service` - (Required) Authorize this IP to access service port. Values can be `true` or `false`
* `sftp` - (Required) Authorize this IP to access SFTP port. Values can be `true` or `false`

## Attributes Reference

The id is set to the value of `service_name`/`ip`.

## Import

OVHcloud database whitelist can be imported using the `service_name` and the `ip`, separated by "/" E.g.,

```
$ terraform import ovh_hosting_privatedatabase_whitelist.ip service_name/ip
```