resource "ovh_cloud_project_network_private" "mypriv" {
  service_name  = "xxxxxxxxxx"
  vlan_id       = "0"
  name          = "mypriv"
  regions       = ["GRA9"]
}

resource "ovh_cloud_project_network_private_subnet" "myprivsub" {
  service_name  = ovh_cloud_project_network_private.mypriv.service_name
  network_id    = ovh_cloud_project_network_private.mypriv.id
  region        = "GRA9"
  start         = "10.0.0.2"
  end           = "10.0.255.254"
  network       = "10.0.0.0/16"
  dhcp          = true
}

resource "ovh_cloud_project_gateway" "gateway" {
  service_name = ovh_cloud_project_network_private.mypriv.service_name
  name          = "my-gateway"
  model         = "s"
  region        = ovh_cloud_project_network_private_subnet.myprivsub.region
  network_id    = tolist(ovh_cloud_project_network_private.mypriv.regions_attributes[*].openstackid)[0]
  subnet_id     = ovh_cloud_project_network_private_subnet.myprivsub.id
}
