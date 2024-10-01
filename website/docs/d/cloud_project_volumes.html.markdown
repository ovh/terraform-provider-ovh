---
subcategory : "Cloud Project Volume"
---

# ovh_cloud_project_volumes

Get all the volume from a region of a public cloud project

## Example Usage

```hcl
data "ovh_cloud_project_volume" "volume" {
   region_name = "xxx"
   service_name = "yyy"
   volume_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}
```

## Argument Reference


* `service_name` - (Required) The id of the public cloud project.
* `region_name` - (Required) A valid OVHcloud public cloud region name in which the volume is available. Ex.: "GRA11".
* `volume_id` - (Required) Volume id to get the informations

## Attributes Reference
* `volumes` - 
  * `name` - The name of the volume
  * `size` - The size of the volume
  * `id` - The id of the volume
* `region_name` - The region name where volume is available
* `service_name` - The id of the public cloud project.
