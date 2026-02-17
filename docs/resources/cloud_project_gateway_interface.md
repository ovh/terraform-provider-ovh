---
subcategory: "Gateway"
---

# ovh_cloud_project_gateway_interface

Create a new Gateway Interface for existing subnet in the specified public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_network_private" "mypriv" {
  service_name  = "xxxxxxxxxx"
  vlan_id       = "0"
  name          = "mypriv"
  regions       = ["GRA9"]
}

resource "ovh_cloud_project_network_private_subnet" "my_privsub" {
  service_name  = ovh_cloud_project_network_private.mypriv.service_name
  network_id    = ovh_cloud_project_network_private.mypriv.id
  region        = "GRA9"
  start         = "10.0.0.2"
  end           = "10.0.0.8"
  network       = "10.0.0.0/24"
  dhcp          = true
}

resource "ovh_cloud_project_network_private_subnet" "my_other_privsub" {
	service_name  = ovh_cloud_project_network_private.mypriv.service_name
	network_id    = ovh_cloud_project_network_private.mypriv.id
	region        = "GRA9"
	start         = "10.0.1.10"
	end           = "10.0.1.254"
	network       = "10.0.1.0/24"
	dhcp          = true
}

resource "ovh_cloud_project_gateway" "gateway" {
  service_name = ovh_cloud_project_network_private.mypriv.service_name
  name          = "my-gateway"
  model         = "s"
  region        = ovh_cloud_project_network_private_subnet.my_privsub.region
  network_id    = tolist(ovh_cloud_project_network_private.mypriv.regions_attributes[*].openstackid)[0]
  subnet_id     = ovh_cloud_project_network_private_subnet.my_privsub.id
}

resource "ovh_cloud_project_gateway_interface" "interface" {
	service_name = ovh_cloud_project_network_private.mypriv.service_name
	region       = ovh_cloud_project_network_private_subnet.my_other_privsub.region
	id           = ovh_cloud_project_gateway.gateway.id
	subnet_id    = ovh_cloud_project_network_private_subnet.my_other_privsub.id
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) ID of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `region` - (Required) Region of the gateway
* `id` - (Required) ID of the gateway
* `subnet_id` - (Required) ID of the subnet to add

## Attributes Reference

The following attributes are exported:

* `service_name` - ID of the cloud project
* `region` - Region of the gateway
* `id` - ID of the gateway
* `subnet_id` - ID of the subnet to add
* `interface_id` - ID of the interface
* `ip` - IP of the interface
* `network_id` - Network ID of the interface

## Import

A gateway interface can be imported using the `service_name`, `region`, `id` (identifier of the gateway) and `interface_id` properties, separated by a `/`.

```bash
$ terraform import ovh_cloud_project_gateway_interface.interface service_name/region/id/interface_id
```
