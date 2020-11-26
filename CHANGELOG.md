## 0.10.0 (Unreleased)
## 0.9.0 (August 26, 2020)

BREAKING CHANGES:

* provider: This release includes a Terraform SDK upgrade with compatibility for Terraform >= v0.12. The provider is not compatible with Terraform < v0.12 anymore. This update should have no significant changes in behavior for the provider. ([#154](https://github.com/terraform-providers/terraform-provider-ovh/pull/154))

FEATURES:

* __New Datasource:__ `ovh_dedicated_ceph` ([#150](https://github.com/terraform-providers/terraform-provider-ovh/pull/150))
* __New Resource:__ `ovh_dedicated_ceph_acl` ([#150](https://github.com/terraform-providers/terraform-provider-ovh/pull/150))

IMPROVEMENTS:

* Fetch all IPs for dedicated servers. ([#149](https://github.com/terraform-providers/terraform-provider-ovh/pull/149))
* r/ovh_cloud_user: Add roles/service_name attributes. Deprecate project_id attribute. ([#151](https://github.com/terraform-providers/terraform-provider-ovh/pull/151))

BUG FIXES:

* r/iploadbalancing_tcp_frontend, r/iploadbalancing_http_frontend: Fix allowed_source,dedicated_ipfo updates ([#155](https://github.com/terraform-providers/terraform-provider-ovh/pull/155))

## 0.8.0 (May 28, 2020)

NOTES/DEPRECATIONS:

* `*/ovh_iploadbalancing_vrack_network_*`: Deprecate `farm_id` attribute as it conflicts with other resources. ([#144](https://github.com/terraform-providers/terraform-provider-ovh/issues/144))

FEATURES:

* __New Datasource:__ `ovh_me_ipxe_script` ([#141](https://github.com/terraform-providers/terraform-provider-ovh/pull/141))
* __New Datasource:__ `ovh_me_ipxe_scripts` ([#141](https://github.com/terraform-providers/terraform-provider-ovh/pull/141))
* __New Resource:__ `ovh_me_ipxe_script` ([#141](https://github.com/terraform-providers/terraform-provider-ovh/pull/141))

IMPROVEMENTS:

* Stop failing sweepers if mandatory env vars are missing. ([#142](https://github.com/terraform-providers/terraform-provider-ovh/pull/142))
* r/iploadbalancing_*: Add importers and tests([#140](https://github.com/terraform-providers/terraform-provider-ovh/pull/140))
* r/iploadbalancing_tcp_farm, r/iploadbalancing_http_farm: Extend read function to get all values ([#140](https://github.com/terraform-providers/terraform-provider-ovh/pull/140))
* r/iploadbalancing_http_route: Extend read function to get action value ([#140](https://github.com/terraform-providers/terraform-provider-ovh/pull/140))
* r/iploadbalancing_http_route_rule: Read all values ([#140](https://github.com/terraform-providers/terraform-provider-ovh/pull/140))
* r/iploadbalancing_tcp_frontend, r/iploadbalancing_http_frontend: Some code refactoring according to linter ([#140](https://github.com/terraform-providers/terraform-provider-ovh/pull/140))

BUG FIXES:

* r/iploadbalancing_vrack_network: fix sweepers ([#142](https://github.com/terraform-providers/terraform-provider-ovh/pull/142))
* r/iploadbalancing_tcp_farm, r/iploadbalancing_http_farm: Fix typo in 'oco' probe type ([#140](https://github.com/terraform-providers/terraform-provider-ovh/pull/140))
* r/iploadbalancing_tcp_farm_server, r/iploadbalancing_http_farm_server: Allow port to have a nil value  ([#140](https://github.com/terraform-providers/terraform-provider-ovh/pull/140))

## 0.7.0 (March 02, 2020)

FEATURES:

* __New Datasource:__ `ovh_vps` ([#126](https://github.com/terraform-providers/terraform-provider-ovh/pull/126))

IMPROVEMENTS:

* r/iploadbalancing_http_farm: add cookie stickiness ([#133](https://github.com/terraform-providers/terraform-provider-ovh/pull/133))
* r/dedicated_server_reboot_task, r/dedicated_server_install_task: retry task on 500/404 errors due to API unstability ([#134](https://github.com/terraform-providers/terraform-provider-ovh/pull/134))

BUG FIXES:

* r/dedicated_server_reboot_task, r/dedicated_server_install_task: fix missing ForcesNew attributes ([#135](https://github.com/terraform-providers/terraform-provider-ovh/pull/135))
* r/me_ssh_key: fix missing ForcesNew attribute ([#136](https://github.com/terraform-providers/terraform-provider-ovh/pull/136))
* r/domain_zone_record, domain_zone_redirection: don't fail sweepers on missing OVH_ZONE env var. ([#138](https://github.com/terraform-providers/terraform-provider-ovh/pull/138))
* r/cloud_network_private: fix sweeper. ([#138](https://github.com/terraform-providers/terraform-provider-ovh/pull/138))

## 0.6.0 (January 15, 2020)

FEATURES:

* __New Datasource:__ `ovh_dedicated_server` ([#100](https://github.com/terraform-providers/terraform-provider-ovh/pull/100))
* __New Datasource:__ `ovh_dedicated_servers` ([#100](https://github.com/terraform-providers/terraform-provider-ovh/pull/100))
* __New Datasource:__ `ovh_dedicated_server_boots` ([#105](https://github.com/terraform-providers/terraform-provider-ovh/pull/105))
* __New Datasource:__ `ovh_dedicated_server_boots` ([#105](https://github.com/terraform-providers/terraform-provider-ovh/pull/105))
* __New Datasource:__ `ovh_dedicated_server_installation_templates` ([#101](https://github.com/terraform-providers/terraform-provider-ovh/pull/101))
* __New Datasource:__ `ovh_iploadbalancing_vrack_network` ([#127](https://github.com/terraform-providers/terraform-provider-ovh/pull/127))
* __New Datasource:__ `ovh_iploadbalancing_vrack_networks` ([#127](https://github.com/terraform-providers/terraform-provider-ovh/pull/127))
* __New Datasource:__ `ovh_me_installation_template` ([#103](https://github.com/terraform-providers/terraform-provider-ovh/pull/103))
* __New Datasource:__ `ovh_me_installation_templates` ([#103](https://github.com/terraform-providers/terraform-provider-ovh/pull/103))
* __New Datasource:__ `ovh_me_ssh_key` ([#93](https://github.com/terraform-providers/terraform-provider-ovh/pull/93))
* __New Datasource:__ `ovh_me_ssh_keys` ([#93](https://github.com/terraform-providers/terraform-provider-ovh/pull/93))
* __New Datasource:__ `ovh_vracks` ([#114](https://github.com/terraform-providers/terraform-provider-ovh/pull/114))
* __New Resource:__ `ovh_dedicated_server_install_task` ([#117](https://github.com/terraform-providers/terraform-provider-ovh/pull/117))
* __New Resource:__ `ovh_dedicated_server_reboot_task` ([#116](https://github.com/terraform-providers/terraform-provider-ovh/pull/116))
* __New Resource:__ `ovh_dedicated_server_update` ([#116](https://github.com/terraform-providers/terraform-provider-ovh/pull/116))
* __New Resource:__ `ovh_iploadbalancing_vrack_network` ([#127](https://github.com/terraform-providers/terraform-provider-ovh/pull/127),[#129](https://github.com/terraform-providers/terraform-provider-ovh/pull/129))
* __New Resource:__ `ovh_me_installation_template` ([#103](https://github.com/terraform-providers/terraform-provider-ovh/pull/103))
* __New Resource:__ `ovh_me_installation_template_partition_scheme` ([#103](https://github.com/terraform-providers/terraform-provider-ovh/pull/103))
* __New Resource:__ `ovh_me_installation_template_partition_scheme_partition` ([#103](https://github.com/terraform-providers/terraform-provider-ovh/pull/103))
* __New Resource:__ `ovh_me_installation_template_partition_scheme_hardware_raid` ([#103](https://github.com/terraform-providers/terraform-provider-ovh/pull/103))
* __New Resource:__ `ovh_me_ssh_key` ([#93](https://github.com/terraform-providers/terraform-provider-ovh/pull/93))
* __New Resource:__ `ovh_vrack_dedicated_server` ([#115](https://github.com/terraform-providers/terraform-provider-ovh/pull/115))
* __New Resource:__ `ovh_vrack_dedicated_server_interface` ([#115](https://github.com/terraform-providers/terraform-provider-ovh/pull/115))
* __New Resource:__ `ovh_vrack_iploadbalancing` ([#127](https://github.com/terraform-providers/terraform-provider-ovh/pull/127))

IMPROVEMENTS:

* provider: bump to go 1.13 ([#104](https://github.com/terraform-providers/terraform-provider-ovh/pull/104),[#118](https://github.com/terraform-providers/terraform-provider-ovh/pull/118))
* provider: migrate to terraform-plugin-sdk ([#98](https://github.com/terraform-providers/terraform-provider-ovh/pull/98))
* provider: skip testacc if required env vars are missing ([#106](https://github.com/terraform-providers/terraform-provider-ovh/pull/106))
* d/cloud_regions: add "has_services_up" filter ([#112](https://github.com/terraform-providers/terraform-provider-ovh/pull/112))
* r/ip_reverse: add sweeper (([#99](https://github.com/terraform-providers/terraform-provider-ovh/pull/99), [#102](https://github.com/terraform-providers/terraform-provider-ovh/pull/102))
* acceptance tests: Add PreCheck for HTTP Loadbalancing ([#94](https://github.com/terraform-providers/terraform-provider-ovh/pull/94))
* r/cloud_network_private_subnet: add importer ([#124](https://github.com/terraform-providers/terraform-provider-ovh/pull/124))

BUG FIXES:

* helpers: Fix nil pointer funcs which return wrong golang values in case of HCL nil values ([#120](https://github.com/terraform-providers/terraform-provider-ovh/pull/120))
* r/cloud_network_private, r/cloud_network_private_subnet: fix acctest & rework ([#113](https://github.com/terraform-providers/terraform-provider-ovh/pull/113))
* handle record id bigger than 32bits ([#109](https://github.com/terraform-providers/terraform-provider-ovh/pull/109))
* docs: Correct variable escaping in ovh_iploadbalancing_http_route example ([#97](https://github.com/terraform-providers/terraform-provider-ovh/pull/97))
* docs: Add "ovh_iploadbalancing_refresh" to the website sidebar ([#96](https://github.com/terraform-providers/terraform-provider-ovh/pull/96))

## 0.5.0 (May 22, 2019)

NOTES:

* provider: This release includes only a Terraform SDK upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and this update should have no significant changes in behavior for the provider. ([#86](https://github.com/terraform-providers/terraform-provider-ovh/issues/86))

## 0.4.0 (May 22, 2019)

NOTES/DEPRECATIONS:

* `*/ovh_publiccloud_*`: Deprecate `publiccloud` data sources & resources (to issue warnings when used) ([#81](https://github.com/terraform-providers/terraform-provider-ovh/issues/81))

FEATURES:

* __New Resource:__ `ovh_iploadbalancing_tcp_frontend` ([#58](https://github.com/terraform-providers/terraform-provider-ovh/issues/58))

IMPROVEMENTS:

* provider: Enable request/response logging in `>=DEBUG` mode ([#77](https://github.com/terraform-providers/terraform-provider-ovh/issues/77))
* provider: Make homedir detection more robust ([#82](https://github.com/terraform-providers/terraform-provider-ovh/issues/82))

BUG FIXES:

* resource/ovh_domain_zone_record: Attempt retries to avoid errors caused by eventual consistency after creation/update/deletion ([#77](https://github.com/terraform-providers/terraform-provider-ovh/issues/77))
* resource/ovh_domain_zone_record: Make fieldtype non-updatable ([#84](https://github.com/terraform-providers/terraform-provider-ovh/issues/84))
* resource/ovh_domain_zone_redirection: Return errors from refreshing after creation/update/deletion ([#77](https://github.com/terraform-providers/terraform-provider-ovh/issues/77))

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
