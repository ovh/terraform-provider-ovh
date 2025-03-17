---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_reinstall_task

Install your Dedicated Server.

~> **WARNING** After some delay, if the task is marked as `done`, the Provider may purge it. To avoid raising errors when terraform refreshes its plan, 404 errors are ignored on Resource Read, thus some information may be lost after a while.

## Examples Usage

### Example 1 - Simple Linux Installation with password-based authentication

```terraform
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_reinstall" {
  service_name = data.ovh_dedicated_server.server.service_name
  os           = data.ovh_dedicated_installation_template.template.template_name
}
```

### Example 2 - Linux Installation with ssh key-based authentication and post-installation script

```terraform
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_reinstall" {
  service_name = data.ovh_dedicated_server.server.service_name
  os           = data.ovh_dedicated_installation_template.template.template_name
  customizations {
    hostname                 = "mon-tux"
    post_installation_script = "IyEvYmluL2Jhc2gKZWNobyAiY291Y291IHBvc3RJbnN0YWxsYXRpb25TY3JpcHQiID4gL29wdC9jb3Vjb3UKY2F0IC9ldGMvbWFjaGluZS1pZCAgPj4gL29wdC9jb3Vjb3UKZGF0ZSAiKyVZLSVtLSVkICVIOiVNOiVTIiAtLXV0YyA+PiAvb3B0L2NvdWNvdQo="
    ssh_key                  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC9xPpdqP3sx2H+gcBm65tJEaUbuifQ1uGkgrWtNY0PRKNNPdy+3yoVOtxk6Vjo4YZ0EU/JhmQfnrK7X7Q5vhqYxmozi0LiTRt0BxgqHJ+4hWTWMIOgr+C2jLx7ZsCReRk+fy5AHr6h0PHQEuXVLXeUy/TDyuY2JPtUZ5jcqvLYgQ== my-nuclear-power-plant"
  }
}
```

Even though the post-installation script could be sent to Terraform provider and API directly in clear text by escaping special characters, it is recommended to send a base64-encoded script to the API.

You can use the following UNIX/Linux command to encode your script:

```bash
cat my-script.sh | base64 -w0
```

Here is the clear-text post-installation bash script from the example above:

```bash
#!/bin/bash
echo "coucou postInstallationScript" > /opt/coucou
cat /etc/machine-id  >> /opt/coucou
date "+%Y-%m-%d %H:%M:%S" --utc >> /opt/coucou
```

### Example 3 - Windows Installation with French display language and PowerShell post-installation script

```terraform
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_install" {
  service_name = data.ovh_dedicated_server.server.service_name
  os           = "win2022core-std_64"
  customizations {
    hostname                 = "ma-fenetre"
    language                 = "fr-fr"
    post_installation_script = "ImNvdWNvdSBwb3N0SW5zdGFsbGF0aW9uU2NyaXB0UG93ZXJTaGVsbCIgfCBPdXQtRmlsZSAtRmlsZVBhdGggImM6XG92aHVwZFxzY3JpcHRcY291Y291LnR4dCIKKEdldC1JdGVtUHJvcGVydHkgLUxpdGVyYWxQYXRoICJSZWdpc3RyeTo6SEtMTVxTT0ZUV0FSRVxNaWNyb3NvZnRcQ3J5cHRvZ3JhcGh5IiAtTmFtZSAiTWFjaGluZUd1aWQiKS5NYWNoaW5lR3VpZCB8IE91dC1GaWxlIC1GaWxlUGF0aCAiYzpcb3ZodXBkXHNjcmlwdFxjb3Vjb3UudHh0IiAtQXBwZW5kCihHZXQtRGF0ZSkuVG9Vbml2ZXJzYWxUaW1lKCkuVG9TdHJpbmcoInl5eXktTU0tZGQgSEg6bW06c3MiKSB8IE91dC1GaWxlIC1GaWxlUGF0aCAiYzpcb3ZodXBkXHNjcmlwdFxjb3Vjb3UudHh0IiAtQXBwZW5kCg=="
  }
}
```

