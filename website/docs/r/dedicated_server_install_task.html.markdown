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

Using a custom template based on an ovh template
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
  default_language   = "en"
  customization      {
	  ssh_key_name = ovh_me_ssh_key.key.key_name
	}
}

resource "ovh_dedicated_server_install_task" "server_install" {
  service_name      = data.ovh_dedicated_server.server.service_name
  template_name     = ovh_me_installation_template.debian.template_name
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
  user_metadata     {
    image_url              = "https://myimage.qcow2"
    image_type             = "qcow2"
    http_headers_0_key     = "Authorization"
    http_headers_0_value   = "Basic bG9naW46xxxxxxx="
    image_checksum_type    = "sha512"
    language               = "en"
    image_checksum         = "047122c9ff4d2a69512212104b06c678f5a9cdb22b75467353613ff87ccd03b57b38967e56d810e61366f9d22d6bd39ac0addf4e00a4c6445112a2416af8f225"
    config_drive_user_data = "#cloud-config\nssh_authorized_keys:\n  - ${data.ovh_me_ssh_key.mykey.key}\n\nusers:\n  - name: patient0\n    sudo: ALL=(ALL) NOPASSWD:ALL\n    groups: users, sudo\n    shell: /bin/bash\n    lock_passwd: false\n    ssh_authorized_keys:\n      - ${data.ovh_me_ssh_key.mykey.key}\ndisable_root: false\npackages:\n  - vim\n  - tree\nfinal_message: The system is finally up, after $UPTIME seconds\n"
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

The `user_metadata` block supports:

* `image_url` - Your Linux image URL
* `image_type` - Your Linux image type (qcow2, raw)
* `image_checksum` - Your image's checksum
* `image_checksum_type` - Your image's checksum type
* `http_headers_0_key` - Your image's HTTP headers key (up to 5)
* `http_headers_0_value` - Your image's HTTP headers value (up to 5)
* `config_drive_user_data` - Your user config drive user data
* `language` - Language.
* `use_spla` - Set to true to use SPLA.


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
