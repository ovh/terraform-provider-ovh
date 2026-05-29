---
subcategory : "VPS"
---

# ovh_vps_service_info (Data Source)

Use this data source to read the service information (lifecycle, renewal, contacts) of a VPS associated with your OVHcloud Account.

## Example Usage

```terraform
data "ovh_vps_service_info" "info" {
  service_name = "XXXXXX"
}
```

## Argument Reference

* `service_name` - (Required) The service name of your VPS.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `service_id` - The OVH service ID
* `status` - Service status
* `creation` - Service creation date (RFC 3339)
* `expiration` - Service expiration date (RFC 3339)
* `engaged_up_to` - Date until which the service is engaged (empty if none)
* `renewal_type` - One of `automaticForcedProduct`, `automaticV2012`, `automaticV2014`, `automaticV2016`, `manual`, `oneShot`, `option`
* `contact_admin` - The administrative contact nichandle
* `contact_billing` - The billing contact nichandle
* `contact_tech` - The technical contact nichandle
* `domain` - The associated domain
* `can_delete_at_expiration` - Whether the service may be marked for deletion at expiration
* `possible_renew_period` - The list of renewal periods (months) accepted by the API
* `renew_automatic` - Whether automatic renewal is enabled
* `renew_delete_at_expiration` - Whether the service is set to be deleted at expiration
* `renew_forced` - Whether the renewal is forced
* `renew_manual_payment` - Whether the renewal requires manual payment
* `renew_period` - The current renewal period in months
