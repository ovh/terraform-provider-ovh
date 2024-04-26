---
subcategory : "Dedicated Server"
---

# ovh_me_installation_template

Use this resource to create a custom installation template available for dedicated servers.

## Example Usage

```hcl
resource "ovh_me_installation_template" "mytemplate" {
  base_template_name = "debian12_64"
  template_name      = "mytemplate"
  custom_hostname    = "mytest"
}
```

## Argument Reference

* `base_template_name`: (Required) The name of an existing installation template, choose one among the list given by `ovh_dedicated_installation_templates` datasource.
* `bit_format`: This template bit format (32 or 64).
* `category`: Category of this template (informative only). (basic, customer, hosting, other, readyToUse, virtualisation).
* `customization`:
  * `custom_hostname`: Set up the server using the provided hostname instead of the default hostname.
  * `post_installation_script_link`: Indicate the URL where your postinstall customisation script is located.
  * `post_installation_script_return`: indicate the string returned by your postinstall customisation script on successful execution. Advice: your script should return a unique validation string in case of succes. A good example is 'loh1Xee7eo OK OK OK UGh8Ang1Gu'.
* `description`: information about this template.
* `distribution`: the distribution this template is based on.
* `family`: this template family type (bsd,linux,solaris,windows).
* `filesystems`: Filesystems available (btrfs,ext3,ext4,ntfs,reiserfs,swap,ufs,xfs,zfs).
* `hard_raid_configuration`: This distribution supports hardware raid configuration through the OVHcloud API. Deprecated, will be removed in next release.
* `remove_default_partition_schemes`: (Required) Remove default partition schemes at creation.
* `template_name`: (Required)  This template name.

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id`: This template name.

## Import

Custom installation template available for dedicated servers can be imported using the `base_template_name`, `template_name` of the cluster, separated by "/" E.g.,

```bash
$ terraform import ovh_me_installation_template.mytemplate base_template_name/template_name
```
