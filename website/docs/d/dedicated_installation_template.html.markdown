---
subcategory : "Dedicated Server"
---

# ovh_dedicated_installation_template (Data Source)

Use this data source to retrieve informations about a specific ovh dedicated server installation template.

## Example Usage

```hcl
data "ovh_dedicated_installation_template" "ovhtemplate" {
  template_name = "debian12_64"
}

output "template" {
  value = data.ovh_dedicated_installation_template.ovhtemplate
}
```

## Argument Reference

* `template_name` - (Required) The name of the template

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the installation template
* `bit_format`: This template bit format (32 or 64)
* `category` - Category of this template (informative only).
* `description` - Information about this template.
* `distribution` - The distribution this template is based on.
* `end_of_install` - The end of install date of the template
* `family` - This template family type.
* `filesystems` - Filesystems available.
* `hardware_raid_configuration` - This distribution supports hardware raid configuration through the OVHcloud API.
* `inputs` - Represents the questions of the expected answers in the userMetadata field
* `license` - The license available for this template
* `lvm_ready` - This template supports LVM.
* `no_partitioning` - Partitioning customization is not available for this OS template.
* `project` - The project
* `soft_raid_only_mirroring` - The template supports RAID0 and RAID1 on 2 disks.
* `subfamily` - The sub family of the template
