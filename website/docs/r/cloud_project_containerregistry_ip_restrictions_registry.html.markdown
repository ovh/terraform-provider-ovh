---
subcategory : "Managed Private Registry"
---

# ovh_cloud_project_containerregistry_ip_restrictions_registry

Apply IP restrictions container registry associated with a public cloud project on artifact manager component.

## Example Usage

```hcl
data "ovh_cloud_project_containerregistry" "registry" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

resource "ovh_cloud_project_containerregistry_ip_restrictions_registry" "my-registry-iprestrictions" {
  service_name = ovh_cloud_project_containerregistry.registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.registry.id

  ip_restrictions = [
    { 
      ip_block = "xxx.xxx.xxx.xxx/xx" 
      description = "xxxxxxx"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `registry_id` - The id of the Managed Private Registry.
* `ip_restrictions` - IP restrictions applied on artifact manager component.
    * `description` - The Description of Whitelisted IpBlock.
    * `ip_block` - Whitelisted IpBlock (CIDR format).

## Attributes Reference

The following attributes are exported:

* `service_name` - The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `registry_id` - The ID of the Managed Private Registry.
* `ip_restrictions` - IP restrictions applied on artifact manager component.
    * `description` - The Description of Whitelisted IpBlock.
    * `ip_block` - Whitelisted IpBlock (CIDR format).