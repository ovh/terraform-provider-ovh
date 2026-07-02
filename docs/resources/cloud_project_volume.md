---
subcategory: "Block Storage"
---

# ovh_cloud_project_volume

~> **DEPRECATED** This resource is deprecated. Use `ovh_cloud_storage_block_volume` instead.

Create volume in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_volume" "volume" {
   region_name  = "xxx"
   service_name = "yyyyy"
   description  = "Terraform volume"
   name         = "terrformName"
   size         = 15
   type         = "classic"
}
```

### Encrypted volume with a customer managed key (CMK)

```terraform
resource "ovh_cloud_project_volume" "encrypted_volume" {
   region_name  = "xxx"
   service_name = "yyyyy"
   description  = "Terraform encrypted volume"
   name         = "encryptedVolume"
   size         = 15
   type         = "classic"

   encryption = {
     encrypted = true
     kms = {
       domain_id      = "<okms domain id>"
       service_key_id = "<okms service key id>"
     }
   }
}
```

Omit the `kms` block to encrypt the volume with OVH managed keys (OMK):

```terraform
resource "ovh_cloud_project_volume" "encrypted_volume" {
   region_name  = "xxx"
   service_name = "yyyyy"
   name         = "encryptedVolume"
   size         = 15
   type         = "classic"

   encryption = {
     encrypted = true
   }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - Optional. The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `region_name` - Required. A valid OVHcloud public cloud region name in which the volume will be available. Ex.: "GRA11". **Changing this value recreates the resource.**
* `availability_zone` - Optional. Availability zone in which the volume is created. Required when `region_name` is a 3AZ region. **Changing this value recreates the resource.**
* `description` - A description of the volume
* `name` - Name of the volume
* `size` - Size (GB) of the volume
* `type` - Type of the volume **Changing this value recreates the resource.** Available types are: classic, classic-luks, classic-multiattach, high-speed, high-speed-luks, high-speed-gen2, high-speed-gen2-luks
* `encryption` - Optional. Volume encryption configuration. Customer managed keys (CMK) are only available in supported regions (3AZ). **Changing this value recreates the resource.**
  * `encrypted` - Whether the volume is encrypted. Setting this auto-derives a LUKS volume type.
  * `kms` - Optional. Customer managed key (CMK) reference. Omit to use OVH managed keys (OMK).
    * `domain_id` - OKMS domain ID holding the customer managed key.
    * `service_key_id` - OKMS service key ID used to encrypt the volume.

## Attributes Reference

The following attributes are exported:

* `service_name` - The id of the public cloud project.
* `region_name` - A valid OVHcloud public cloud region name in which the volume will be available.
* `description` - A description of the volume
* `name` - Name of the volume
* `size` - Size of the volume
* `id` - id of the volume
* `encryption` - Volume encryption configuration (see above).

## Import

The resource can be imported using the public cloud project ID, region and the volume ID, e.g.,

```terraform
import {
  to = ovh_cloud_project_volume.volume
  id = "<public cloud project ID>/<region>/<volume ID>"
}
```

```bash
$ terraform plan -generate-config-out=volume.tf
$ terraform apply
```

The file `volume.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above. See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.