---
subcategory : "Key Manager"
---

# ovh_cloud_keymanager_container_consumer

Registers a consumer on a Barbican Key Manager container for a public cloud project. Consumers track which OpenStack resources (instances, load balancers, images) are using a given container.

~> All fields on this resource force replacement. Consumers cannot be updated in place.

## Example Usage

Register a load balancer as a consumer of a container:

```terraform
resource "ovh_cloud_keymanager_container_consumer" "consumer" {
  service_name  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  container_id  = "00000000-0000-0000-0000-000000000000"
  service       = "LOADBALANCER"
  resource_type = "LOADBALANCER"
  resource_id   = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `container_id` - (Required, Forces new resource) UUID of the container to register the consumer on.
* `service` - (Required, Forces new resource) OpenStack service type of the consumer. Possible values: `COMPUTE`, `IMAGE`, `LOADBALANCER`, `NETWORK`.
* `resource_type` - (Required, Forces new resource) Type of the resource consuming the container. Possible values: `IMAGE`, `INSTANCE`, `LOADBALANCER`.
* `resource_id` - (Required, Forces new resource) UUID of the resource consuming the container.

## Attributes Reference

The following attributes are exported:

* `id` - Consumer ID (composite: `service_name/container_id/service/resource_type/resource_id`).

## Import

A Key Manager container consumer can be imported using a composite ID:

```bash
$ terraform import ovh_cloud_keymanager_container_consumer.my_consumer service_name/container_id/service/resource_type/resource_id
```
