---
subcategory : "Load Balancer"
---

# ovh_cloud_loadbalancer_listeners (Data Source)

Use this data source to list the listeners of a load balancer in a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_loadbalancer_listeners" "listeners" {
  service_name    = "<public cloud project ID>"
  loadbalancer_id = "<load balancer ID>"
}

output "listener_ports" {
  value = [for l in data.ovh_cloud_loadbalancer_listeners.listeners.listeners : l.protocol_port]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.
* `loadbalancer_id` - (Required) ID of the load balancer.

## Attributes Reference

The following attributes are exported:

* `listeners` - List of listeners. Each element exports:
  * `id` - Listener ID.
  * `name` - Name of the listener.
  * `description` - Description of the listener.
  * `protocol` - Protocol of the listener (`HTTP`, `HTTPS`, `TCP`, `UDP`, ...).
  * `protocol_port` - Port number the listener listens on.
  * `connection_limit` - Maximum number of connections allowed.
  * `allowed_cidrs` - List of CIDRs allowed to access the listener.
  * `timeout_client_data` - Timeout for client data in milliseconds.
  * `timeout_member_data` - Timeout for member data in milliseconds.
  * `timeout_member_connect` - Timeout for member connection in milliseconds.
  * `timeout_tcp_inspect` - Timeout for TCP inspect in milliseconds.
  * `insert_headers` - Headers to insert into requests:
    * `x_forwarded_for` - Insert X-Forwarded-For header.
    * `x_forwarded_port` - Insert X-Forwarded-Port header.
    * `x_forwarded_proto` - Insert X-Forwarded-Proto header.
    * `x_ssl_client_verify` - Insert X-SSL-Client-Verify header.
    * `x_ssl_client_has_cert` - Insert X-SSL-Client-Has-Cert header.
    * `x_ssl_client_dn` - Insert X-SSL-Client-DN header.
  * `default_tls_container_ref` - Reference to the default TLS container.
  * `default_pool_id` - ID of the default pool for this listener.
  * `sni_container_refs` - List of SNI container references.
  * `tls_versions` - List of TLS versions allowed.
  * `checksum` - Computed hash representing the current target specification value.
  * `created_at` - Creation date of the listener.
  * `updated_at` - Last update date of the listener.
  * `resource_status` - Listener readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
  * `current_state` - Current state of the listener:
    * `name` - Listener name.
    * `description` - Listener description.
    * `protocol` - Listener protocol.
    * `protocol_port` - Port number the listener listens on.
    * `connection_limit` - Maximum number of connections allowed.
    * `timeout_client_data` - Timeout for client data in milliseconds.
    * `timeout_member_data` - Timeout for member data in milliseconds.
    * `timeout_member_connect` - Timeout for member connection in milliseconds.
    * `timeout_tcp_inspect` - Timeout for TCP inspect in milliseconds.
    * `operating_status` - Operating status of the listener.
    * `provisioning_status` - Provisioning status of the listener.
    * `default_tls_container_ref` - Reference to the default TLS container.
    * `default_pool_id` - ID of the default pool for this listener.
    * `region` - Region.
    * `availability_zone` - Availability zone.
    * `insert_headers` - Headers inserted into requests (same schema as above).
    * `allowed_cidrs` - List of CIDRs allowed to access the listener.
    * `sni_container_refs` - List of SNI container references.
    * `tls_versions` - List of TLS versions allowed.
