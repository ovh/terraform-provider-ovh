---
subcategory : "Load Balancer (IPLB)"
---

# ovh\_iploadbalancing\_tcp_farm\_server

Creates a backend server entry linked to loadbalancing group (farm)

## Example Usage

```hcl
data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_tcp_farm" "farmname" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  port         = 8080
  zone         = "all"
}

resource "ovh_iploadbalancing_tcp_farm_server" "backend" {
  service_name           = "${data.ovh_iploadbalancing.lb.service_name}"
  farm_id                = "${ovh_iploadbalancing_tcp_farm.farmname.id}"
  display_name           = "mybackend"
  address                = "4.5.6.7"
  status                 = "active"
  port                   = 80
  proxy_protocol_version = "v2"
  weight                 = 2
  probe                  = true
  ssl                    = false
  backup                 = true
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `farm_id` - ID of the farm this server is attached to
* `display_name` - Label for the server
* `address` - Address of the backend server (IP from either internal or OVHcloud network)
* `status` - backend status - `active` or `inactive`
* `port` - Port that backend will respond on
* `proxy_protocol_version` - version of the PROXY protocol used to pass origin connection information from loadbalancer to receiving service (`v1`, `v2`, `v2-ssl`, `v2-ssl-cn`)
* `on_marked_down` - enable action when backend marked down. (`shutdown-sessions`)
* `weight` - used in loadbalancing algorithm
* `probe` - defines if backend will be probed to determine health and keep as active in farm if healthy
* `ssl` - is the connection ciphered with SSL (TLS)
* `backup` - is it a backup server used in case of failure of all the non-backup backends

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `farm_id` - See Argument Reference above.
* `display_name` - See Argument Reference above.
* `address` - See Argument Reference above.
* `status` - See Argument Reference above.
* `port` - See Argument Reference above.
* `proxy_protocol_version` - See Argument Reference above.
* `on_marked_down` - See Argument Reference above.
* `weight` - See Argument Reference above.
* `probe` - See Argument Reference above.
* `ssl` - See Argument Reference above.
* `backup` - See Argument Reference above.
* `cookie` - Value of the stickiness cookie used for this backend.

## Import 

TCP farm server can be imported using the following format `service_name`, the `id` of the farm and the `id` of the server separated by "/" e.g.

```bash
$ terraform import ovh_iploadbalancing_tcp_farm_server.backend service_name/farm_id/server_id
```