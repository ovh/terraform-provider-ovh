---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_install_task

Install your Dedicated Server.

~> __WARNING__ After some delay, if the task is marked as `done`, the Provider
may purge it. To avoid raising errors when terraform refreshes its plan, 
404 errors are ignored on Resource Read, thus some information may be lost
after a while.

## Example Usage

Using a custom template based on an OVHCloud template
```hcl
data ovh_dedicated_server_boots "rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_type    = "rescue"
}

resource "ovh_me_installation_template" "debian" {
  base_template_name = "debian12_64"
  template_name      = "mydebian12"
  customization {
    post_installation_script_link = "http://test"
    post_installation_script_return = "ok"
  }
}

resource "ovh_dedicated_server_install_task" "server_install" {
  service_name      = "nsxxxxxxx.ip-xx-xx-xx.eu"
  template_name     = ovh_me_installation_template.debian.template_name
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
  details {
      custom_hostname = "mytest"
  }
  user_metadata {
    key = "sshKey"
    value = "ssh-ed25519 AAAAC3..."
  }
}
```

Using a BringYourOwnLinux (BYOLinux) template (with userMetadata)

```hcl
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data ovh_dedicated_server_boots "rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_type    = "rescue"
}

resource "ovh_dedicated_server_install_task" "server_install" {
  service_name      = data.ovh_dedicated_server.server.service_name
  template_name     = "byolinux_64"
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
  details {
      custom_hostname = "mytest"
  }
  user_metadata {
    key       = "imageURL"
    value = "https://myimage.qcow2"
  }
  user_metadata {  
    key = "imageType"
    value = "qcow2"
  }
  user_metadata {  
    key = "httpHeaders0Key"     
    value = "Authorization"
  }
  user_metadata {  
    key = "httpHeaders0Value"     
    value = "Basic bG9naW46xxxxxxx="
  }
  user_metadata {  
    key = "imageChecksumType"    
    value = "sha512"
  }
  user_metadata {  
    key = "imageCheckSum"
    value = "047122c9ff4d2a69512212104b06c678f5a9cdb22b75467353613ff87ccd03b57b38967e56d810e61366f9d22d6bd39ac0addf4e00a4c6445112a2416af8f225"
  }
  user_metadata {  
    key = "configDriveUserData" 
    value = "#cloud-config\nssh_authorized_keys:\n  - ${data.ovh_me_ssh_key.mykey.key}\n\nusers:\n  - name: patient0\n    sudo: ALL=(ALL) NOPASSWD:ALL\n    groups: users, sudo\n    shell: /bin/bash\n    lock_passwd: false\n    ssh_authorized_keys:\n      - ${data.ovh_me_ssh_key.mykey.key}\ndisable_root: false\npackages:\n  - vim\n  - tree\nfinal_message: The system is finally up, after $UPTIME seconds\n"
  }
}
```

Using a Microsoft Windows server OVHcloud template with a specific language

```hcl
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data ovh_dedicated_server_boots "rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_type    = "rescue"
}

resource "ovh_dedicated_server_install_task" "server_install" {
  service_name      = data.ovh_dedicated_server.server.service_name
  template_name     = "win2019-std_64"
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
  details {
    custom_hostname = "mytest"
  }
  user_metadata {
    key  = "language"
    value ="fr-fr"
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
* `user_metadata` - see `user_metadata` block below.


The `details` block supports:

* `custom_hostname` - Set up the server using the provided hostname instead of the default hostname.
* `disk_group_id` - Disk group id.
* `language` - Deprecated, will be removed in next release.
* `no_raid` - Set to true to disable RAID.
* `soft_raid_devices` - soft raid devices.
* `use_spla` - Deprecated, will be removed in next release.

The `user_metadata` block supports many arguments, here is a non-exhaustive list depending on the OS:

-[see OS questions](https://help.ovhcloud.com/csm/en-dedicated-servers-api-os-installation?id=kb_article_view&sysparm_article=KB0061951#os-questions)

-[see api](https://eu.api.ovh.com/console-preview/?section=%2Fdedicated%2FinstallationTemplate&branch=v1#get-/dedicated/installationTemplate/-templateName-) 

-[see documentation](https://help.ovhcloud.com/csm/en-ie-dedicated-servers-api-os-installation?id=kb_article_view&sysparm_article=KB0061950#create-an-os-installation-task) to get more information


## Attributes Reference

The following attributes are exported:

* `id` - The task id
* `comment` - Details of this task. (should be `Install asked`)
* `done_date` - Completion date in RFC3339 format.
* `function` - Function name (should be `hardInstall`).
* `last_update` - Last update in RFC3339 format.
* `start_date` - Task creation date in RFC3339 format.
* `status` - Task status (should be `done`)

## Import

Installation task can be imported using the `service_name` (`nsXXXX.ip...`) of the baremetal server, the `template_name` used  and ths `task_id`, separated by "/" E.g.,

```bash
$ terraform import ovh_dedicated_server_install_task nsXXXX.ipXXXX/template_name/12345
```
