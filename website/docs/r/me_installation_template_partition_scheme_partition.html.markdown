---
subcategory : "Dedicated Server"
---

# ovh_me_installation_template_partition_scheme_partition

Use this resource to create a partition in the partition scheme of a custom installation template available for dedicated servers.

## Example Usage

```hcl
resource "ovh_me_installation_template" "mytemplate" {
  base_template_name = "debian12_64"
  template_name      = "mytemplate"
}

resource "ovh_me_installation_template_partition_scheme" "scheme" {
  template_name = ovh_me_installation_template.mytemplate.template_name
  name          = "myscheme"
  priority      = 1
}

resource "ovh_me_installation_template_partition_scheme_partition" "root" {
  template_name = ovh_me_installation_template_partition_scheme.scheme.template_name
  scheme_name   = ovh_me_installation_template_partition_scheme.scheme.name
  mountpoint    = "/"
  filesystem    = "ext4"
  size          = "400"
  order         = 1
  type          = "primary"
}
```

## Argument Reference

* `filesystem`: Partition filesystem. Enum with possibles values:
	- btrfs
	- ext3
	- ext4
	- ntfs
	- reiserfs
	- swap
	- ufs
	- xfs
	- zfs
* `mountpoint`: (Required) partition mount point.
* `order`: step or order. specifies the creation order of the partition on the disk
* `raid`: raid partition type. Enum with possible values: 
  - raid0
  - raid1
  - raid10
  - raid5
  - raid6
* `scheme_name`: (Required) The partition scheme name.
* `size`: size of partition in MB, 0 => rest of the space.
* `template_name`: (Required) The template name of the partition scheme.
* `type`: partition type. Enum with possible values:
	- lv
	- primary
	- logical
* `volume_name`: The volume name needed for proxmox distribution


## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id`: a fake id associated with this partition scheme partition formatted as follow: `template_name/scheme_name/mountpoint`

## Import

The resource can be imported using the `template_name`, `scheme_name`, `mountpoint` of the cluster, separated by "/" E.g.,

```bash
$ terraform import ovh_me_installation_template_partition_scheme_partition.root template_name/scheme_name/mountpoint
```
