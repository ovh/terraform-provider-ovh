## 0.1.1 (Unreleased)

BACKWARDS INCOMPATIBILITIES / NOTES:

* d/ovh_publiccloud_region: Deprecated fields which don't comply
  with lowercase & underscore convention (`continentCode`, `datacenterLocation`).
  Use `continent_code` and `datacenter_location` instead. [GH-4]

FEATURES

* __New Resource:__ `ovh_domain_zone_record` ([#3](https://github.com/terraform-providers/terraform-provider-ovh/issues/3))

IMPROVEMENTS

* The provider config can now source its credentials from `~/.ovh.conf` ([#10](https://github.com/terraform-providers/terraform-provider-ovh/issues/10))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
