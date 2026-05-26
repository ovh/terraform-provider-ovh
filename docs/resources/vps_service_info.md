---
subcategory : "VPS"
---

# ovh_vps_service_info (Resource)

Manage the auto-renew configuration of an existing OVHcloud VPS. This resource does **not** own the underlying VPS (the `ovh_vps` resource does). It only manipulates the writable `renew` fields exposed by `/vps/{serviceName}/serviceInfos`.

~> **NOTE** Deleting this resource will not delete the VPS. It will restore the default renew configuration (`renew_automatic = false`, `renew_delete_at_expiration = false`, `renew_forced = false`) and remove the resource from Terraform state.

## Example Usage

```terraform
resource "ovh_vps_service_info" "info" {
  service_name               = "vpsXXXXX.ovh.net"
  renew_automatic            = true
  renew_period               = 1
  renew_delete_at_expiration = false
  renew_forced               = false
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) The service name of your VPS.
* `renew_automatic` - (Required, Bool) Whether automatic renewal is enabled.
* `renew_delete_at_expiration` - (Optional, Bool) Whether the service should be deleted at expiration. Defaults to `false`.
* `renew_forced` - (Optional, Bool) Whether renewal is forced. Defaults to `false`.
* `renew_manual_payment` - (Optional, Bool) Whether renewal requires manual payment.
* `renew_period` - (Optional, Int) Renewal period in months. Must be one of the values returned in `possible_renew_period`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `service_id` - The OVH service ID.
* `status` - Service status.
* `creation` - Service creation date.
* `expiration` - Service expiration date.
* `engaged_up_to` - Date until which the service is engaged.
* `renewal_type` - One of `automaticForcedProduct`, `automaticV2012`, `automaticV2014`, `automaticV2016`, `manual`, `oneShot`, `option`.
* `possible_renew_period` - The list of accepted renewal periods (months).

## Import

The resource can be imported using the VPS service name:

```bash
terraform import ovh_vps_service_info.info vpsXXXXX.ovh.net
```
