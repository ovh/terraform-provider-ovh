---
subcategory : "Ovhcloud Connect (OCC)"
---

# ovh_ovhcloud_connect (Data Source)

Use this data source to retrieve information about an Ovhcloud Connect product

## Example Usage

```terraform
data "ovh_ovhcloud_connect" "occ" {
  service_name = "XXX"
}
```

## Argument Reference

The following arguments are supported:

- `service_name` - (Required) The uuid of the Ovhcloud connect.

## Attributes Reference

The following attributes are exported:

- `uuid` - uuid of the Ovhcloud Connect service
- `bandwidth` - Service bandwidth
- `description` - Service description
- `status` - Service status
- `interface_list` - List of interfaces linked to a service
- `pop` - Pop reference where the service is delivered
- `port_quantity` - Port quantity
- `product` - Product name of the service
- `provider_name` - Service provider
- `service_name` - Service name
- `vrack` - vrack linked to the service
- `iam` - IAM resource information
  - `urn` - URN of the private database, used when writing IAM policies
  - `display_name` - Resource display name
  - `id` - Unique identifier of the resource in the IAM
  - `tags` - Resource tags. Tags that were internally computed are prefixed with `ovh:`
