---
subcategory: "Gateway"
---

# ovh_cloud_project_gateway

Create a new Gateway for existing subnet in the specified public cloud project.

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The ID of the private network.
* `name` - (Required) The name of the gateway. \*` model` - (Required) The model of the gateway.
* `model` - (Required) The model of the gateway.
* `network_id` - (Required) The ID of the private network.
* `subnet_id` - (Required) The ID of the subnet.

## Attributes Reference
The following attributes are exported:

* `service_name` - See Argument Reference above.
* `name` - See Argument Reference above.
* `model` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `status` - The status of the gateway.
