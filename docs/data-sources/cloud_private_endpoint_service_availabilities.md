---
subcategory: "Private Endpoint"
---

# ovh_cloud_private_endpoint_service_availabilities (Data Source)

List the services reachable through a private endpoint in a public
cloud project, with their enabled tiers.

## Example Usage

```terraform
data "ovh_cloud_private_endpoint_service_availabilities" "availabilities" {
  service_name = "<public cloud project ID>"
  region       = "GRA11"
}
```

## Argument Reference

* `service_name` - (Required) The id of the public cloud project.
* `region` - (Optional) Region filter. When set, only tiers available in that region are returned.

## Attributes Reference

* `service_availabilities` - List of discoverable service availabilities:
  * `name` - Service name.
  * `description` - Service description.
  * `tiers` - Enabled tiers of the service:
    * `name` - Tier name.
    * `region` - Region where this tier is available.
    * `transport` - Transport protocol of the service.
    * `port` - Port the service listens on through the endpoint.
    * `bandwidth` - Bandwidth of the tier, in Mbit/s. Null when the catalog does not declare it.
