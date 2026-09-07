---
subcategory: "Private Endpoint"
---

# ovh_cloud_private_endpoint (Data Source)

Use this data source to retrieve information about a private endpoint in a
public cloud project, including the placement and the address the platform
resolved for it.

## Example Usage

```terraform
data "ovh_cloud_private_endpoint" "endpoint" {
  service_name = "<public cloud project ID>"
  id           = "<private endpoint ID>"
}

output "endpoint_hostname" {
  value = data.ovh_cloud_private_endpoint.endpoint.current_state.hostname
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `id` - (Required) The id of the private endpoint.

## Attributes Reference

* `name` - Private endpoint name.
* `description` - Private endpoint description.
* `region` - Region where the private endpoint is created.
* `service` - Service exposed through the endpoint:
  * `name` - Name of the service.
  * `tier` - Tier of the service.
* `network_id` - Network the endpoint port is allocated in.
* `subnet_id` - Subnet the endpoint port is allocated in.
* `vip_address` - Private IP address pinned on the endpoint port.
* `public_dns_enabled` - Whether the endpoint hostname is registered in the public DNS zone.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the private endpoint.
* `updated_at` - Last update date of the private endpoint.
* `resource_status` - Private endpoint readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the private endpoint:
  * `name` - Private endpoint name.
  * `description` - Private endpoint description.
  * `endpoint_id` - Identifier of the endpoint on the service side, as `pe-<id>`.
  * `location` - Location details:
    * `region` - Region.
  * `service` - Service exposed through the endpoint:
    * `name` - Name of the service.
    * `tier` - Tier of the service.
    * `port` - Port the service listens on through the endpoint.
    * `transport` - Transport protocol of the service.
    * `bandwidth` - Bandwidth of the tier, in Mbit/s. Null when the catalog does not declare it.
  * `network` - Network the endpoint port is attached to:
    * `id` - Network identifier.
  * `subnet` - Subnet the endpoint port is attached to:
    * `id` - Subnet identifier.
  * `vip_address` - Private IP address of the endpoint in the customer network.
  * `hostname` - Fully qualified domain name of the endpoint. Null until the service publishes it. See `public_dns_enabled` for public resolvability.
  * `public_dns_enabled` - Whether the endpoint hostname is registered in the public DNS zone.
