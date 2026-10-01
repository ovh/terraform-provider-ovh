resource "ovh_cloud_project_volume" "encrypted_volume" {
   region_name  = "xxx"
   service_name = "yyyyy"
   description  = "Terraform encrypted volume"
   name         = "encryptedVolume"
   size         = 15
   type         = "classic"

   encryption = {
     encrypted = true
     kms = {
       domain_id      = "<okms domain id>"
       service_key_id = "<okms service key id>"
     }
   }
}
