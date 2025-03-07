---
subcategory : "vRack"
---

# ovh_cloud_project_vrack (Data Source)

Use this data source to get the linked vrack on your public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_vrack" "vrack" {
  service_name = "XXXXXX"
}

output "vrack" {
  value = data.ovh_cloud_project_vrack.vrack
}
```

## Argument Reference

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

## Attributes Reference

The following attributes are exported
* `id` - The id of the vrack
* `name` - The name of the vrack
* `decription` - The description of the vrack
* `service_name` - The id of the public cloud project
