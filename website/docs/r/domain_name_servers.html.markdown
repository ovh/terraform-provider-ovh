---
subcategory : "Domain names"
---

# ovh_domain_name_servers

Use this resource to manage a domain's name servers.

## Example Usage

```hcl
resource "ovh_domain_name_servers" "name_servers" {
  domain = "mydomain.ovh"

  servers {
    host = "dns105.ovh.net"
    ip = "213.251.188.144"
  }

  servers {
    host = "ns105.ovh.net"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) Domain name for which to manage name servers
* `servers` - (Required) Details about a name server
  * `host` - (Required) The server hostname
  * `ip` - (Optional) The server IP


## Attributes Reference

* `domain` - Domain name and resource ID
* `servers` - Details about a name server
  * `host` - The server hostname
  * `ip` - The server IP

## Import

Name servers can be imported using their `domain`. 

Using the following configuration:

```hcl
import {
  to = ovh_domain_name_servers.name_servers
  id = "<domain name>"
}
```

You can then run:

```bash
$ terraform plan -generate-config-out=name_servers.tf
$ terraform apply
```

The file `name_servers.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above.
See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
