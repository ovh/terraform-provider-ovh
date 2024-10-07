---
subcategory : "Public Cloud Network"
---

# ovh_cloud_project_network_private_subnet_v2

Creates a subnet in a private network of a public cloud region.

## Example Usage

```hcl
resource "ovh_cloud_project_network_private_subnet_v2" "subnet" {
  service_name      = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  network_id        = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  name              = "my_private_subnet"
  region            = "XXX1"
  dns_nameservers   = ["1.1.1.1"]
  cidr              = "192.168.168.0/24"
  dhcp              = true
  enable_gateway_ip = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The id of the public cloud project. If omitted,
    the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. 

* `network_id` - (Required) The id of the network.
   Changing this forces a new resource to be created.

* `dhcp` - (Optional) Enable DHCP.
   Changing this forces a new resource to be created. Defaults to true.

* `cidr` - (Required) IP range of the subnet
   Changing this value recreates the subnet.

* `name` - (Required) Name of the subnet
   Changing this value recreates the subnet.

* `region` - (Required) The region in which the network subnet will be created.
   Ex.: "GRA1". Changing this value recreates the resource.

* `enable_gateway_ip` - Set to true if you want to set a default gateway IP.
   Changing this value recreates the resource. Defaults to true.

* `dns_nameservers` - DNS nameservers used by DHCP
   Changing this value recreates the resource. Defaults to OVH default DNS nameserver.

* `host_routes` - List of custom host routes.
   Changing this value recreates the resource.

* `allocation_pools` - List of IP allocation pools
   Changing this value recreates the resource.

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `cidr` - See Argument Reference above.
* `dhcp` - See Argument Reference above.
* `region` - See Argument Reference above.
* `gateway_ip` - See Argument Reference above.
* `enable_gateway_ip` - See Argument Reference above.
* `dns_nameservers` - See Argument Reference above.
* `host_routes` - See Argument Reference above.
* `allocation_pools` - See Argument Reference above.

## Import

Subnet in a private network of a public cloud project can be imported using the `service_name`, `region`, `network_id` and `subnet_id`, separated by "/" E.g.,

```bash
$ terraform import ovh_cloud_project_network_private_subnet_v2.mysubnet 5ceb661434891538b54a4f2c66fc4b746e/BHS5/25807101-8aaa-4ea5-b507-61f0d661b101/0f0b73a4-403b-45e4-86d0-b438f1291909
```
