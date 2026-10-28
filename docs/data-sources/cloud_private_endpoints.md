---
subcategory: "Private Endpoint"
---

# ovh_cloud_private_endpoints (Data Source)

Use this data source to list the private endpoints of a public cloud project.

To discover which services *can* be exposed through a private endpoint rather
than which ones already are, use
[`ovh_cloud_private_endpoint_service_availabilities`](../data-sources/cloud_private_endpoint_service_availabilities.md).

## Example Usage

```terraform
data "ovh_cloud_private_endpoints" "endpoints" {
  service_name = "<public cloud project ID>"
}

output "endpoint_hostnames" {
  value = {
    for endpoint in data.ovh_cloud_private_endpoints.endpoints.private_endpoints :
    endpoint.name => endpoint.current_state.hostname
  }
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.

## Attributes Reference

* `private_endpoints` - List of private endpoints:
  * `id` - Unique identifier of the private endpoint.
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
