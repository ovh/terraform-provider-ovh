---
subcategory : "Account Management"
---

# ovh_me_ssh_key

Creates an SSH Key.

-> __NOTE__ This resource will be removed in next release, please use userMetada in a server install task instead.
```hcl
resource "ovh_dedicated_server_install_task" "server_install" {
  service_name      = "nsxxxxxxx.ip-xx-xx-xx.eu"
  template_name     = ovh_me_installation_template.debian.template_name
  bootid_on_destroy = data.ovh_dedicated_server_boots.rescue.result[0]
  details {
      custom_hostname = "mytest"
  }
  user_metadata {
    key = "sshKey"
    value = "ssh-ed25519 AAAAC3..."
  }
}
```

-> __NOTE__ An SSH key in OVH provider cannot be currently used with Public Cloud instances through Terraform. We advise to use [OpenStack provider](https://registry.terraform.io/providers/terraform-provider-openstack/openstack/latest) to manage Public Cloud instances. Hence, if you need to associate an SSH key to a Public Cloud instance, you need to use [openstack_compute_keypair_v2](https://registry.terraform.io/providers/terraform-provider-openstack/openstack/latest/docs/resources/compute_keypair_v2) resource.

## Example Usage

-> __NOTE__ This resource will be removed in next release, please use userMetada instead.

```hcl
resource "ovh_me_ssh_key" "mykey" {
  key_name = "mykey"
  key      = "ssh-ed25519 AAAAC3..."
}
```

## Argument Reference

-> __NOTE__ This resource will be removed in next release, please use userMetada instead.

The following arguments are supported:

* `key_name` - (Required) The friendly name of this SSH key.

* `key` - (Required) The content of the public key in the form "ssh-algo content", e.g. "ssh-ed25519 AAAAC3...".

* `default` - True when this public SSH key is used for rescue mode and reinstallations.

## Attributes Reference

-> __NOTE__ This resource will be removed in next release, please use userMetada instead.

The following attributes are exported:

* `key_name` - See Argument Reference above.
* `key` - See Argument Reference above.
* `default` - See Argument Reference above.
