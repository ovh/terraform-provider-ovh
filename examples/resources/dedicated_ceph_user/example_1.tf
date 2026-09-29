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
