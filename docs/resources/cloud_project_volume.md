---
subcategory : "Cloud Project"
---

# ovh_cloud_project_volume

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

## Argument Reference

The following arguments are supported:

* `service_name` - Optional. The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `region_name` - Required. A valid OVHcloud public cloud region name in which the volume will be available. Ex.: "GRA11". **Changing this value recreates the resource.**
* `description` - A description of the volume
* `name` - Name of the volume
* `size` - Size (GB) of the volume
* `type` - Type of the volume **Changing this value recreates the resource.** Available types are: classic, classic-luks, classic-multiattach, high-speed, high-speed-luks, high-speed-gen2, high-speed-gen2-luks

## Attributes Reference

The following attributes are exported:

* `service_name` - The id of the public cloud project.
* `region_name` - A valid OVHcloud public cloud region name in which the volume will be available.
* `description` - A description of the volume
* `name` - Name of the volume
* `size` - Size of the volume
* `id` - id of the volume

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