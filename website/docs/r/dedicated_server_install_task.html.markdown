---
layout: "ovh"
page_title: "OVH: dedicated_server_install_task"
sidebar_current: "docs-ovh-resource-dedicated-server-install-task"
description: |-
  Install your Dedicated Server
---

# ovh_dedicated_server_install_task

Install your Dedicated Server.

~> __WARNING__ After some delay, if the task is marked as `done`, the Provider
may purge it. To avoid raising errors when terraform refreshes its plan, 
404 errors are ignored on Resource Read, thus some information may be lost
after a while.

## Example Usage

```hcl
data "ovh_dedicated_server_boots" "rescue" {
  service_name = "ns00000.ip-1-2-3.eu"
  boot_type    = "rescue"
}

resource "ovh_me_ssh_key" "key" {
  key_name = "mykey"
  key      = "ssh-ed25519 AAAAC3..."
}

resource "ovh_me_installation_template" "debian" {
  base_template_name = "debian10_64"
  template_name      = "mydebian10"
  default_language   = "en"

  customization {
    change_log      = "v1"
    ssh_key_name    = ovh_me_ssh_key.key.key_name
  }
}

resource "ovh_dedicated_server_install_task" "server_install" {
  service_name      = "ns00000.ip-1-2-3.eu"
  template_name     = ovh_me_installation_template.debian.template_name
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]

  details {
      custom_hostname = "mytest"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service_name of your dedicated server.
* `partition_scheme_name` - Partition scheme name.
* `template_name` - (Required) Template name.
* `bootid_on_destroy` - If set, reboot the server on the specified boot id during destroy phase.
* `details` - see `details` block below.

The `details` block supports:

* `change_log` - Template change log details.
* `custom_hostname` - Set up the server using the provided hostname instead of the default hostname.
* `disk_group_id` - Disk group id.
* `install_rtm` - set to true to install RTM.
* `install_sql_server` - set to true to install sql server (Windows template only).
* `language` - language.
* `no_raid` - set to true to disable RAID.
* `post_installation_script_link` - Indicate the URL where your postinstall customisation script is located.
* `post_installation_script_return` - Indicate the string returned by your postinstall customisation script on successful execution. Advice: your script should return a unique validation string in case of succes. A good example is 'loh1Xee7eo OK OK OK UGh8Ang1Gu'.
* `reset_hw_raid` - set to true to make a hardware raid reset.
* `soft_raid_devices` - soft raid devices.
* `ssh_key_name` - Name of the ssh key that should be installed. Password login will be disabled.
* `use_spla` - set to true to use SPLA.
* `use_distrib_kernel` - Use the distribution's native kernel instead of the recommended OVHcloud Kernel.

## Attributes Reference

The following attributes are exported:

* `id` - The task id
* `comment` - Details of this task. (should be `Install asked`)
* `done_date` - Completion date in RFC3339 format.
* `function` - Function name (should be `hardInstall`).
* `last_update` - Last update in RFC3339 format.
* `start_date` - Task creation date in RFC3339 format.
* `status` - Task status (should be `done`)
