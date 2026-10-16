---
subcategory : "Additional IP"
---

# ovh_cloud_project_failover_ip_attach

~> **NOTE** This resource relies on the Public Cloud API v1. Prefer the `networks` attribute of the new `ovh_cloud_instance` resource. See the [migration guide](../guides/public_cloud_resources_migration) to move an existing resource without recreating it.

Attaches a failover IP address to a compute instance

## Example Usage

```terraform
resource "ovh_cloud_project_failover_ip_attach" "failover_ip" {
  service_name = "XXXXXX"
  ip           = "XXXXXX"
  routed_to    = "XXXXXX"
}
```

## Argument Reference

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `ip` - The failover ip address to attach
* `routed_to` - The GUID of an instance to which the failover IP address is be attached

## Attributes Reference

* `block` - The IP block
* `continentCode` - The Ip continent
* `geoloc` - The Ip location
* `id` - The Ip id
* `ip` - The Ip Address
* `progress` - Current operation progress in percent
* `routedTo` - Instance where ip is routed to
* `status` - Ip status, can be `ok` or `operationPending`
* `subType` - IP sub type, can be `cloud` or `ovh`

## Import

A cloud project failover IP attachment can be imported using the `id` (failover IP block ID) E.g.,

```bash
$ terraform import ovh_cloud_project_failover_ip_attach.my_ip id
```
