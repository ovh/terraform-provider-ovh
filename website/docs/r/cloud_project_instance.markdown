---
subcategory : "Cloud Project"
---

# ovh_cloud_project_instance
**This resource uses a Beta API**
Creates an instance associated with a public cloud project.

## Example Usage

Create a instance.

```hcl
resource "ovh_cloud_project_instance" "instance" {
  service_name  = "XXX"
    region = "RRRR"
    billing_period = "hourly"
    boot_from {
        image_id = "UUID"
    }
    flavor {
        flavor_id = "UUID"
    }
    name = "sshkeyname"
    ssh_key {
        name = "sshname"
    }
    network {
        public = true
    }  
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, Forces new resource) The id of the public cloud project. If omitted,
  the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.
* `region` - (Required, Forces new resource) Instance region.
* `billing_period` - (Required, Forces new resource) Billing period - hourly or monthly
* `network` - (Required, Forces new resource) Create network interfaces
  * `public` - (Optional, Forces new resource) Set the new instance as public boolean
* `flavor` - (Required, Forces new resource) Flavor information
  * `flavor_id` - (Required, Forces new resource) Flavor ID
* `boot_from` - (Required, Forces new resource) Boot the instance from an image or a volume
  * `image_id` - (Mandatory only if volume_id not used, Forces new resource) Instance image id
  * `volume_id` - (Mandatory only if image_id not used, Forces new resource) Instance volume id
* `group`- (Optional, Forces new resource) Start instance in group
  * `group_id` - (Optional, Forces new resource) Group id
* `name` - (Required, Forces new resource) Instance name
* `ssh_key` - (Mandatory only if ssh_key_create not used, Forces new resource) Existing SSH Keypair
  * `name` - (Optional, Forces new resource) SSH Keypair name
* `ssh_key_create` - (Mandatory only if ssh_key not used, Forces new resource) Unix cron pattern
  * `name` - (Optional, Forces new resource) SSH Keypair name
  * `public_key` - (Optional, Forces new resource) SSH Public key
* `user_data`- (Optional, Forces new resource) Configuration information or scripts to use upon launch
* `auto_backup` - (Optional, Forces new resource) Create an autobackup workflow after instance start up.
  * `cron` - (Optional, Forces new resource) Unix cron pattern
  * `rotation` - (Optional, Forces new resource) Number of backup to keep

## Attributes Reference

The following attributes are exported:

* `addresses` - Instance IP addresses
  * `ip` - IP address
  * `version` - IP version
* `attached_volumes` - Volumes attached to the instance
  * `id` - Volume Id
* `flavor_id` - Flavor id
* `flavor_name` - Flavor name
* `id` - Instance id
* `image_id` - Image id
* `task_state` - Instance task state
