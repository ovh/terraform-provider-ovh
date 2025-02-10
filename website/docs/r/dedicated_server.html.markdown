---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server

Use this resource to order and manage a dedicated server.

## Example Usage

```hcl
data "ovh_me" "account" {}

resource "ovh_dedicated_server" "server" {
  ovh_subsidiary = data.ovh_me.account.ovh_subsidiary
  display_name = "My server display name"
  operating_system = "debian12_64"
  plan = [
    {
      plan_code = "22rise01"
      duration = "P1M"
      pricing_mode = "default"

      configuration = [
        {
          label = "dedicated_datacenter"
          value = "bhs"
        },
        {
          label = "dedicated_os"
          value = "none_64.en"
        },
        {
          label = "region"
          value = "canada"
        }
      ]
    }
  ]

  plan_option = [
    {
      duration = "P1M"
      plan_code = "ram-32g-rise13"
      pricing_mode = "default"
      quantity = 1
    },
    {
      duration = "P1M"
      plan_code = "bandwidth-500-included-rise"
      pricing_mode = "default"
      quantity = 1
    },
    {
      duration = "P1M"
      plan_code = "softraid-2x512nvme-rise"
      pricing_mode = "default"
      quantity = 1
    },
    {
      duration = "P1M"
      plan_code = "vrack-bandwidth-100-rise-included"
      pricing_mode = "default"
      quantity = 1
    }
  ]
}
```

~> __WARNING__ After ordering a dedicated server, the provider will wait for 1 hour for it to be delivered. If it is still not delivered after this time, the apply will end in error, but the delivery process will still continue on OVHcloud's side. In this case you just need to manually untaint the resource and re-run an apply: `terraform untaint ovh_dedicated_server.server && terraform apply`. This can be repeated as many times as needed while waiting for the server to be delivered.

## Argument Reference

### Arguments used to order a dedicated server

* `ovh_subsidiary` - (Required) OVHcloud Subsidiary. Country of OVHcloud legal entity you'll be billed by. List of supported subsidiaries available on API at [/1.0/me.json under `models.nichandle.OvhSubsidiaryEnum`](https://eu.api.ovh.com/1.0/me.json)
* `plan` - (Required) Product Plan to order
  * `duration` - (Required) Duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing mode identifier
  * `catalog_name` - Catalog name
  * `configuration` - Representation of a configuration item to personalize product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in API.OVH.COM
* `plan_option` - Product Plan to order
  * `duration` - (Required) duration
  * `plan_code` - (Required) Plan code
  * `pricing_mode` - (Required) Pricing model identifier
  * `catalog_name` - Catalog name
  * `configuration` - Representation of a configuration item to personalize product
    * `label` - (Required) Identifier of the resource
    * `value` - (Required) Path to the resource in API.OVH.COM

### Editable fields of a dedicated server

* `display_name` - Display name of your dedicated server
* `monitoring` - Icmp monitoring state
* `no_intervention` - Prevent datacenter intervention
* `rescue_mail` - Custom email used to receive rescue credentials
* `rescue_ssh_key` - Public SSH Key used in the rescue mode
* `root_device` - Root device of the server
* `state` - All states a Dedicated can in (error, hacked, hackedBlocked, ok)

### Arguments used to reinstall a dedicated server
* `customizations` - Customization of the OS configuration
  * `configDriveUserData` -Config Drive UserData
  * `efiBootloaderPath` - Path of the EFI bootloader from the OS installed on the server
  * `hostname` - Custom hostname
  * `httpHeaders` - Image HTTP Headers
  * `imageCheckSum` - Image checksum
  * `imageCheckSumType` - Checksum type
  * `imageType` - Image Type
  * `imageURL` - Image URL
  * `language` - Display Language
  * `postInstallationScript` - Post-Installation Script
  * `postInstallationScriptExtension` - Post-Installation Script File Extension
  * `sshKey` - SSH Public Key
* `storage` - Storage customization
  * `diskGroupId` -  Disk group id
  * `hardwareRaid` - Hardware Raid configurations
    * `arrays` - Number of arrays
    * `disks` - Total number of disks in the disk group involved in the hardware raid configuration
    * `raidLevel` - Hardware raid type
    * `spares` -  Number of disks in the disk group involved in the spare
  * `partitioning` - Partitioning configuration
    * `disks` - Total number of disks in the disk group involved in the partitioning configuration
    * `layout` - Custom partitioning layout 
      * `extras` - Partition extras parameters
        * `lv` - LVM-specific parameters
        * `zp` - ZFS-specific parameters
      * `fileSystem` -  File system type
      * `mountPoint` - Mount point
      * `raidLevel` - Software raid type
      * `size` - Partition size in MiB
    * `schemeName` - Partitioning scheme (if applicable with selected operating system)
* `properties` - Arbitrary properties to pass to cloud-init's config drive datasource

## Attributes Reference

* `service_name` - The service_name of your dedicated server
* `display_name` - Dedicated server display name
* `name` - Dedicated server name
* `server_id` - Server id
* `commercial_range` - Dedicated server commercial range
* `os` - Operating system
* `ip` - Dedicated server ip (IPv4)
* `region` - Dedicated region localisation
* `availability_zone` - Dedicated AZ localisation
* `datacenter` - Dedicated datacenter localisation (bhs1,bhs2,...)
* `rack` - Rack id of the server
* `state` - All states a Dedicated can be in (error, hacked, hackedBlocked, ok)
* `power_state` - Power state of the server (poweron, poweroff)
* `support_level` - Dedicated server support level (critical, fastpath, gs, pro)
* `iam` - IAM resource information
  * `urn` - URN of the private database, used when writing IAM policies
  * `display_name` - Resource display name
  * `id` - Unique identifier of the resource in the IAM
  * `tags` - Resource tags. Tags that were internally computed are prefixed with `ovh:`
* `boot_id` - Boot id of the server
* `boot_script` - Boot script of the server
* `efi_bootloader_path` - Path of the EFI bootloader
* `link_speed` - Link speed of the server
* `monitoring` - Icmp monitoring state
* `no_intervention` - Prevent datacenter intervention
* `new_upgrade_system` -
* `partition_scheme_name` - Partition scheme name
* `professional_use` - Does this server have professional use option
* `rescue_mail` - Rescue mail of the server
* `rescue_ssh_key` - Public SSH Key used in the rescue mode
* `reverse` - Dedicated server reverse
* `root_device` - Root device of the server

## Import

Dedicated servers can be imported using the `service_name`.
Using the following configuration:

```hcl
import {
  to = ovh_dedicated_server.server
  id = "<service name>"
}
```

You can then run:

```bash
terraform plan -generate-config-out=dedicated.tf
terraform apply
```

The file `dedicated.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above.
See <https://developer.hashicorp.com/terraform/language/import/generating-configuration> for more details.
