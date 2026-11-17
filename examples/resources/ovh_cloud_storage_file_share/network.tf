resource "ovh_cloud_storage_file_share" "share_on_subnet" {
  service_name = "<public cloud project ID>"
  name         = "my-share-on-subnet"
  size         = 150
  region       = "GRA11"
  protocol     = "NFS"
  share_type   = "STANDARD_1AZ"

  network = {
    subnet_id = "<subnet id>"
  }
}
