---
layout: "ovh"
page_title: "OVH: ovh_me_installation_template"
sidebar_current: "docs-ovh-resource-me-installation-template-x"
description: |-
  Creates a custom installation template available for dedicated servers.
---

# ovh_me_installation_template

Use this resource to create a custom installation template available for dedicated servers.

## Example Usage

```hcl
resource "ovh_me_installation_template" "mytemplate" {
  base_template_name = "centos7_64"
  template_name      = "mytemplate"
  default_language   = "fr"
}
```

## Argument Reference

* `available_languages`: List of all language available for this template.
* `base_template_name`: (Required) OVH template name yours will be based on, choose one among the list given by compatibleTemplates function.
* `beta`: This distribution is new and, although tested and functional, may still display odd behaviour.
* `bit_format`: This template bit format (32 or 64).
* `category`: Category of this template (informative only). (basic, customer, hosting, other, readyToUse, virtualisation).
* `customization`:
  * `change_log`: Template change log details.
  * `custom_hostname`: Set up the server using the provided hostname instead of the default hostname.
  * `post_installation_script_link`: Indicate the URL where your postinstall customisation script is located.
  * `post_installation_script_return`: indicate the string returned by your postinstall customisation script on successful execution. Advice: your script should return a unique validation string in case of succes. A good example is 'loh1Xee7eo OK OK OK UGh8Ang1Gu'.
  * `rating`: Rating.
  * `ssh_key_name`: Name of the ssh key that should be installed. Password login will be disabled.
  * `use_distribution_kernel`: Use the distribution's native kernel instead of the recommended OV
* `default_language`: (Required)  The default language of this template.
* `deprecated`: is this distribution deprecated.
* `description`: information about this template.
* `distribution`: the distribution this template is based on.
* `family`: this template family type (bsd,linux,solaris,windows).
* `filesystems`: Filesystems available (btrfs,ext3,ext4,ntfs,reiserfs,swap,ufs,xfs,zfs).
* `hard_raid_configuration`: This distribution supports hardware raid configuration through the OVH API.
* `last_modification`: Date of last modification of the base image.
* `remove_default_partition_schemes`: (Required) Remove default partition schemes at creation.
* `supports_distribution_kernel`: This distribution supports installation using the distribution's native kernel instead of the recommended OVH kernel.
* `supports_gpt_label`: This distribution supports the GUID Partition Table (GPT), providing up to 128 partitions that can have more than 2TB.
* `supports_rtm`: This distribution supports RTM software.
* `supports_sql_server`: This distribution supports the microsoft SQL server.
* `supports_uefi`: This distribution supports UEFI setup (no,only,yes).
* `template_name`: (Required)  This template name.


## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id`: This template name.

## Import

Use the following id format to import the resource : `base_template_name/template_name`
