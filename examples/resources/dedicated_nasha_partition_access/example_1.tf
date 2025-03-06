resource "ovh_dedicated_nasha_partition_access" "my_partition" {
  service_name    = "zpool-12345"
  partition_name  = "my-partition"
  ip              = "123.123.123.123/32"
  type            = "readwrite"
  acl_description = "Description of the ACL"
}
