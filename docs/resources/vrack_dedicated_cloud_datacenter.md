---
subcategory : "vRack"
---

# ovh_vrack_dedicated_cloud_datacenter

Move a Dedicated Cloud Datacenter to a vrack.

## Example Usage

~> **WARNING** You have to import the resource first as it cannot be created, see Import section below.

```terraform
resource "ovh_vrack_dedicated_cloud_datacenter" "vrack_dedicated_cloud_datacenter" {
    service_name         = "<vRack service name>"
    datacenter           = "<Dedicated Cloud Datacenter service name>"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your vrack
* `datacenter` - (Required) Your Dedicated Cloud Datacenter.

## Attributes Reference

No additional attribute is exported.

## Import

A Datacenter will always be in a vRack, first import the resource, this will move the Dedicated Cloud Datacenter to the vRack target.

```bash
$ terraform import ovh_vrack_dedicated_cloud_datacenter.vrack-dedicatedCloudDatacenter "<vRack service name>/<Dedicated Cloud Datacener service name>/<vRack target service name>"
```