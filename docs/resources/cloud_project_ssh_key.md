---
subcategory : "Cloud Project"
---

# ovh_cloud_project_ssh_key (Resource)

**This resource uses a Beta API.** Creates an SSH key in a Public Cloud project. Keys are stored in the database and lazily synced to OpenStack when an instance referencing the key is created.

~> SSH keys are **immutable** — both `name` and `public_key` force resource replacement if changed.

## Example Usage

```terraform
resource "ovh_cloud_project_ssh_key" "my_key" {
  service_name = "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  name         = "my-deploy-key"
  public_key   = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample user@host"
}
```

### Use with an instance

```terraform
resource "ovh_cloud_project_ssh_key" "my_key" {
  service_name = "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  name         = "my-deploy-key"
  public_key   = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExample user@host"
}

resource "ovh_cloud_project_instance" "server" {
  service_name = "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  region       = "GRA11"
  # ...
  ssh_key {
    name = ovh_cloud_project_ssh_key.my_key.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Optional) The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used. **Changing this value recreates the resource.**
* `name` - (Required) SSH key name. Must be unique within the project. Used as the resource identifier. **Changing this value recreates the resource.**
* `public_key` - (Required) SSH public key content (e.g. the contents of `~/.ssh/id_ed25519.pub`). **Changing this value recreates the resource.**

## Attributes Reference

The following attributes are exported:

* `service_name` - The id of the public cloud project.
* `name` - SSH key name.
* `public_key` - SSH public key content.
* `created_at` - Creation date of the SSH key (RFC 3339 format).
* `updated_at` - Last update date of the SSH key (RFC 3339 format).
