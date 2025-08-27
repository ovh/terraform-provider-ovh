resource "ovh_cloud_project_volume" "volume" {
   region_name  = "xxx"
   service_name = "yyyyy"
   description  = "Terraform volume"
   name         = "terrformName"
   size         = 15
   type         = "classic"
}
