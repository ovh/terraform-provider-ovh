resource "ovh_domain_ds_records" "ds_records" {
  domain = "mydomain.ovh"
  
  ds_records {
      algorithm = "RSASHA1_NSEC3_SHA1"
      flags = "KEY_SIGNING_KEY"
      public_key = "my_base64_encoded_public_key"
      tag = 12345
  }
}
