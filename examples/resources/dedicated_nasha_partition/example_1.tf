resource "ovh_dedicated_nasha_partition" "my_partition" {
  service_name = "zpool-12345"
  name = "my-partition"
  size = 20
  protocol = "NFS"
}
