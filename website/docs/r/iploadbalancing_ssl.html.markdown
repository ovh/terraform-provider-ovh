---
subcategory : "Load Balancer (IPLB)"
---

# ovh_iploadbalancing_ssl

Creates a new custom SSL certificate on your IP Load Balancing

## Example Usage

```hcl
data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_ssl" "sslname" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  display_name = "test"
  certificate  = "..."
  key          = "..."
  chain        = "..."

  # use this if ssl is configured as frontend default_ssl
  lifecycle {
    create_before_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The internal name of your IP load balancing
* `display_name` - Readable label for loadbalancer ssl
* `certificate` - Certificate
* `chain` - Certificate chain
* `key` - Certificate key

## Attributes Reference

The following attributes are exported:

* `service_name` - See Argument Reference above.
* `display_name` - See Argument Reference above.
* `certificate` - See Argument Reference above.
* `chain` - See Argument Reference above.
* `expire_date` - Expire date of your SSL certificate.
* `fingerprint` - Fingerprint of your SSL certificate.
* `id` - Id of your SSL certificate.
* `san` - Subject Alternative Name of your SSL certificate.
* `serial` - Serial of your SSL certificate (Deprecated, use fingerprint instead !)
* `subject` - Subject of your SSL certificate.
* `type` - Type of your SSL certificate. 'built' for SSL certificates managed by the IP Load Balancing. 'custom' for user manager certificates.

## Import

SSL can be imported using the following format `service_name` and the `id` of the ssl, separated by "/" e.g.

```bash
$ terraform import ovh_iploadbalancing_ssl.sslname service_name/ssl_id
```