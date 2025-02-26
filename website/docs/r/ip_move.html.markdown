---
subcategory : "Additional IP"
---

# ovh_ip_move

Moves a given IP to a different service, or inversely, parks it if empty service is given

## Example Usage

## Move IP `1.2.3.4` to service loadbalancer-XXXXX

```hcl
resource "ovh_ip_move" "move_ip_to_load_balancer_xxxxx" {
  ip = "1.2.3.4"
  routed_to {
    service_name = "loadbalancer-XXXXX"
  }
}
```

## Park IP/Detach IP `1.2.3.4` from any service

```hcl
resource "ovh_ip_move" "park_ip" {
  ip = "1.2.3.4"
  routed_to {
    service_name = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `ip` - (Required) IP block that we want to attach to a different service
* `routed_to` - (Required) Service to route the IP to. If null, the IP will be [parked](https://api.ovh.com/console/#/ip/%7Bip%7D/park~POST)
  instead of [moved](https://api.ovh.com/console/#/ip/%7Bip%7D/move~POST)
  * `service_name` - (Required) Name of the service to route the IP to. IP will be parked if this value is an empty string

## Attributes Reference

Attributes are mostly the same as for [ovh_ip_service](https://registry.terraform.io/providers/ovh/ovh/latest/docs/resources/ip_service#attributes-reference):

* `can_be_terminated` - Whether IP service can be terminated
* `country` - Country
* `description` - Description attached to the IP
* `ip` - IP block
* `organisation_id` - IP block organisation Id
* `routed_to` - Routage information
  * `service_name` - Service where ip is routed to
* `service_name`: Service name in the form of `ip-<part-1>.<part-2>.<part-3>.<part-4>`
* `type` - Possible values for ip type
* `task_status` - Status field of the current IP task that is in charge of changing the service the IP is attached to
* `task_start_date` - Starting date and time field of the current IP task that is in charge of changing the service the IP is attached to

## Import

The resource can be imported using the `ip` field, e.g.,

```bash
$ terraform import ovh_ip_move.mv '1.2.3.4/32'
```
