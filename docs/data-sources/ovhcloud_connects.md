---
subcategory : "Ovhcloud Connect (OCC)"
---

# ovh_ovhcloud_connects (Data Source)

Get the details of your Ovhcloud Connect products.

## Example Usage

```terraform
data "ovh_ovhcloud_connect" "occs" {}
```

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
