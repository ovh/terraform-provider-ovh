---
subcategory : "Dedicated Server"
---

# ovh_dedicated_server_bringyourownimage

Install the image provided by an URL on your Dedicated Server.

## Example Usage

```hcl
data ovh_dedicated_server "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
  boot_type    = "harddisk"
}

resource ovh_dedicated_server_bringyourownimage "os_install" {
  service_name = data.ovh_dedicated_server.server.name
  url = "https://url_to_your_server/your_image.qcow2"
  type = "qcow2"
  config_drive {
    enable = true
    hostname = "my-server"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) The service_name of your dedicated server.
* `url` - (Required) The URL of the OS image you wish to install.
* `check_sum` - The checksum of your image for validation.
* `check_sum_type` - The type of the checksum used for validation (either `md5`, `sha1`, `sha256`, `sha512`).
* `config_drive` - The configuration to provide the OS installation.
  * `enable` - Whether to enable configuration through the config_drive or not.
  * `hostname` - The hostname to set on your dedicated server.
  * `ssh_key` - The SSH Key to install on your dedicated server.
  * `user_data` - Custom data to provide to the installation.
  * `user_metadatas` - Custom metadata to provide to the installation.
* `description` - Description of the image to install.
* `disk_group_id` - Disk group id to process install on (only available for some templates).
* `http_headers` - Headers to provide when querying the image.
* `type` - Type of the image to install (either `qcow2` or `raw`).

## Attributes Reference

The following attributes are exported:

* `id` - The task id related to the image setup.
* `comment` - Details of this task. (should be `Install asked`)
* `done_date` - Completion date in RFC3339 format.
* `function` - Function name (should be `hardInstall`).
* `last_update` - Last update in RFC3339 format.
* `start_date` - Task creation date in RFC3339 format.
* `status` - Task status (should be `done`)