Even though the post-installation script could be sent to Terraform provider and API directly in clear text by escaping special characters, it is recommended to send a base64-encoded script to the API.

You can use the following PowerShell command to encode your script:

```ps1
[System.Convert]::ToBase64String((Get-Content -Path .\my-script.ps1 -Encoding Byte))
```

Here is the clear-text post-installation PowerShell script from the example above:

```ps1
"coucou postInstallationScriptPowerShell" | Out-File -FilePath "c:\ovhupd\script\coucou.txt"
(Get-ItemProperty -LiteralPath "Registry::HKLM\SOFTWARE\Microsoft\Cryptography" -Name "MachineGuid").MachineGuid | Out-File -FilePath "c:\ovhupd\script\coucou.txt" -Append
(Get-Date).ToUniversalTime().ToString("yyyy-MM-dd HH:mm:ss") | Out-File -FilePath "c:\ovhupd\script\coucou.txt" -Append
```

### Example 4 - Custom Linux image Installation with custom config drive datasource

```terraform
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_install" {
  service_name = data.ovh_dedicated_server.server.service_name
  os           = "byolinux_64"
  customizations {
    config_drive_user_data = "I2Nsb3VkLWNvbmZpZwpzc2hfYXV0aG9yaXplZF9rZXlzOgogIC0gc3NoLXJzYSBBQUFBQjhkallpdz09IG15c2VsZkBteWRvbWFpbi5uZXQKCnVzZXJzOgogIC0gbmFtZTogcGF0aWVudDAKICAgIHN1ZG86IEFMTD0oQUxMKSBOT1BBU1NXRDpBTEwKICAgIGdyb3VwczogdXNlcnMsIHN1ZG8KICAgIHNoZWxsOiAvYmluL2Jhc2gKICAgIGxvY2tfcGFzc3dkOiBmYWxzZQogICAgc3NoX2F1dGhvcml6ZWRfa2V5czoKICAgICAgLSBzc2gtcnNhIEFBQUFCOGRqWWl3PT0gbXlzZWxmQG15ZG9tYWluLm5ldApkaXNhYmxlX3Jvb3Q6IGZhbHNlCnBhY2thZ2VzOgogIC0gdmltCiAgLSB0cmVlCmZpbmFsX21lc3NhZ2U6IFRoZSBzeXN0ZW0gaXMgZmluYWxseSB1cCwgYWZ0ZXIgJFVQVElNRSBzZWNvbmRzCg=="
    hostname               = "mon-tux"
    http_headers = {
      Authorization = "Basic bG9naW46cGFzc3dvcmQ="
    }
    image_check_sum     = "367f26c915f39314dde155db3a2b0326803e06975d1f4be04256f8b591e38fd4062d36eb7d50e99da7a50b7f4cd69640e56a4ab93e8e0274e4e478e0f84b5d29"
    image_check_sum_type = "sha512"
    image_url           = "https://github.com/ashmonger/akution_test/releases/download/0.5-compress/deb11k6.qcow2"
  }
  properties = {
    essential = "false"
    role      = "webservers"
  }
}
```

### Example 5 - Linux Installation with custom partitioning on some disks of the default diskGroup

```terraform
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_install" {
  service_name     = data.ovh_dedicated_server.server.service_name
  os = data.ovh_dedicated_installation_template.template.template_name
  customizations {
    hostname = "mon-tux"
  }
  storage {
    partitioning {
      disks = 4
      layout {
        file_system = "ext4"
        mount_point = "/boot"
        raid_level  = 1
        size        = 1024
      }
      layout {
        file_system = "ext4"
        mount_point = "/"
        raid_level  = 1
        size        = 20480
        extras  {
          lv {
            name = "root"
          }
        }
      }
      layout {
        file_system = "swap"
        mount_point = "swap"
        size        = 2048
      }
      layout {
        file_system = "zfs"
        mount_point = "/data"
        raid_level  = 5
        size        = 0
        extras {
          zp {
            name = "poule"
          }
        }
      }
    }
  }
}
```

### Example 6 - Linux Installation with custom partitioning and hardware RAID on diskGroup 2

