---
layout: "ovh"
page_title: "OVH: iploadbalancing_tcp_farm"
sidebar_current: "docs-ovh-resource-iploadbalancing-tcp-farm-x"
description: |-
  Creates a backend server group (farm).
---

# ovh\_iploadbalancing\_tcp_farm

Creates a backend server group (farm) to be used by loadbalancing frontend(s)

## Example Usage

```
data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
   state       = "ok"  
}

resource "ovh_iploadbalancing_tcp_farm" "farmname" {
  service_name = "${data.ovh_iploadbalancing.lb.id}"
  display_name = "ingress-8080-gra"
  zone = "GRA"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `balance` - Load balancing algorithm. `roundrobin` if null (`first`, `leastconn`, `roundrobin`, `source`)
* `display_name` - Readable label for loadbalancer farm
* `port` - Port attached to your farm ([1..49151]). Inherited from frontend if null
* `stickiness` - 	Stickiness type. No stickiness if null (`sourceIp`)
* `vrack_network_id` - Internal Load Balancer identifier of the vRack private network to attach to your farm, mandatory when your Load Balancer is attached to a vRack
* `zone` - (Required) Zone where the farm will be defined (ie. `GRA`, `BHS` also supports `ALL`)
* `probe` - define a backend healthcheck probe
  * `type` - (Required) Valid values : `http`, `internal`, `mysql`, `oko`, `pgsql`, `smtp`, `tcp`
  * `interval` - probe interval, Value between 30 and 3600 seconds, default 30
  * `match` - What to mach `pattern` against (`contains`, `default`, `internal`, `matches`, `status`)
  * `port` - Port for backends to recieve traffic on.
  * `negate` - Negate probe result
  * `pattern` - Pattern to match against `match`
  * `force_ssl` - Force use of SSL (TLS)
  * `url` - URL for HTTP probe type.
  * `method` - HTTP probe method (`GET`, `HEAD`, `OPTIONS`, `internal`)

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `balance` - See Argument Reference above.
* `display_name` - See Argument Reference above.
* `port` - See Argument Reference above.
* `stickiness` - See Argument Reference above.
* `vrack_network_id` - See Argument Reference above.
* `zone` - See Argument Reference above.
* `probe` - See Argument Reference above.
  * `type` - See Argument Reference above.
  * `interval` - See Argument Reference above.
  * `match` - See Argument Reference above.
  * `port` - See Argument Reference above.
  * `negate` - See Argument Reference above.
  * `pattern` - See Argument Reference above.
  * `force_ssl` - See Argument Reference above.
  * `url` - See Argument Reference above.
  * `method` - See Argument Reference above.
