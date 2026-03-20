---
subcategory : "Cloud"
---

# ovh_cloud_security_group

Creates a security group in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_security_group" "sg" {
  service_name = "xxxxxxxxxx"
  region       = "GRA1"
  name         = "my-security-group"
  description  = "Allow SSH and HTTPS"

  rule {
    direction        = "INGRESS"
    ethernet_type       = "IPV4"
    protocol         = "TCP"
    port_range_min   = 22
    port_range_max   = 22
    remote_ip_prefix = "0.0.0.0/0"
    description      = "SSH"
  }

  rule {
    direction        = "INGRESS"
    ethernet_type       = "IPV4"
    protocol         = "TCP"
    port_range_min   = 443
    port_range_max   = 443
    remote_ip_prefix = "0.0.0.0/0"
    description      = "HTTPS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project. **Changing this value recreates the resource.**
* `region` - (Required) Region where the security group will be created. **Changing this value recreates the resource.**
* `name` - (Required) Name of the security group.
* `description` - (Optional) Description of the security group.
* `rule` - (Optional) List of security group rules. Each rule supports:
  * `direction` - (Required) Direction of the rule (`INGRESS` or `EGRESS`).
  * `ethernet_type` - (Required) Ether type (`IPV4` or `IPV6`).
  * `protocol` - (Optional) Protocol (`TCP`, `UDP`, `ICMP`, etc.).
  * `port_range_min` - (Optional) Minimum port number.
  * `port_range_max` - (Optional) Maximum port number.
  * `remote_group_id` - (Optional) Remote security group ID.
  * `remote_ip_prefix` - (Optional) Remote IP prefix (CIDR notation).
  * `description` - (Optional) Description of the rule.

## Attributes Reference

The following attributes are exported:

* `id` - Security group ID.
* `checksum` - Computed hash representing the current target specification value.
* `created_at` - Creation date of the security group.
* `updated_at` - Last update date of the security group.
* `resource_status` - Security group readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
* `current_state` - Current state of the security group:
  * `name` - Name of the security group.
  * `description` - Description of the security group.
  * `region` - Region of the security group.
  * `rules` - Current security group rules:
    * `id` - Rule ID.
    * `direction` - Direction of the rule.
    * `ethernet_type` - Ether type.
    * `protocol` - Protocol.
    * `port_range_min` - Minimum port number.
    * `port_range_max` - Maximum port number.
    * `remote_group_id` - Remote security group ID.
    * `remote_ip_prefix` - Remote IP prefix.
    * `description` - Description of the rule.

## Import

A cloud security group can be imported using the `service_name` and `security_group_id`, separated by `/`:

```terraform
import {
  to = ovh_cloud_security_group.sg
  id = "<service_name>/<security_group_id>"
}
```

```bash
$ terraform import ovh_cloud_security_group.sg service_name/security_group_id
```
