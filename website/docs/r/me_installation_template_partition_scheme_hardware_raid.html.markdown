---
subcategory : "Dedicated Server"
---

# ovh_me_installation_template_partition_scheme_hardware_raid

Use this resource to create a hardware raid group in the partition scheme of a custom installation template available for dedicated servers.

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

resource "ovh_me_installation_template_partition_scheme_hardware_raid" "group1" {
  template_name = ovh_me_installation_template_partition_scheme.scheme.template_name
  scheme_name   = ovh_me_installation_template_partition_scheme.scheme.name
  name          = "group1"
  disks         = ["[c1:d1,c1:d2,c1:d3]", "[c1:d10,c1:d20,c1:d30]"]
  mode          = "raid50"
  step          = 1
}
```

## Argument Reference

* `disks`: Disk List. Syntax is cX:dY for disks and [cX:dY,cX:dY] for groups. With X and Y resp. the controller id and the disk id.
* `mode`: RAID mode (raid0, raid1, raid10, raid5, raid50, raid6, raid60).
* `name`: Hardware RAID name.
* `scheme_name`: (Required) The partition scheme name.
* `step`: Specifies the creation order of the hardware RAID.
* `template_name`: (Required) The template name of the partition scheme.


## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id`: a fake id associated with this partition scheme hardware raid group formatted as follow: `template_name/scheme_name/name`

## Import

The resource can be imported using the `template_name`, `scheme_name`, `name` of the cluster, separated by "/" E.g.,

```bash
$ terraform import ovh_me_installation_template_partition_scheme_hardware_raid.group1 template_name/scheme_name/name
```
