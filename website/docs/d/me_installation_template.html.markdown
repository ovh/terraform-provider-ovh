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

* `template_name`: This template name

## Attributes Reference

The following attributes are exported:

* `bit_format`: This template bit format (32 or 64).
* `category`: Category of this template (informative only). (basic, customer, hosting, other, readyToUse, virtualisation).
* `customization`: 
  * `custom_hostname`: Set up the server using the provided hostname instead of the default hostname.
  * `post_installation_script_link`: Indicate the URL where your postinstall customisation script is located.
  * `post_installation_script_return`: indicate the string returned by your postinstall customisation script on successful execution. Advice: your script should return a unique validation string in case of succes. A good example is 'loh1Xee7eo OK OK OK UGh8Ang1Gu'.
    * `ssh_key_name`: Name of the ssh key that should be installed. Password login will be disabled. Deprecated, will be removed in next release, use userMetada instead.
* `description`: information about this template.
* `distribution`: the distribution this template is based on.
* `family`: this template family type (bsd,linux,solaris,windows).
* `hard_raid_configuration`: This distribution supports hardware raid configuration through the OVHcloud API.
* `filesystems`: Filesystems available (btrfs,ext3,ext4,ntfs,reiserfs,swap,ufs,xfs,zfs).
* `partition_scheme`: 
  * `name`: name of this partitioning scheme.
  * `priority`: on a reinstall, if a partitioning scheme is not specified, the one with the higher priority will be used by default, among all the compatible partitioning schemes (given the underlying hardware specifications).
  * `hardware_raid`: 
     * `name`: Hardware RAID name.
     * `disks`: Disk List. Syntax is cX:dY for disks and [cX:dY,cX:dY] for groups. With X and Y resp. the controller id and the disk id.
     * `mode`: RAID mode (raid0, raid1, raid10, raid5, raid50, raid6, raid60).
     * `step`: Specifies the creation order of the hardware RAID.
  * `partition`:
     * `filesystem`: Partition filesystem.
     * `mountpoint`: partition mount point.
     * `raid`: raid partition type.
     * `size`: size of partition in MB, 0 => rest of the space.
     * `order`: step or order. specifies the creation order of the partition on the disk
     * `type`: partition type.
     * `volume_name`: The volume name needed for proxmox distribution