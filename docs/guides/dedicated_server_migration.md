---
page_title: "Migrating dedicated servers from previous versions to v2.0.0"
---

Version v2.0.0 of OVHcloud Terraform provider introduced a breaking change on the resources related to dedicated servers, mainly:
- Deletion of resource `ovh_dedicated_server_install_task` that has been replaced by a new resource `ovh_dedicated_server_reinstall_task`
- The parameters of resource `ovh_dedicated_server` have been updated to reflect the changes on parameters needed to reinstall a dedicated server

The complete changelog can be found [here](https://github.com/ovh/terraform-provider-ovh/releases/tag/v2.0.0).

This guide explains how to migrate your existing configuration relying on resource `ovh_dedicated_server_install_task` to a new configuration compatible with version v2.0.0 of the provider without triggering a reinstallation of your dedicated servers.

!> This documentation presents a direct migration path between v1.6.0 and v2.0.0. To make this transition easier, we released a version v1.7.0 of the provider that includes both deprecated resources and the new ones. The migration steps are the same, but this version allows a more gradual shift towards v2.0.0.

## First step: import your dedicated servers in the state

From version v2.0.0 and later, the preferred way to manage dedicated servers installation details is through usage of resource `ovh_dedicated_server`. As a result, if you don't already have your dedicated servers declared in your Terraform configuration, you must import them.

You can use an `import` block like the following:

```terraform
import {
  id = "nsxxxxxxx.ip-xx-xx-xx.eu"
  to = ovh_dedicated_server.srv
}

resource "ovh_dedicated_server" "srv" {}
```

-> If you are doing a first migration to v1.7.0, you should add the following parameter `prevent_install_on_import = true` to the dedicated server resource. This guarantees you that the server won't be reinstalled after import, even if you have a diff on the reinstall-related parameters.

To finish importing the resource into your Terraform state, you should run:

```sh
terraform apply
```

## Second step: backport your previous task details into the imported resource

This step is manual and requires you to convert the previous installation details from resource `ovh_dedicated_server_install_task` to the new fields of resource `ovh_dedicated_server`: `os`, `customizations`, `properties` and `storage`.

Let's take an example: if you previously used the following configuration:

```terraform
resource "ovh_dedicated_server_install_task" "server_install" {
  service_name      = "nsxxxxxxx.ip-xx-xx-xx.eu"
  template_name     = "debian12_64"
  details {
      custom_hostname = "mytest"
  }
  user_metadata {
    key = "sshKey"
    value = "ssh-ed25519 AAAAC3..."
  }
  user_metadata {
    key = "postInstallationScript"
    value = <<-EOF
        #!/bin/bash
          echo "coucou postInstallationScript" > /opt/coucou
          cat /etc/machine-id  >> /opt/coucou
          date "+%Y-%m-%d %H:%M:%S" --utc >> /opt/coucou
        EOF
  }
}
```

You can replace it by the following one:

```terraform
resource "ovh_dedicated_server" "srv" {
  customizations = {
    hostname                 = "mytest"
    post_installation_script = "IyEvYmluL2Jhc2gKZWNobyAiY291Y291IHBvc3RJbnN0YWxsYXRpb25TY3JpcHQiID4gL29wdC9jb3Vjb3UKY2F0IC9ldGMvbWFjaGluZS1pZCAgPj4gL29wdC9jb3Vjb3UKZGF0ZSAiKyVZLSVtLSVkICVIOiVNOiVTIiAtLXV0YyA+PiAvb3B0L2NvdWNvdQo="
    ssh_key                  = "ssh-ed25519 AAAAC3..."
  }
  os         = "debian12_64"
  properties = null
  storage    = null
}
```

You can check the documentation of resource `ovh_dedicated_server` to see what inputs are available for the reinstallation-related fields.
The documentation of resource `ovh_dedicated_server_reinstall_task` also includes several examples of configuration.

## Third step: make sure your server is not reinstalled unintentionally

You should add the following piece of configuration in the declaration of your dedicated server resource in order to avoid a reinstallation on the next `terraform apply`:

```terraform
resource "ovh_dedicated_server" "srv" {
  #
  # ... resource fields
  #

  lifecycle {
    ignore_changes = [os, customizations, properties, storage]
  }
}
```

This is needed because there is no API endpoint that returns the previous installation parameters of a dedicated server. The goal here is to migrate your old configuration to the new format without triggering a reinstallation.

## Fourth step: remove the lifecycle block

After a while, whenever you need to trigger a reinstallation of your dedicated server, you can just remove the `lifecycle` field from your configuration and run `terraform apply`.