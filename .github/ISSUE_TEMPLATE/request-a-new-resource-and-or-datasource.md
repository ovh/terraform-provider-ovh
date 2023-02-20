---
name: Request a New Resource and/or Datasource
about: Request an entirely new resource and/or data source to add to the provider.
title: "[NEW]"
labels: ''
assignees: ''

---

Please take a look to the list of [existing resources and datasources](https://registry.terraform.io/providers/ovh/ovh/latest/docs) before creating a request.

If you're looking for a change to be made to an existing resource or data source, consider submitting either the "Request a Feature" or Report a Bug" forms instead.

### Title

Please update the title to match what you're requesting, e.g.:

[New Resource]: - for new resource requests
[New Data Source]: - for new datasource requests

### Description

Please describe the new resource or datasource you want and why.

### Requested Resource(s) and/or Data Source(s)

Please list any new resource(s) and/or data source(s). The naming format is ovh_<service>_<resource_name>, e.g., ovh_cloud_project_database_user.

The naming should reflect an [OVHcloud API](https://api.ovh.com/) endpoint.

* ovh_xx_xx

### Potential Terraform Configuration

If this request was implemented, what might the Terraform configuration look like? Similar to above, a best guess is helpful, even if you're unsure of exactly what the end result will look like.
