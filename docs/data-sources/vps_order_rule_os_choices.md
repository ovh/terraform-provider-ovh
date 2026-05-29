---
subcategory : "VPS"
---

# ovh_vps_order_rule_os_choices (Data Source)

Use this data source to discover the operating system choices for a given datacenter / OS family before ordering a VPS.

## Example Usage

```terraform
data "ovh_vps_order_rule_os_choices" "linux_gra" {
  datacenter = "GRA"
  os         = "linux"
}
```

## Argument Reference

* `datacenter` - (Required) Datacenter code (e.g. `GRA`, `BHS`, `SBG`).
* `os` - (Required) OS family (e.g. `linux`, `windows`).

## Attributes Reference

* `choices` - List of OS choices:
  * `name` - OS name.
  * `status` - Availability status.
