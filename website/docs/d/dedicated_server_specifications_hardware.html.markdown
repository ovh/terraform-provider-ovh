---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_specifications_hardware (Data Source)

Use this data source to get the hardward information about a dedicated server associated with your OVHcloud Account.

## Example Usage

```hcl
data "ovh_dedicated_server_specifications_hardware" "spec" {
  service_name = "myserver"
}
```

## Argument Reference

* `service_name` - (Required) The internal name of your dedicated server.

## Attributes Reference

The following attributes are exported:

* `boot_mode` - Server boot mode
* `cores_per_processor` - Number of cores per processor
* `default_hardware_raid_size` - Default hardware raid size for this server
  * `unit`
  * `value`
* `default_hardware_raid_type` - Default hardware raid type configured on this server
* `description` - Commercial name of this server
* `disk_groups` - Details about the groups of disks in the server
  * `default_hardware_raid_size` - Default hardware raid size for this disk group
    * `unit`
    * `value`
  * `default_hardware_raid_type` - Default hardware raid type for this disk group
  * `description` - Human readable description of this disk group
  * `disk_group_id` - Identifier of this disk group
  * `disk_size` - Disk capacity
    * `unit`
    * `value`
  * `disk_type` - Type of the disk (SSD, SATA, SAS, ...)
  * `number_of_disks` - Number of disks in this group
  * `raid_controller` - Raid controller, if any, managing this group of disks
* `expansion_cards` - Details about the server's expansion cards
  * `description` - Expansion card description
  * `type` - Expansion card type enum
* `form_factor` - Server form factor
* `memory_size` - RAM capacity
  * `unit`
  * `value`
* `motherboard` - Server motherboard
* `number_of_processors` - Number of processors in this dedicated server
* `processor_architecture` - Processor architecture bit
* `processor_name` - Processor name
* `threads_per_processor` - Number of threads per processor
* `usb_keys` - Capacity of the USB keys installed on your server, if any
  * `unit`
  * `value`