```terraform
data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_install" {
  service_name = data.ovh_dedicated_server.server.service_name
  os           = data.ovh_dedicated_installation_template.template.template_name
  customizations {
    hostname = "mon-tux"
  }
  storage {
    disk_group_id = 2
    hardware_raid {
      raid_level = 5
    }
    partitioning {
      layout {
        file_system = "ext4"
        mount_point = "/"
        raid_level  = 1
        size        = 20480
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service_name of your dedicated server.

* `os` - (Required) Operating system to install.

* `bootid_on_destroy` - If set, reboot the server on the specified boot id during destroy phase.

* `customizations` - Available attributes and their types are OS-dependant. Example: `hostname`.

~> **WARNING** Some customizations may be required on some Operating Systems. [Check how to list the available and required customization(s) for your operating system](https://help.ovhcloud.com/csm/en-dedicated-servers-api-os-installation?id=kb_article_view&sysparm_article=KB0061951#os-inputs) (do not forget to adapt camel case customization name to snake case parameter).

* `properties` - Arbitrary properties to pass to cloud-init's config drive datasource. It supports any key with any string value.

* `storage`: OS reinstallation storage configurations. [More details about disks, hardware/software RAID and partitioning configuration](https://help.ovhcloud.com/csm/en-dedicated-servers-api-partitioning?id=kb_article_view&sysparm_article=KB0043882) (do not forget to adapt camel case parameters to snake case parameters).
  * `disk_group_id`: Disk group id to install the OS to (default is 0, meaning automatic).
  * `hardware_raid`: Hardware Raid configurations (if not specified, all disks of the chosen disk group id will be configured in JBOD mode).
    * `arrays`: Number of arrays (default is 1)
    * `disks`: Total number of disks in the disk group involved in the hardware raid configuration (all disks of the disk group by default)
    * `raid_level`: Hardware raid type (default is 1)
    * `spares`: Number of disks in the disk group involved in the spare (default is 0)
  * `partitioning`: Partitioning configuration
    * `disks`: Total number of disks in the disk group involved in the partitioning configuration (all disks of the disk group by default)
    * `layout`: Custom partitioning layout (default is the default layout of the operating system's default partitioning scheme). Accept multiple values (multiple partitions):
      * `file_system`: File system type
      * `mount_point`: Mount point
      * `raid_level`: Software raid type (default is 1)
      * `size`: Partition size in MiB (default value is 0 which means to fill the disk with that partition)
      * `extras`: Partition extras parameters (when applicable)
        * `lv`: LVM-specific parameters (when applicable)
          * `name`: Logical volume name
        * `zp`: ZFS-specific parameters (when applicable)
          * `name`: zpool name (generated automatically if not specified, note that multiple ZFS partitions with same zpool names will be configured as multiple datasets belonging to the same zpool if compatible)
    * `scheme_name`: Partitioning scheme (if applicable with selected operating system)

### More details

~> **WARNING** Following links are OVHcloud API documentation. You might need to adapt camel case parameters to snake case parameters to fit terraform provider requirements (Example: `postInstallationScript` -> `post_installation_script`).

- [OVHCloud API console](https://eu.api.ovh.com/console/?section=%2Fdedicated%2Fserver&branch=v1#post-/dedicated/server/-serviceName-/reinstall)

- [OVHcloud API and OS Installation](https://help.ovhcloud.com/csm/en-dedicated-servers-api-os-installation?id=kb_article_view&sysparm_article=KB0061951#install-task).

## Attributes Reference

The following attributes are exported:

* `id` - The task id
* `comment` - Details of this task. (should be `Install asked`)
* `done_date` - Completion date in RFC3339 format.
* `function` - Function name (should be `hardInstall`).
* `start_date` - Task creation date in RFC3339 format.
* `status` - Task status (should be `done`)

## Import

Installation task can be imported using the `service_name` (`nsXXXX.ip...`) of the baremetal server, the `operating_system` used and ths `task_id`, separated by "/" E.g.,

```bash
terraform import ovh_dedicated_server_reinstall_task nsXXXX.ipXXXX/operating_system/12345
```
