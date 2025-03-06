---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_udp_farm_server

Creates a backend server entry linked to loadbalancing group (farm)

## Example Usage

```terraform
data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_udp_farm" "farm_name" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name = "ingress-8080-gra"
  zone         = "gra"
  port         = 80
}

resource "ovh_iploadbalancing_udp_farm_server" "backend" {
  service_name           = "${data.ovh_iploadbalancing.lb.service_name}"
  farm_id                = "${ovh_iploadbalancing_udp_farm.farm_name.farm_id}"
  display_name           = "mybackend"
  address                = "4.5.6.7"
  status                 = "active"
  port                   = 80
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `farm_id` - ID of the farm this server is attached to
* `address` - Address of the backend server (IP from either internal or OVHcloud network)
* `display_name` - Label for the server
* `port` - Port that backend will respond on
* `status` - backend status - `active` or `inactive`

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `farm_id` - See Argument Reference above.
* `backend_id` - Synonym for `farm_id`.
* `address` - See Argument Reference above.
* `display_name` - See Argument Reference above.
* `port` - See Argument Reference above.
* `server_id` - Id of your server.
* `status` - See Argument Reference above.

## Import

UDP farm server can be imported using the following format `service_name`, the `id` of the farm and the `id` of the server separated by "/" e.g.

```bash
$ terraform import ovh_iploadbalancing_udp_farm_server.backend service_name/farm_id/server_id
```
