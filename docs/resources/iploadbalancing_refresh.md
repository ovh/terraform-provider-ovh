---
subcategory : "Load Balancer (IPLB)"
---

# ovh\_iploadbalancing\_refresh

Applies changes from other `ovh_iploadbalancing_*` resources to the production configuration of loadbalancers.

## Example Usage

```terraform
data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_tcp_farm" "farm_name" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  port         = 8080
  zone         = "all"
}

resource "ovh_iploadbalancing_tcp_farm_server" "backend" {
  service_name           = "${data.ovh_iploadbalancing.lb.service_name}"
  farm_id                = "${ovh_iploadbalancing_tcp_farm.farm_name.id}"
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

resource "ovh_iploadbalancing_refresh" "mylb" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  keepers = [
    "${ovh_iploadbalancing_tcp_farm_server.backend.*.address}",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `keepers` - List of values tracked to trigger refresh, used also to form implicit dependencies

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `keepers` - See Argument Reference above.
