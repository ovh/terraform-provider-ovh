---
subcategory : "Cloud Disk Array"
---

# ovh_dedicated_ceph_pool

Creates a pool in a dedicated CEPH cluster (Cloud Disk Array).

Destroying the resource deletes the pool and all of its data: the provider opens the deletion window of the pool, then deletes it.

## Example Usage

```terraform
data "ovh_dedicated_ceph" "my_ceph" {
  service_name = "94d423da-0e55-45f2-9812-836460a19939"
}

resource "ovh_dedicated_ceph_pool" "my_pool" {
  service_name = data.ovh_dedicated_ceph.my_ceph.id
  name         = "my-pool"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your dedicated CEPH. **Changing this value recreates the resource.**
* `name` - (Required) Name of the pool. **Changing this value recreates the resource.**

## Attributes Reference

The following attributes are exported:

* `id` - Identifier of the pool, formatted as `service_name/name`.
* `pool_type` - Type of the pool (`ERASURE_CODED` or `REPLICATED`).
* `backup` - Whether the pool is backed up.
* `min_active_replicas` - Minimum number of active replicas.
* `replica_count` - Number of replicas.

## Import

A dedicated CEPH pool can be imported using the `service_name` and the pool `name`, separated by "/" E.g.,

```bash
$ terraform import ovh_dedicated_ceph_pool.my_pool service_name/pool_name
```
