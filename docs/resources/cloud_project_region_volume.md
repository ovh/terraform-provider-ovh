---
subcategory : "Cloud Project"
---

# ovh_cloud_project_volume

Create volume in a public cloud project.

## Example Usage

Create a subscription

```terraform
resource "ovh_cloud_project_volume" "vol" {
   region_name = "xxx"
   service_name = "yyyyy"
   description = "Terraform volume"
   name = "terrformName"
   size = 15
   type = "classic"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - Required. The id of the public cloud project. **Changing this value recreates the resource.**
* `region_name` - Required. A valid OVHcloud public cloud region name in which the volume will be available. Ex.: "GRA11". **Changing this value recreates the resource.**
* `description` - A description of the volume
* `name` - Name of the volume
* `size` - Size (GB) of the volume
* `type` - Type of the volume **Changing this value recreates the resource.**

## Attributes Reference

The following attributes are exported:

* `service_name` - The id of the public cloud project.
* `region_name` - A valid OVHcloud public cloud region name in which the volume will be available.
* `description` - A description of the volume
* `name` - Name of the volume
* `size` - Size of the volume
* `id` - id of the volume
