---
subcategory : "Gateway"
---

# ovh_cloud_gateway (Data Source)

Use this data source to retrieve information about a gateway in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_gateway" "gateway" {
  service_name = "<public cloud project ID>"
  id           = "<gateway ID>"
}

output "gateway_name" {
  value = data.ovh_cloud_gateway.gateway.name
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `id` - (Required) Gateway ID.

## Attributes Reference

The following attributes are exported:

* `region` - Region of the gateway.
* `availability_zone` - Availability zone of the gateway.
* `name` - Gateway name.
* `description` - Gateway description.
* `external_gateway` - External gateway configuration:
  * `enabled` - Whether the external gateway is enabled.
  * `model` - External gateway sizing model.
* `subnet_ids` - List of subnet IDs attached as router interfaces.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the gateway.
* `updated_at` - Last update date of the gateway.
* `resource_status` - Gateway readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the gateway:
  * `name` - Gateway name.
  * `description` - Gateway description.
  * `status` - OpenStack router status (`ACTIVE`, `BUILD`, `DOWN`, `ERROR`).
  * `external_gateway` - External gateway configuration:
    * `enabled` - Whether the external gateway is enabled.
    * `model` - External gateway sizing model.
  * `external_ip` - External IP address assigned to the gateway.
  * `subnets` - Currently attached subnets:
    * `id` - Subnet ID.
  * `region` - Region.
  * `availability_zone` - Availability zone.
