---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_update

Update various properties of your Dedicated Server.

~> __WARNING__ `rescue_mail` and `root_device` properties aren't
updated consistently. This is an issue on the OVHcloud API which 
has been reported. Meanwhile, these properties are not mapped
on this terraform resource.

## Example Usage

```hcl
data "ovh_dedicated_server_boots" "rescue" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_type    = "rescue"
  kernel       = "rescue64-pro"
}

resource "ovh_dedicated_server_update" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_id      = data.ovh_dedicated_server_boots.rescue.result[0]
  monitoring   = true
  state        = "ok"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service_name of your dedicated server.
* `boot_id` - boot id of the server
* `boot_script` - boot script of the server
* `monitoring` - Icmp monitoring state
* `state` - error, hacked, hackedBlocked, ok

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `boot_id` - See Argument Reference above.
* `monitoring` - See Argument Reference above.
* `state` - See Argument Reference above.
