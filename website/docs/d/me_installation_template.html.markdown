---
subcategory : "Dedicated Server"
---

# ovh_me_installation_template (Data Source)

Use this data source to get a custom installation template available for dedicated servers.

## Example Usage

```hcl
data "ovh_me_installation_template" "mytemplate" {
  template_name = "mytemplate"
}
```

## Argument Reference

* `template_name`: Template name.

## Attributes Reference

The following attributes are exported:

* `bit_format`: Template bit format (32 or 64).
* `category`: Category of this template (informative only).
* `customization`: 
  * `custom_hostname`: Set up the server using the provided hostname instead of the default hostname.
  * `post_installation_script_link`: Indicate the URL where your postinstall customisation script is located.
  * `post_installation_script_return`: Indicate the string returned by your postinstall customisation script on successful execution. Advice: your script should return a unique validation string in case of succes. A good example is 'loh1Xee7eo OK OK OK UGh8Ang1Gu'.
* `description`: Information about this template.
* `distribution`: Distribution this template is based on.
* `end_of_install` - End of install date of the template.
* `family`: Template family type (bsd,linux,solaris,windows).
* `filesystems`: Filesystems available.
* `hard_raid_configuration`: Distribution supports hardware raid configuration through the OVHcloud API.
* `lvm_ready` - Whether this template supports LVM.
* `no_partitioning` - Partitioning customization is not available for this OS template.
* `partition_scheme`: 
  * `name`: Name of this partitioning scheme.
  * `priority`: On a reinstall, if a partitioning scheme is not specified, the one with the higher priority will be used by default, among all the compatible partitioning schemes (given the underlying hardware specifications).
  * `hardware_raid`: 
     * `name`: Hardware RAID name.
     * `disks`: Disk List. Syntax is cX:dY for disks and [cX:dY,cX:dY] for groups. With X and Y resp. the controller id and the disk id.
     * `mode`: RAID mode (raid0, raid1, raid10, raid5, raid50, raid6, raid60).
     * `step`: Specifies the creation order of the hardware RAID.
  * `partition`:
     * `filesystem`: Partition filesystem.
     * `mountpoint`: Partition mount point.
     * `raid`: Raid partition type.
     * `size`: Size of partition in MB, 0 => rest of the space.
     * `order`: Step or order. Specifies the creation order of the partition on the disk.
     * `type`: Partition type.
     * `volume_name`: Volume name needed for proxmox distribution.
* `soft_raid_only_mirroring` - Template supports RAID0 and RAID1 on 2 disks.
* `subfamily` - Subfamily of the template.
