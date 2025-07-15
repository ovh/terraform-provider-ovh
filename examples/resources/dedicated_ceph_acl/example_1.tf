data "ovh_dedicated_ceph" "my_ceph" {
  service_name = "94d423da-0e55-45f2-9812-836460a19939"
}

resource "ovh_dedicated_ceph_acl" "my_acl" {
  service_name = data.ovh_dedicated_ceph.my_ceph.id
  network      = "1.2.3.4"
  netmask      = "255.255.255.255"
}
