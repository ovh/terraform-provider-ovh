---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_udp_farm

Creates a backend server group (farm) to be used by loadbalancing frontend(s)

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `display_name` - Readable label for loadbalancer farm
* `port` - Port attached to your farm ([1..49151]). Inherited from frontend if null
* `vrack_network_id` - Internal Load Balancer identifier of the vRack private network to attach to your farm, mandatory when your Load Balancer is attached to a vRack
* `zone` - (Required) Zone where the farm will be defined (ie. `gra`, `bhs` also supports `all`)

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `farm_id` - Id of your farm.
* `display_name` - See Argument Reference above.
* `port` - See Argument Reference above.
* `vrack_network_id` - See Argument Reference above.
* `zone` - See Argument Reference above.

## Import

UDP Farm can be imported using the following format `service_name` and the `id` of the farm, separated by "/" e.g.

```bash
$ terraform import ovh_iploadbalancing_udp_farm.farmname service_name/farm_id
```