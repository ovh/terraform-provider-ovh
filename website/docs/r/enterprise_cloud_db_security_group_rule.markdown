---
layout: "ovh"
page_title: "OVH: ovh_enterprise_cloud_db_security_group_rule"
sidebar_current: "docs-ovh-resource-enterprise-cloud-db-security-group-rule"
description: |-
  Add a new rule in a Security Group for an Enterprise Cloud DB.
---

# ovh_enterprise_cloud_db_security_group

Add a new rule in an Enterprise Cloud DB Security Group

## Example Usage

```hcl
data "ovh_enterprise_cloud_db" "db" {
	cluster_id = "%s"
}
	
resource "ovh_enterprise_cloud_db_security_group" "sg" {
  cluster_id = data.ovh_enterprise_cloud_db.db.id
  name = "example"
}

resource "ovh_enterprise_cloud_db_security_group_rule" "rule" {
  cluster_id = data.ovh_enterprise_cloud_db.db.id
  security_group_id = ovh_enterprise_cloud_db_security_group.sg.id
  source = "10.0.0.0/8"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The Enterprise Cloud DB ID
* `security_group_id` - (Required) The security group ID
* `source` - (Required) The IPV4 network mask to apply

The following attributes are exported:

* `cluster_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `source` - See Argument Reference above.
