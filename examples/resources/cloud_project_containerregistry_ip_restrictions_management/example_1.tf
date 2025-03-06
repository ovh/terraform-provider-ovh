data "ovh_cloud_project_containerregistry" "registry" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

resource "ovh_cloud_project_containerregistry_ip_restrictions_management" "my_mgt_iprestrictions" {
  service_name = ovh_cloud_project_containerregistry.registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.registry.id

  ip_restrictions = [
    { 
      ip_block = "xxx.xxx.xxx.xxx/xx"
      description = "xxxxxxx"
    }
  ]
}
