resource "ovh_hosting_privatedatabase_whitelist" "ip" {
  service_name = "XXXXXX"
  ip           = "1.2.3.4"
  name         = "A name for your IP address"
  service      = true
  sftp         = true
}
