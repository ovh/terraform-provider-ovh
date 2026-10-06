---
subcategory : "Cloud Disk Array"
---

# ovh_dedicated_ceph_user

Creates a user in a dedicated CEPH cluster (Cloud Disk Array) and manages its permissions on the pools of the cluster.

## Example Usage

```terraform
data "ovh_dedicated_ceph" "my_ceph" {
  service_name = "94d423da-0e55-45f2-9812-836460a19939"
}

resource "ovh_dedicated_ceph_pool" "my_pool" {
  service_name = data.ovh_dedicated_ceph.my_ceph.id
  name         = "my-pool"
}

resource "ovh_dedicated_ceph_user" "my_user" {
  service_name = data.ovh_dedicated_ceph.my_ceph.id
  name         = "my-user"

  pool_permissions = [
    {
      pool_name = ovh_dedicated_ceph_pool.my_pool.name
      read      = true
      write     = true
    },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your dedicated CEPH. **Changing this value recreates the resource.**
* `name` - (Required) Name of the user. **Changing this value recreates the resource.**
* `pool_permissions` - (Optional) List of the permissions of the user on the pools of the cluster, one per pool. Each entry must grant at least one permission. The permissions of the user on the pools not listed are cleared. Changing the order of the list makes an in-place update that sets the same permissions.
  * `pool_name` - (Required) Name of the pool.
  * `read` - (Optional) Read permission. Defaults to `false`.
  * `write` - (Optional) Write permission. Defaults to `false`.
  * `execute` - (Optional) Execute permission. Defaults to `false`.
  * `class_read` - (Optional) Class read permission. Defaults to `false`.
  * `class_write` - (Optional) Class write permission. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - Identifier of the user, formatted as `service_name/name`.
* `key` - (Sensitive) Key of the user to connect to the cluster.
* `mon_caps` - Capabilities of the user on the MON daemons.
* `osd_caps` - Capabilities of the user on the OSD daemons, derived from its pool permissions.
* `mds_caps` - Capabilities of the user on the MDS daemons.

## Import

A dedicated CEPH user can be imported using the `service_name` and the user `name`, separated by "/" E.g.,

```bash
$ terraform import ovh_dedicated_ceph_user.my_user service_name/user_name
```
