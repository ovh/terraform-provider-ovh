resource "ovh_dedicated_nasha_partition_snapshot" "my_partition" {
  service_name = "zpool-12345"
  partition_name = "my-partition"
  type = "day-3"
}
