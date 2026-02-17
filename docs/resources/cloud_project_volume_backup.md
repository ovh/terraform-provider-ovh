---
subcategory : "Cloud Project"
---

# ovh_cloud_project_volume_backup

Manage backups for the given volume in a public cloud project.

## Example Usage

```terraform
resource "ovh_cloud_project_volume_backup" "backup" {
  service_name = "<public cloud project ID>"
  region_name  = "GRA9"
  volume_id    = "<volume ID>"
  name         = "ExampleBackup"
}
```

## Schema

### Required

- `region_name` (String) Region name
- `volume_id` (String) ID of the volume to backup

### Optional

- `name` (String) name of the backup
- `service_name` (String) Service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

### Read-Only

- `creation_date` (String) Creation date of the backup
- `id` (String) Volume backup id
- `region` (String) Volume backup region
- `size` (Number) Size of the backup in GiB
- `status` (String) Staus of the backup

## Import

A volume backup in a public cloud project can be imported using the `service_name`, `region_name` and `id` attributes. Using the following configuration:

```terraform
import {
  id = "<service_name>/<region_name>/<id>"
  to = ovh_cloud_project_volume_backup.backup
}
```

You can then run:

```bash
$ terraform plan -generate-config-out=backup.tf
$ terraform apply
```

The file `backup.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above. See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
