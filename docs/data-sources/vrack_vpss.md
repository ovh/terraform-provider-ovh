---
subcategory : "vRack"
---

# ovh_vrack_vpss (Data Source)

Use this data source to get the list of VPSes attached to a vRack.

## Example Usage

```terraform
data "ovh_vrack_vpss" "vrackvpss" {
  service_name = "pn-xxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service name of the vRack.

## Attributes Reference

The following attributes are exported:

* `result` - The list of VPS service names attached to the vRack.
