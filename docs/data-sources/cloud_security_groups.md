---
subcategory: "Instances"
---

# ovh_cloud_security_groups (Data Source)

Use this data source to list the security groups of a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_security_groups" "security_groups" {
  service_name = "<public cloud project ID>"
}

output "security_group_names" {
  value = [for sg in data.ovh_cloud_security_groups.security_groups.security_groups : sg.name]
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name of the resource representing the id of the cloud project.

## Attributes Reference

The following attributes are exported:

* `security_groups` - List of security groups. Each element exports:
  * `id` - Security group ID.
  * `name` - Name of the security group.
  * `description` - Description of the security group.
  * `location` - Location of the security group:
    * `region` - Region.
  * `rule` - List of security group rules:
    * `direction` - Direction of the rule (`INGRESS` or `EGRESS`).
    * `ethernet_type` - Ethernet type (`IPV4` or `IPV6`).
    * `protocol` - Protocol (`TCP`, `UDP`, `ICMP`, etc.).
    * `port_range_min` - Minimum port number.
    * `port_range_max` - Maximum port number.
    * `remote_group_id` - Remote security group ID.
    * `remote_ip_prefix` - Remote IP prefix (CIDR notation).
    * `description` - Description of the rule.
  * `checksum` - Computed hash representing the current target specification value.
  * `created_at` - Creation date of the security group.
  * `updated_at` - Last update date of the security group.
  * `resource_status` - Security group readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`).
  * `current_state` - Current state of the security group:
    * `name` - Name of the security group.
    * `description` - Description of the security group.
    * `location` - Location details:
      * `region` - Region.
    * `rules` - User-specified security group rules with their IDs:
      * `id` - Rule ID.
      * `direction` - Direction of the rule.
      * `ethernet_type` - Ethernet type.
      * `protocol` - Protocol.
      * `port_range_min` - Minimum port number.
      * `port_range_max` - Maximum port number.
      * `remote_group_id` - Remote security group ID.
      * `remote_ip_prefix` - Remote IP prefix.
      * `description` - Description of the rule.
    * `default_rules` - Default egress rules auto-created by OpenStack (same schema as `rules`).
