## 0.4.0 (Unreleased)

NOTES/DEPRECATIONS:

* `*/ovh_publiccloud_*`: Deprecate `publiccloud` data sources & resources (to issue warnings when used) [GH-81]

FEATURES:

* __New Resource:__ `ovh_iploadbalancing_tcp_frontend` [GH-58]

IMPROVEMENTS:

* provider: Enable request/response logging in `>=DEBUG` mode [GH-77]
* provider: Make homedir detection more robust [GH-82]

BUG FIXES:

* resource/ovh_domain_zone_record: Attempt retries to avoid errors caused by eventual consistency after creation/update/deletion [GH-77]
* resource/ovh_domain_zone_record: Make fieldtype non-updatable [GH-84]
* resource/ovh_domain_zone_redirection: Return errors from refreshing after creation/update/deletion [GH-77]

## 0.3.0 (July 11, 2018)

DEPRECATIONS / NOTES:

Resources and datasources names now reflects the OVH API endpoints. As such,
resources & datasources that doesn't comply are now deprecated and will be removed
in next release.

* data-source/ovh_publiccloud_region: Deprecated in favor of data-source/ovh_cloud_region ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* data-source/ovh_publiccloud_regions: Deprecated in favor of data-source/ovh_cloud_regions ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* resource/ovh_publiccloud_private_network: Deprecated in favor of resource/ovh_cloud_network_private ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* resource/ovh_publiccloud_private_network_subnet: Deprecated in favor of resource/ovh_cloud_network_private_subnet ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* resource/ovh_publiccloud_user: Deprecated in favor of resource/ovh_cloud_user ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* resource/ovh_vrack_publiccloud_attachment: Deprecated in favor of resource/ovh_vrack_cloudproject ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))

FEATURES

* __New Datasource:__ `ovh_cloud_region` ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* __New Datasource:__ `ovh_cloud_regions` ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* __New Datasource:__ `ovh_domain_zone` ([#39](https://github.com/terraform-providers/terraform-provider-ovh/pull/39))
* __New Datasource:__ `ovh_iploadbalancing` ([#40](https://github.com/terraform-providers/terraform-provider-ovh/pull/40))
* __New Datasource:__ `ovh_me_paymentmean_creditcard` ([#34](https://github.com/terraform-providers/terraform-provider-ovh/pull/34),[#52](https://github.com/terraform-providers/terraform-provider-ovh/pull/52))
* __New Datasource:__ `ovh_me_paymentmean_bankaccount` ([#34](https://github.com/terraform-providers/terraform-provider-ovh/pull/34),[#52](https://github.com/terraform-providers/terraform-provider-ovh/pull/52))
* __New Resource:__ `ovh_iploadbalancing_tcp_farm` ([#32](https://github.com/terraform-providers/terraform-provider-ovh/pull/32))
* __New Resource:__ `ovh_iploadbalancing_tcp_farm_server` ([#33](https://github.com/terraform-providers/terraform-provider-ovh/pull/33))
* __New Resource:__ `ovh_iploadbalancing_http_route` ([#35](https://github.com/terraform-providers/terraform-provider-ovh/pull/35))
* __New Resource:__ `ovh_iploadbalancing_http_route_rule` ([#35](https://github.com/terraform-providers/terraform-provider-ovh/pull/35))
* __New Resource:__ `ovh_domain_zone_redirection` ([#36](https://github.com/terraform-providers/terraform-provider-ovh/pull/36))
* __New Resource:__ `ovh_cloud_network_private` ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* __New Resource:__ `ovh_cloud_network_private_subnet` ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* __New Resource:__ `ovh_cloud_user` ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))
* __New Resource:__ `ovh_vrack_cloudproject` ([#31](https://github.com/terraform-providers/terraform-provider-ovh/pull/31))

IMPROVEMENTS

* Various doc improvements ([#14](https://github.com/terraform-providers/terraform-provider-ovh/pull/14),[#16](https://github.com/terraform-providers/terraform-provider-ovh/pull/16),[#17](https://github.com/terraform-providers/terraform-provider-ovh/pull/17),[#26](https://github.com/terraform-providers/terraform-provider-ovh/pull/26),[#53](https://github.com/terraform-providers/terraform-provider-ovh/pull/51),[#51](https://github.com/terraform-providers/terraform-provider-ovh/pull/53))
* provider: Fallback to get current home directory ([#19](https://github.com/terraform-providers/terraform-provider-ovh/pull/19))
* provider: bump to terraform v0.10.8 ([#49](https://github.com/terraform-providers/terraform-provider-ovh/pull/49))
* r/ovh_domain_zone_record: add sweeper ([#50](https://github.com/terraform-providers/terraform-provider-ovh/pull/50))
* r/ovh_domain_zone_redirection: add sweeper ([#50](https://github.com/terraform-providers/terraform-provider-ovh/pull/50))


BUG FIXES:

* resource/ovh_domain_zone_record: Fixes [[#25](https://github.com/terraform-providers/terraform-provider-ovh/issues/25)] by id removal, cleans up naming and struct repetition for domain zone record resource
* provider: Fixes [[#27](https://github.com/terraform-providers/terraform-provider-ovh/issues/27)] by switching to `/auth/currentCredential` for client validation

## 0.2.0 (January 10, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* d/ovh_publiccloud_region: Deprecated fields which don't comply
  with lowercase & underscore convention (`continentCode`, `datacenterLocation`).
  Use `continent_code` and `datacenter_location` instead. ([#4](https://github.com/terraform-providers/terraform-provider-ovh/issues/4))

FEATURES

* __New Resource:__ `ovh_domain_zone_record` ([#3](https://github.com/terraform-providers/terraform-provider-ovh/issues/3))

IMPROVEMENTS

* The provider config can now source its credentials from `~/.ovh.conf` ([#10](https://github.com/terraform-providers/terraform-provider-ovh/issues/10))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
