---
subcategory : "Additional IP"
---

# ovh_ip_reverse

Provides a OVHcloud IP reverse.

## Example Usage

```terraform
# Set the reverse of an IP
resource "ovh_ip_reverse" "test" {
  readiness_timeout_duration = "1m"
  ip                         = "192.0.2.0/24"
  ip_reverse                 = "192.0.2.1"
  reverse                    = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `ip` - (Required) The IP block to which the IP belongs
* `reverse` - (Required) The value of the reverse
* `ip_reverse` - (Required) The IP to set the reverse of
* `readiness_timeout_duration` - (Optional) The maximum duration that the provider will wait for a successful response (while retrying every 5s). If the record cannot be verified within this timeout, the operation will fail (default value: 60s)

## Attributes Reference

The id is set to the value of ip_reverse.

## Import

The resource can be imported using the `ip`, `ip_reverse` of the address, separated by "|" E.g.,

```bash
$ terraform import ovh_ip_reverse.my_reverse '2001:0db8:c0ff:ee::/64|2001:0db8:c0ff:ee::42'
```
