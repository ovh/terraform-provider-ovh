---
subcategory : "Dedicated Server"
---

# ovh_dedicated_installation_template (Data Source)

Use this data source to retrieve information about a specific OVH dedicated server installation template.

## Example Usage

```terraform
data "ovh_dedicated_installation_template" "ovh_template" {
  template_name = "debian12_64"
}

output "template" {
  value = data.ovh_dedicated_installation_template.ovh_template
}
```

## Argument Reference

* `template_name` - (Required) The name of the template.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the installation template.
* `bit_format`: Template bit format (32 or 64).
* `category` - Category of this template (informative only).
* `description` - Information about this template.
* `distribution` - Distribution this template is based on.
* `end_of_install` - End of install date of the template.
* `family` - Template family type.
* `filesystems` - Filesystems available.
* `hardware_raid_configuration` - Distribution supports hardware raid configuration through the OVHcloud API.
* `inputs` - Represents the questions of the expected answers in the userMetadata field.
* `license` - License available for this template.
* `lvm_ready` - Whether this template supports LVM.
* `no_partitioning` - Partitioning customization is not available for this OS template.
* `project` - Distribution project details.
* `soft_raid_only_mirroring` - Template supports RAID0 and RAID1 on 2 disks.
* `subfamily` - Subfamily of the template.
