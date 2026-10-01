---
subcategory : "Web Cloud Private SQL"
---

# ovh_hosting_privatedatabase_webhosting_network

Allow or deny the OVHcloud webhosting network to connect to your private database.

A private database has this access enabled, so `enabled = true` matches the state of a brand new
instance and does nothing. Destroying the resource restores that default rather than leaving the
access closed.

## Example Usage

```terraform
resource "ovh_hosting_privatedatabase_webhosting_network" "network" {
  service_name = "XXXXXX"
  enabled      = false
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your private database.
* `enabled` - (Required) Allow the OVHcloud webhosting network to connect to the private database. Values can be `true` or `false`

## Attributes Reference

* `status` - Webhosting network status (`disabled`, `disabling`, `enabled` or `enabling`)

The id is set to the value of `service_name`.

## Import

The webhosting network access can be imported using the `service_name`. E.g.,

```
$ terraform import ovh_hosting_privatedatabase_webhosting_network.network service_name
```
