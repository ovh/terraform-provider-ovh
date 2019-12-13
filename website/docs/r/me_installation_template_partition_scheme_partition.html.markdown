---
layout: "ovh"
page_title: "OVH: ovh_me_installation_template_partition_scheme_partition"
sidebar_current: "docs-ovh-resource-me-installation-template-partition-scheme-partition"
description: |-
  Creates a partition in the partition scheme of a custom installation template available for dedicated servers.
---

# ovh_me_installation_template_partition_scheme_partition

Use this resource to create a partition in the partition scheme of a custom installation template available for dedicated servers.

## Example Usage

```hcl
resource "ovh_me_installation_template" "mytemplate" {
  base_template_name = "centos7_64"
  template_name      = "mytemplate"
  default_language   = "fr"
}

resource "ovh_me_installation_template_partition_scheme" "scheme" {
  template_name      = ovh_me_installation_template.mytemplate.template_name
  name               = "myscheme"
  priority           = 1
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

* `filesystem`: Partition filesystem.
* `mountpoint`: (Required) partition mount point.
* `order`: step or order. specifies the creation order of the partition on the disk
* `raid`: raid partition type.
* `scheme_name`: (Required) The partition scheme name.
* `size`: size of partition in MB, 0 => rest of the space.
* `template_name`: (Required) The template name of the partition scheme.
* `type`: partition type.
* `volume_name`: The volume name needed for proxmox distribution


## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id`: a fake id associated with this partition scheme partition formatted as follow: `template_name/scheme_name/mountpoint`

## Import

Use the fake id format to import the resource : `template_name/scheme_name/mountpoint` (example: "mytemplate/myscheme//").
