---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_role

Reference a DBaaS logs role.

## Example Usage

```hcl
resource "ovh_dbaas_logs_role" "ro" {
  service_name     = "ldp-xx-xxxxx"

  name = "Devops - RO"
  description = "Devops - RO"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service name
* `name` -  (Required) The role name
* `description` - (Required) The role description

## Attributes Reference

Id is set to the role Id. In addition, the following attributes are exported:
* `nb_member` -  number of member for the role
* `nb_permission` - number of configured permission for the role

## Import

OVHcloud DBaaS Log Role can be imported using the `service_name` and `role_id` of the role, separated by "/" E.g.,

```bash
$  terraform import ovh_dbaas_logs_role.ro ldp-ra-XX/dc145bc2-eb01-4efe-a802-XXXXXX
```
