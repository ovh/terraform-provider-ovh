data "ovh_cloud_project_gateway" "gateway" {
    service_name = "<public cloud project ID>"
    region       = "GRA11"
    id           = "<gateway ID>"
}