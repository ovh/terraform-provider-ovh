---
layout: "ovh"
page_title: "OVH: cloud_project_database_opensearch_user"
sidebar_current: "docs-ovh-resource-cloud-project-database-opensearch-user"
description: |-
  Creates an user for a opensearch cluster associated with a public cloud project.
---

# ovh_cloud_project_database_opensearch_user

Creates an user for a opensearch cluster associated with a public cloud project.

## Example Usage

```hcl
data "ovh_cloud_project_database" "opensearch" {
  service_name  = "XXX"
  engine        = "opensearch"
  id            = "ZZZ"
}

resource "ovh_cloud_project_database_opensearch_user" "user" {
  service_name  = data.ovh_cloud_project_database.opensearch.service_name
  cluster_id    = data.ovh_cloud_project_database.opensearch.id
  acls {
		pattern = "logs_*"
		permission = "read"
	}
	acls {
		pattern = "data_*"
		permission = "deny"
	}
  name          = "johndoe"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `cluster_id` - (Required, Forces new resource) Cluster ID.

* `acls` - (Optional) Acls of the user.
  * `pattern` - (Required) Pattern of the ACL.
  * `permission` - (Required) Permission of the ACL:
    * `admin`
    * `read`
    * `write`
    * `readwrite`
    * `deny`

* `name` - (Required, Forces new resource) Username affected by this acl.

## Attributes Reference

The following attributes are exported:

* `acls` - See Argument Reference above.
  * `pattern` - See Argument Reference above.
  * `permission` - See Argument Reference above.
* `cluster_id` - See Argument Reference above.
* `created_at` - Date of the creation of the user.
* `id` - ID of the user.
* `name` - See Argument Reference above.
* `password` - (Sensitive) Password of the user.
* `service_name` - See Argument Reference above.
* `status` - Current status of the user.

## Import

OVHcloud Managed opensearch clusters users can be imported using the `service_name`, `cluster_id` and `id` of the user, separated by "/" E.g.,

```
$ terraform import ovh_cloud_project_database_opensearch_user.my_user service_name/cluster_id/id