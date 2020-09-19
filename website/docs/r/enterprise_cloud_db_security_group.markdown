---
layout: "ovh"
page_title: "OVH: ovh_enterprise_cloud_db_security_group"
sidebar_current: "docs-ovh-resource-enterprise-cloud-db-security-group"
description: |-
  Creates a new Security Group for an Enterprise Cloud DB.
---

# ovh_enterprise_cloud_db_security_group

Add a new Security Group in an Enterprise Cloud DB

## Example Usage

```hcl
data "ovh_enterprise_cloud_db" "db" {
	cluster_id = "%s"
}
	
resource "ovh_enterprise_cloud_db_security_group" "sg" {
  cluster_id = data.ovh_enterprise_cloud_db.db.id
  name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The Enterprise Cloud DB ID
* `name` - (Required) The security group name

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `name` - See Argument Reference above.
