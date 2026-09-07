---
subcategory: "Private Endpoint"
---

# ovh_cloud_private_endpoint

Creates a private endpoint in a public cloud project.

A private endpoint exposes a service inside a customer network through a
private IP address. Placement is expressed with `network_id`, `subnet_id`, or
both; at least one is required. The platform resolves the one that is omitted
and exports it as a computed value, so `network_id` and `subnet_id` always end
up populated in the state:

* `network_id` alone: the network must hold exactly one IPv4 subnet, which is
  used for the endpoint.
* `subnet_id` alone: the subnet's parent network is resolved.
* both: the subnet must belong to the network, otherwise the API rejects the
  request with a `400`.

The port carrying the endpoint address is always created by the platform.

## Example Usage

```terraform
resource "ovh_cloud_private_endpoint" "endpoint" {
  service_name       = "<public cloud project ID>"
  name               = "my-endpoint"
  description        = "Endpoint for the service"
  region             = "GRA11"
  public_dns_enabled = false
  subnet_id          = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  service = {
    name = "s3"
    tier = "S"
  }
}
```

### Let the platform resolve the subnet

```terraform
resource "ovh_cloud_private_endpoint" "endpoint" {
  service_name       = "<public cloud project ID>"
  name               = "my-endpoint"
  region             = "GRA11"
  public_dns_enabled = false
  network_id         = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  service = {
    name = "s3"
    tier = "S"
  }
}
```

### Pin the endpoint address

```terraform
resource "ovh_cloud_private_endpoint" "endpoint" {
  service_name       = "<public cloud project ID>"
  name               = "my-endpoint"
  region             = "GRA11"
  public_dns_enabled = false
  subnet_id          = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  vip_address        = "10.0.0.42"

  service = {
    name = "s3"
    tier = "S"
  }
}
```

## Compatible services

The `ovh_cloud_private_endpoint_service_availabilities` data source lists the
services and tiers enabled in each region

### S3

An endpoint towards S3 gives the subnet a private address on which the S3 API
answers. To force a user through it, restrict its S3 policy with an
`aws:SourceVpce` condition matched against the identifier the service assigns
to the endpoint, exported as `current_state.endpoint_id`. The `Deny` policy
narrows the user to requests entering through the endpoint:

```terraform
resource "ovh_cloud_private_endpoint" "s3" {
  service_name       = "<public cloud project ID>"
  name               = "s3-endpoint"
  region             = "GRA11"
  public_dns_enabled = false
  subnet_id          = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  service = {
    name = "s3"
    tier = "S"
  }
}

resource "ovh_cloud_project_user" "app" {
  service_name = "<public cloud project ID>"
  description  = "Reaches the bucket through the private endpoint only"
  role_names   = ["objectstore_operator"]
}

resource "ovh_cloud_project_user_s3_policy" "app" {
  service_name = ovh_cloud_project_user.app.service_name
  user_id      = ovh_cloud_project_user.app.id
  policy = jsonencode({
    Statement = [{
      Sid      = "DenyUnlessViaPrivateEndpoint"
      Effect   = "Deny"
      Action   = "s3:*"
      Resource = "*"
      Condition = {
        StringNotEquals = {
          "aws:SourceVpce" = ovh_cloud_private_endpoint.s3.current_state.endpoint_id
        }
      }
    }]
  })
}
```

Requests made with this user's credentials from outside the private endpoint
are denied. The same requests going through the private endpoint succeed, with
the S3 client pointed at the endpoint address.

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `name` - (Required) Private endpoint name, unique in the project for a given subnet (case-insensitive).
* `description` - (Optional) Private endpoint description.
* `region` - (Required) Region where the private endpoint is created. **Changing this value recreates the resource.**
* `service` - (Required) Service to expose through the endpoint. **Changing this value recreates the resource.**
  * `name` - (Required) Name of the service.
  * `tier` - (Required) Tier of the service. See the `ovh_cloud_private_endpoint_service_availabilities` data source for the tiers enabled per service and region.
* `network_id` - (Optional, Computed) Network to allocate the endpoint port in. At least one of `network_id` or `subnet_id` must be set; both may be set, and the subnet must then belong to the network. When omitted, the network owning the given subnet is resolved by the platform and exported as a computed value. **Changing this value recreates the resource.**
* `subnet_id` - (Optional, Computed) Subnet to allocate the endpoint port in. At least one of `network_id` or `subnet_id` must be set; both may be set, and the subnet must then belong to the network. When omitted, the network must hold exactly one IPv4 subnet, which is resolved by the platform and exported as a computed value. **Changing this value recreates the resource.**
* `vip_address` - (Optional, Computed) Private IP address to pin on the endpoint port. Requires `subnet_id` to be set in the same configuration. When omitted, an address is allocated from the subnet allocation pool and exported as a computed value. The address must be inside the subnet CIDR, inside one of its allocation pools, must not be the subnet gateway address, and must not already be used by another port; a violation is rejected with a `400`, a collision with an existing port with a `409`. **Changing this value recreates the resource.**
* `public_dns_enabled` - (Required) Register the endpoint hostname in the public DNS zone. Registering publishes the endpoint hostname in a public zone, so the choice must be stated explicitly. **Changing this value recreates the resource.**

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the private endpoint.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the private endpoint.
* `updated_at` - Last update date of the private endpoint.
* `resource_status` - Private endpoint readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the private endpoint:
  * `name` - Private endpoint name.
  * `description` - Private endpoint description.
  * `endpoint_id` - Identifier of the endpoint on the service side, as `pe-<id>`. The value to match in an S3 policy `aws:SourceVpce` condition.
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

## Import

A private endpoint can be imported using the `service_name` and `private_endpoint_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_private_endpoint.endpoint
  id = "<service_name>/<private_endpoint_id>"
}
```

```bash
$ terraform plan -generate-config-out=endpoint.tf
$ terraform apply
```
