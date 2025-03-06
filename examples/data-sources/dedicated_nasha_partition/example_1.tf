data "ovh_dedicated_nasha_partition" "my_nas_ha_partition" {
  service_name = "zpool-12345"
  name         = "my-zpool-partition"
}
