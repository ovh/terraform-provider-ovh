---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer_listener

Creates a listener on a load balancer in a public cloud project.

## Example Usage

### Basic HTTP Listener

```terraform
resource "ovh_cloud_loadbalancer_listener" "http" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  name            = "http-listener"
  protocol        = "HTTP"
  protocol_port   = 80
}
```

### HTTPS Listener with TLS and Insert Headers

```terraform
resource "ovh_cloud_loadbalancer_listener" "https" {
  service_name    = "xxxxxxxxxx"
  loadbalancer_id = ovh_cloud_loadbalancer.lb.id
  name            = "https-listener"
  protocol        = "TERMINATED_HTTPS"
  protocol_port   = 443
  description     = "HTTPS listener with TLS termination"

  default_tls_container_ref = "https://key-manager.cloud.ovh.net/v1/containers/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  tls_versions              = ["TLSv1.2", "TLSv1.3"]

  insert_headers {
    x_forwarded_for   = true
    x_forwarded_port  = true
    x_forwarded_proto = true
  }

  timeout_client_data    = 50000
  timeout_member_data    = 50000
  timeout_member_connect = 5000
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `loadbalancer_id` - (Required) ID of the load balancer. **Changing this value recreates the resource.**
* `name` - (Required) Name of the listener.
* `protocol` - (Required) Protocol of the listener (`HTTP`, `HTTPS`, `SCTP`, `TCP`, `TERMINATED_HTTPS`, `UDP`). **Changing this value recreates the resource.**
* `protocol_port` - (Required) Port number the listener listens on. **Changing this value recreates the resource.**
* `description` - (Optional) Description of the listener.
* `connection_limit` - (Optional) Maximum number of connections allowed.
* `allowed_cidrs` - (Optional) List of CIDRs allowed to access the listener.
* `timeout_client_data` - (Optional) Timeout for client data in milliseconds.
* `timeout_member_data` - (Optional) Timeout for member data in milliseconds.
* `timeout_member_connect` - (Optional) Timeout for member connection in milliseconds.
* `timeout_tcp_inspect` - (Optional) Timeout for TCP inspect in milliseconds.
* `insert_headers` - (Optional) Headers to insert into requests:
  * `x_forwarded_for` - (Optional) Insert X-Forwarded-For header.
  * `x_forwarded_port` - (Optional) Insert X-Forwarded-Port header.
  * `x_forwarded_proto` - (Optional) Insert X-Forwarded-Proto header.
  * `x_ssl_client_verify` - (Optional) Insert X-SSL-Client-Verify header.
  * `x_ssl_client_has_cert` - (Optional) Insert X-SSL-Client-Has-Cert header.
  * `x_ssl_client_dn` - (Optional) Insert X-SSL-Client-DN header.
* `default_tls_container_ref` - (Optional) Reference to the default TLS container.
* `sni_container_refs` - (Optional) List of SNI container references.
* `tls_versions` - (Optional) List of TLS versions allowed.

## Attributes Reference

The following attributes are exported:

* `id` - Listener ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the listener.
* `updated_at` - Last update date of the listener.
* `resource_status` - Listener readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the listener:
  * `name` - Listener name.
  * `description` - Listener description.
  * `protocol` - Listener protocol.
  * `protocol_port` - Port number.
  * `connection_limit` - Maximum number of connections.
  * `timeout_client_data` - Timeout for client data in milliseconds.
  * `timeout_member_data` - Timeout for member data in milliseconds.
  * `timeout_member_connect` - Timeout for member connection in milliseconds.
  * `timeout_tcp_inspect` - Timeout for TCP inspect in milliseconds.
  * `operating_status` - Operating status of the listener.
  * `provisioning_status` - Provisioning status of the listener.
  * `default_tls_container_ref` - Reference to the default TLS container.
  * `region` - Region.
  * `availability_zone` - Availability zone.
  * `insert_headers` - Headers inserted into requests:
    * `x_forwarded_for` - X-Forwarded-For header enabled.
    * `x_forwarded_port` - X-Forwarded-Port header enabled.
    * `x_forwarded_proto` - X-Forwarded-Proto header enabled.
    * `x_ssl_client_verify` - X-SSL-Client-Verify header enabled.
    * `x_ssl_client_has_cert` - X-SSL-Client-Has-Cert header enabled.
    * `x_ssl_client_dn` - X-SSL-Client-DN header enabled.
  * `allowed_cidrs` - List of allowed CIDRs.
  * `sni_container_refs` - List of SNI container references.
  * `tls_versions` - List of TLS versions.

## Import

A cloud load balancer listener can be imported using the `service_name`, `loadbalancer_id`, and `listener_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_loadbalancer_listener.http
  id = "<service_name>/<loadbalancer_id>/<listener_id>"
}
```

```bash
$ terraform import ovh_cloud_loadbalancer_listener.http service_name/loadbalancer_id/listener_id
```
