---
subcategory : "vRack"
---

# ovh_vrack_dedicated_cloud_datacenter

Move a Dedicated Cloud Datacenter to a vrack.

## Example Usage

```terraform
resource "ovh_vrack_dedicated_cloud_datacenter" "vrack-dedicatedCloudDatacenter" {
    service_name         = "<vRack service name>"
    datacenter           = "<Dedicated Cloud Datacenter service name>"
	target_service_name  = "<vRack target service name>"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your vrack
* `datacenter` - (Required) Your Dedicated Cloud Datacenter.
* `target_service_name` - (Required) The internal name of the vrack where you want to move your datacenter.

## Attributes Reference

No additional attribute is exported.