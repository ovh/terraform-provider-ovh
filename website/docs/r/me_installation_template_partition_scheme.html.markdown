---
subcategory : "Dedicated Server"
---

# ovh_me_installation_template_partition_scheme

Use this resource to create partition scheme for a custom installation template available for dedicated servers.

## Example Usage

```hcl
resource "ovh_me_installation_template" "mytemplate" {
  base_template_name = "debian12_64"
  template_name      = "mytemplate"
}

resource "ovh_me_installation_template_partition_scheme" "scheme" {
  template_name = ovh_me_installation_template.mytemplate.template_name
  name          = "myscheme"
  priority      = 1
}
```

## Argument Reference

* `template_name`: (Required) The template name of the partition scheme.
* `name`: (Required) (Required) This partition scheme name.
* `priority`: on a reinstall, if a partitioning scheme is not specified, the one with the higher priority will be used by default, among all the compatible partitioning schemes (given the underlying hardware specifications).


## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id`: a fake id associated with this partition scheme formatted as follow: `template_name/name`

## Import

The resource can be imported using the `template_name`, `name` of the cluster, separated by "/" E.g.,

```bash
$ terraform import ovh_me_installation_template_partition_scheme.scheme template_name/name
```
