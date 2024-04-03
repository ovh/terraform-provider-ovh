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

resource "ovh_me_ssh_key" "key" {
  key_name = "mykey"
  key      = "ssh-ed25519 AAAAC3..."
}

resource "ovh_me_installation_template" "debian" {
  base_template_name = "debian11_64"
  template_name      = "mydebian11"
  default_language   = "en"

  customization {
    ssh_key_name    = ovh_me_ssh_key.key.key_name
  }
}

resource "ovh_dedicated_server_install_task" "server_install" {
  service_name      = "nsxxxxxxx.ip-xx-xx-xx.eu"
  template_name     = ovh_me_installation_template.debian.template_name
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]

  details {
      custom_hostname = "mytest"
  }
}
```
Using a BringYourOwnLinux (BYOL) template (with userMetadata)
```hcl
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

resource "ovh_me_ssh_key" "key" {
  key_name = "mykey"
  key      = "ssh-ed25519 AAAAC3..."
}

data ovh_dedicated_server_boots "rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_type    = "rescue"
}

resource "ovh_me_installation_template" "mytemplate" {
  base_template_name = "byolinux_64"
  template_name      = "mybyol_test"
}

resource "ovh_dedicated_server_install_task" "server_install" {
  service_name      = data.ovh_dedicated_server.server.service_name
  template_name     = ovh_me_installation_template.mytemplate.template_name
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
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
    key = "language"     
    value = "en"
  }
  user_metadata {  
    key = "imageCheckSum"
    value = "047122c9ff4d2a69512212104b06c678f5a9cdb22b75467353613ff87ccd03b57b38967e56d810e61366f9d22d6bd39ac0addf4e00a4c6445112a2416af8f225"
  }
  user_metadata {  
    key = "configDriveUserData" 
    value = "#cloud-config\nssh_authorized_keys:\n  - ${data.ovh_me_ssh_key.mykey.key}\n\nusers:\n  - name: patient0\n    sudo: ALL=(ALL) NOPASSWD:ALL\n    groups: users, sudo\n    shell: /bin/bash\n    lock_passwd: false\n    ssh_authorized_keys:\n      - ${data.ovh_me_ssh_key.mykey.key}\ndisable_root: false\npackages:\n  - vim\n  - tree\nfinal_message: The system is finally up, after $UPTIME seconds\n"
  }
  details {
      custom_hostname = "mytest"
  }
}
```
Using a Microsoft Windows server OVHcloud template with a specific language
hcl```
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
  user_metadata {
    key  = "language"
    value ="fr-fr"
  }
 user_metadata {
    key  = "useSpla"
    value = "true"
 }
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
* `user_metadata` - see `user_metadata` block below.


The `details` block supports:

* `custom_hostname` - Set up the server using the provided hostname instead of the default hostname.
* `disk_group_id` - Disk group id.
* `install_sql_server` - Set to true to install sql server (Windows template only).
* `language` - Language.
* `no_raid` - Set to true to disable RAID.
* `post_installation_script_link` - Indicate the URL where your postinstall customisation script is located.
* `post_installation_script_return` - Indicate the string returned by your postinstall customisation script on successful execution. Advice: your script should return a unique validation string in case of succes. A good example is 'loh1Xee7eo OK OK OK UGh8Ang1Gu'.
* `soft_raid_devices` - soft raid devices.
* `ssh_key_name` - Name of the ssh key that should be installed. Password login will be disabled.
* `use_spla` - Set to true to use SPLA.

The `user_metadata` block supports :
(but is not limited to ! : [see documentation](https://help.ovhcloud.com/csm/world-documentation-bare-metal-cloud-dedicated-servers-managing-servers?id=kb_browse_cat&kb_id=203c4f65551974502d4c6e78b7421996&kb_category=97feff9459b0a510f078155c0c16be9b))

* `imageUrl` - Your Linux image URL.
* `imageType` - Your Linux image type (qcow2, raw).
* `imageCheckSum` - Your image's checksum.
* `imageCheckSumType` - Your image's checksum type.
* `httpHeadersNKey` - Your image's HTTP headers key (where N is a integer).
* `httpHeadersNValue` - Your image's HTTP headers value (where N is a integer).
* `configDriveUserData` - Your user config drive user data.
* `configDriveMetadata0Key` - Your user config drive user metadata key(where N is a integer).
* `configDriveMetadata0Value` - Your user config drive user metadata value (where N is a integer).
* `language` - Language.
* `useSpla` - Set to true to use SPLA.


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
