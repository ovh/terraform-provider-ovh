---
layout: "ovh"
page_title: "OVH: ovh_domain_zone_redirection"
sidebar_current: "docs-ovh-resource-domain-zone-redirection"
description: |-
  Provides a OVH domain zone resource.
---

# ovh_domain_zone_redirection

Provides a OVH domain zone redirection.

## Example Usage

```hcl
# Add a redirection to a sub-domain
resource "ovh_domain_zone_redirection" "test" {
    zone = "testdemo.ovh"
    subdomain = "test"
    type = "visiblePermanent"
    target = "http://www.ovh"
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The domain to add the redirection to
* `subdomain` - (Optional) The name of the redirection
* `target` - (Required) The value of the redirection
* `type` - (Required) The type of the redirection, with values:
  * `visible` -> Redirection by http code 302
  * `visiblePermanent` -> Redirection by http code 301
  * `invisible` -> Redirection by html frame
* `description` - (Optional) A description of this redirection
* `keywords` - (Optional) Keywords to describe this redirection
* `title` - (Optional) Title of this redirection

## Attributes Reference

The following attributes are exported:

* `id` - The redirection ID
* `zone` - The domain to add the redirection to
* `subDomain` - The name of the redirection
* `target` - The value of the redirection
* `type` - The type of the redirection
* `description` - The description of the redirection
* `keywords` - Keywords  of the redirection
* `title` - The title of the redirection
