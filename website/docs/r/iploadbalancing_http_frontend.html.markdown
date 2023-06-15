---
subcategory : "Load Balancer (IPLB)"
---

# ovh\_iploadbalancing\_http_frontend

Creates a backend HTTP server group (frontend) to be used by loadbalancing frontend(s)

## Example Usage

```hcl
data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_http_farm" "farm80" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name = "ingress-8080-gra"
  zone         = "all"
  port         = 80
}

resource "ovh_iploadbalancing_http_frontend" "testfrontend" {
  service_name    = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name    = "ingress-8080-gra"
  zone            = "all"
  port            = "80,443"
  default_farm_id = "${ovh_iploadbalancing_http_farm.farm80.id}"
}
```

## Example Usage with HTTP header

```
data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_http_farm" "farm80" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name = "ingress-8080-gra"
  zone         = "all"
  port         = 80
}

resource "ovh_iploadbalancing_http_frontend" "testfrontend" {
  service_name    = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name    = "ingress-8080-gra"
  zone            = "all"
  port            = "80,443"
  default_farm_id = "${ovh_iploadbalancing_http_farm.farm80.id}"
  http_header     = ["X-Ip-Header %%ci", "X-Port-Header %%cp"]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `display_name` - Human readable name for your frontend, this field is for you
* `port` - Port(s) attached to your frontend. Supports single port (numerical value),
   range (2 dash-delimited increasing ports) and comma-separated list of 'single port'
   and/or 'range'. Each port must be in the [1;49151] range
* `zone` - (Required) Zone where the frontend will be defined (ie. `gra`, `bhs` also supports `all`)
* `allowed_source` - Restrict IP Load Balancing access to these ip block. No restriction if null. List of IP blocks.
* `dedicated_ipfo` - Only attach frontend on these ip. No restriction if null. List of Ip blocks.
* `default_farm_id` - Default TCP Farm of your frontend
* `default_ssl_id` - Default ssl served to your customer
* `disabled` - Disable your frontend. Default: 'false'
* `ssl` - SSL deciphering. Default: 'false'
* `hsts` - HTTP Strict Transport Security. Default: 'false'
* `redirect_location` - Redirection HTTP'
* `http_header` - HTTP headers to add to the frontend. List of string.

## Attributes Reference

The following attributes are exported:

* `id` - Id of your frontend
* `service_name` - See Argument Reference above.
* `port` - See Argument Reference above.
* `zone` - See Argument Reference above.
* `http_header` - See Argument Reference above.
* `display_name` - See Argument Reference above.
* `allowed_source` - See Argument Reference above.
* `dedicated_ipfo` - See Argument Reference above.
* `default_farm_id` - See Argument Reference above.
* `default_ssl_id` - See Argument Reference above.
* `disabled` - See Argument Reference above.
* `ssl` - See Argument Reference above.
* `hsts` - See Argument Reference above.

## Import 

HTTP frontend can be imported using the following format `service_name` and the `id` of the frontend separated by "/" e.g.

```bash
$ terraform import ovh_iploadbalancing_http_frontend.testfrontend service_name/http_frontend_id
```