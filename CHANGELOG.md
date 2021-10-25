## 0.17.0 (Unreleased)

## 0.16.0 (October 25, 2021)


FEATURES:

* __New Resources:__ `r/ovh_iploadbalancing_tcp_route`, `r/ovh_iploadbalancing_tcp_route_rule` ([#222](https://github.com/ovh/terraform-provider-ovh/pull/222))

Improvements:

* Add regions regions_status deprecation. ([#198](https://github.com/ovh/terraform-provider-ovh/pull/198), [#227](https://github.com/ovh/terraform-provider-ovh/pull/198))
* fix & improve: data sources and kubernetes resources ([#226](https://github.com/ovh/terraform-provider-ovh/pull/226))
* `r/cloud_project_user`: add importer ([#220](https://github.com/ovh/terraform-provider-ovh/pull/220))
* Add missing Ovh subsidiaries ([#224](https://github.com/ovh/terraform-provider-ovh/pull/224))

BUG FIXES:

* fix: use the right json annotation ([#218](https://github.com/ovh/terraform-provider-ovh/pull/218))
* data/cloud/project/kube: fix acctest match on version ([#223](https://github.com/ovh/terraform-provider-ovh/pull/223))



## 0.15.0 (July 7, 2021)

BREAKING CHANGES:

* `r/ip_reverse`: `ipreverse` is renamed `ip_reverse` and is now mandatory ([#209](https://github.com/ovh/terraform-provider-ovh/pull/209))


BUG FIXES:

* `r/dbaas_logs_input`: fix import function ([#205](https://github.com/ovh/terraform-provider-ovh/pull/205))

Improvements:

* Provider is now built against go v1.16. ([#206](https://github.com/ovh/terraform-provider-ovh/pull/206))


## 0.13.1 (June 28, 2021)

NOTE:

This release fixes the v0.13.0 release, which should have included #194 patchset, but
due to an issue during the release process, the resulting binaries published
on the terraform registry didn't include it.

BUG FIXES:

* `r/cloud_project_kube`: fix issue with empty version([#194](https://github.com/ovh/terraform-provider-ovh/pull/194))

## 0.14.0 (June 23, 2021)

IMPORTANT: This release introduces a new kind of resources which are able to order and terminate OVH products.
OVH products are generally not on demand products, and thus may generate heavy costs. To use these
resources, you have to register a default payment mean on your account. These resources are still in
beta, and should be used with care.


FEATURES:

* __New Datasource:__ `d/dbaas_logs_input_engine` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Datasource:__ `d/dbaas_logs_output_graylog_stream` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Datasource:__ `d/ip_service` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Datasource:__ `d/order_cart` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Datasource:__ `d/order_cart_product` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Datasource:__ `d/order_cart_product_options` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Datasource:__ `d/order_cart_product_options_plan` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Datasource:__ `d/order_cart_product_plan` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Resource:__ `r/cloud_project` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Resource:__ `r/dbaas_logs_input` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Resource:__ `r/dbaas_logs_output_graylog_stream` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Resource:__ `r/domain_zone` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Resource:__ `r/ip_service` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Resource:__ `r/iploadbalancing` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Resource:__ `r/vrack` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))
* __New Resource:__ `r/vrack_ip` ([#202](https://github.com/ovh/terraform-provider-ovh/pull/202))


Improvements:

* Forgotten documentation change about private networking in Kubernetes. ([#195](https://github.com/ovh/terraform-provider-ovh/pull/195))
* Do not search for iploadbalancing service when exact service name was passed ([#200](https://github.com/ovh/terraform-provider-ovh/pull/200))

## 0.13.0 (May 11, 2021)

BREAKING CHANGES:

`OVH_VRACK_ID` and `OVH_PROJECT_ID` environment variables support are removed in favor of `OVH_VRACK_SERVICE` and `OVH_CLOUD_PROJECT_SERVICE`. Accordingly, `project_id` and `vrack_id` arguments are removed from resources in favor of `service_name`.

* `d/cloud_region`: removed ([#193](https://github.com/ovh/terraform-provider-ovh/pull/193))
* `d/cloud_regions`: removed ([#193](https://github.com/ovh/terraform-provider-ovh/pull/193))
* `r/cloud_private_network`: removed ([#193](https://github.com/ovh/terraform-provider-ovh/pull/193))
* `r/cloud_private_network_subnet`: removed ([#193](https://github.com/ovh/terraform-provider-ovh/pull/193))
* `r/cloud_user`: removed ([#193](https://github.com/ovh/terraform-provider-ovh/pull/193))

Improvements:

* increase vrack task timeouts fomr 20m to 60m

BUG FIXES:

* `r/cloud_project_containerregistry`: fix region and plan_id arguments issue ([#193](https://github.com/ovh/terraform-provider-ovh/pull/193))
* `r/dedicated_server_install_task`: fix issue when OVH cleans its tasks database ([#193](https://github.com/ovh/terraform-provider-ovh/pull/193))
* `r/dedicated_server_reboot_task`: fix issue when OVH cleans its tasks database ([#193](https://github.com/ovh/terraform-provider-ovh/pull/193))
* `r/cloud_project_kube`: fix issue with empty version([#194](https://github.com/ovh/terraform-provider-ovh/pull/194))
* `r/cloud_project_kube_nodepool`: fix issue with computed arguments 
* Documentation: fix typo ([#191](https://github.com/ovh/terraform-provider-ovh/pull/191))


## 0.12.0 (May 7, 2021)

FEATURES:

* __New Datasource:__ `ovh_cloud_project_capabilities_containerregistry` ([#184](https://github.com/ovh/terraform-provider-ovh/pull/184))
* __New Datasource:__ `ovh_cloud_project_capabilities_containerregistry_filter` ([#184](https://github.com/ovh/terraform-provider-ovh/pull/184))
* __New Datasource:__ `ovh_cloud_project_containerregistries` ([#184](https://github.com/ovh/terraform-provider-ovh/pull/184))
* __New Datasource:__ `ovh_cloud_project_containerregistry` ([#184](https://github.com/ovh/terraform-provider-ovh/pull/184))
* __New Datasource:__ `ovh_cloud_project_containerregistry_users` ([#184](https://github.com/ovh/terraform-provider-ovh/pull/184))
* __New Resource:__ `ovh_cloud_project_containerregistry` ([#184](https://github.com/ovh/terraform-provider-ovh/pull/184))
* __New Resource:__ `ovh_cloud_project_containerregistry_user` ([#184](https://github.com/ovh/terraform-provider-ovh/pull/184))

IMPROVEMENTS:

* Various documentation fixes, improvements. ([#182](https://github.com/ovh/terraform-provider-ovh/pull/182), [#185](https://github.com/ovh/terraform-provider-ovh/pull/185))
* Better handle error type casting for ovh api errors ([#186](https://github.com/ovh/terraform-provider-ovh/pull/186))
* Add `private_network_id` and `autoscale` arguments to ovh_cloud_project_kube resources.
([#189](https://github.com/ovh/terraform-provider-ovh/pull/189))
* Add missing attributes to ovh_cloud_project_kube_nodepool resource.
([#189](https://github.com/ovh/terraform-provider-ovh/pull/189))

BUG FIXES:

* Updated the GPG key used to verify Terraform installs in response to the Terraform GPG key rotation. ([#750](https://github-redirect.dependabot.com/hashicorp/terraform-plugin-sdk/issues/739))
cf [Terraform Updates for HCSEC-2021-12](https://discuss.hashicorp.com/t/terraform-updates-for-hcsec-2021-12/23570)

## 0.11.0 (March 3, 2021)

FEATURES:

* __New Datasource:__ `ovh_cloud_project_kube` ([#180](https://github.com/ovh/terraform-provider-ovh/pull/180))
* __New Resource:__ `ovh_cloud_project_kube` ([#180](https://github.com/ovh/terraform-provider-ovh/pull/180))
* __New Resource:__ `ovh_cloud_project_kube_nodepool` ([#180](https://github.com/ovh/terraform-provider-ovh/pull/180))

IMPROVEMENTS:

* Documentation improvements on provider setup ([#181](https://github.com/ovh/terraform-provider-ovh/pull/181))


## 0.10.0 (December 7, 2020)

BREAKING CHANGES:

* d/publiccloud*: all datasources ovh_publiccloud_* are removed in favor of ovh_cloud_project_* ([#175](https://github.com/ovh/terraform-provider-ovh/pull/175))
* r/publiccloud*: all resources ovh_publiccloud_* are removed in favor of ovh_cloud_project_* ([#175](https://github.com/ovh/terraform-provider-ovh/pull/175))

NOTES/DEPRECATIONS:

* d/cloud*: all datasources ovh_cloud_* are deprecated in favor of ovh_cloud_project_* ([#175](https://github.com/ovh/terraform-provider-ovh/pull/175))
* r/cloud*: all resources ovh_cloud_* are deprecated in favor of ovh_cloud_project_* ([#175](https://github.com/ovh/terraform-provider-ovh/pull/175))
* d/cloud*: use service_name for identifier ([#173](https://github.com/ovh/terraform-provider-ovh/pull/173))
* r/cloud*: use service_name for identifier ([#173](https://github.com/ovh/terraform-provider-ovh/pull/173))

FEATURES:

* __New Datasource:__ `ovh_me_identity_user` ([#166](https://github.com/ovh/terraform-provider-ovh/pull/166))
* __New Datasource:__ `ovh_me_identity_users` ([#166](https://github.com/ovh/terraform-provider-ovh/pull/166))
* __New Resource:__ `ovh_me_identity_user` ([#166](https://github.com/ovh/terraform-provider-ovh/pull/166))

IMPROVEMENTS:

* enforce CheckDeleted on all resources read operations ([#176](https://github.com/ovh/terraform-provider-ovh/pull/176))
* cicd acceptance tests now run on OVH CDS build system, travis-ci is removed ([#174](https://github.com/ovh/terraform-provider-ovh/pull/174))
* migrate to new lib ovh/terraform-ovh-provider ([#172](https://github.com/ovh/terraform-provider-ovh/pull/172))
* r/iploadbalancing*: add missing sweepers ([#171](https://github.com/ovh/terraform-provider-ovh/pull/171))
* go-ovh lib: bump to v1.1.0 ([#170](https://github.com/ovh/terraform-provider-ovh/pull/170))
* add freebsd support ([#164](https://github.com/ovh/terraform-provider-ovh/pull/164))
* increase vrack task timeout to 20 minutes ([#38b610e](https://github.com/ovh/terraform-provider-ovh/commit/38b610e310b7478d5cbe53bdc2e3dd09581b1340))

BUG FIXES:

* r/iploadbalancing_http_farm: fix probe handling ([#178](https://github.com/ovh/terraform-provider-ovh/pull/178))
* r/iploadbalancing_tcp_farm: fix probe handling ([#178](https://github.com/ovh/terraform-provider-ovh/pull/178))
* r/dedicated_server_update: fix monitoring update ([#178](https://github.com/ovh/terraform-provider-ovh/pull/178))
* d/vps: Fix erroneous types([#164](https://github.com/ovh/terraform-provider-ovh/pull/164))
* r/me_ssh_key: fix setting key default property and handle key not found error ([#158](https://github.com/ovh/terraform-provider-ovh/pull/158))


## 0.9.0 (August 26, 2020)

BREAKING CHANGES:

* provider: This release includes a Terraform SDK upgrade with compatibility for Terraform >= v0.12. The provider is not compatible with Terraform < v0.12 anymore. This update should have no significant changes in behavior for the provider. ([#154](https://github.com/ovh/terraform-provider-ovh/pull/154))

FEATURES:

* __New Datasource:__ `ovh_dedicated_ceph` ([#150](https://github.com/ovh/terraform-provider-ovh/pull/150))
* __New Resource:__ `ovh_dedicated_ceph_acl` ([#150](https://github.com/ovh/terraform-provider-ovh/pull/150))

IMPROVEMENTS:

* Fetch all IPs for dedicated servers. ([#149](https://github.com/ovh/terraform-provider-ovh/pull/149))
* r/ovh_cloud_user: Add roles/service_name attributes. Deprecate project_id attribute. ([#151](https://github.com/ovh/terraform-provider-ovh/pull/151))

BUG FIXES:

* r/iploadbalancing_tcp_frontend, r/iploadbalancing_http_frontend: Fix allowed_source,dedicated_ipfo updates ([#155](https://github.com/ovh/terraform-provider-ovh/pull/155))

## 0.8.0 (May 28, 2020)

NOTES/DEPRECATIONS:

* `*/ovh_iploadbalancing_vrack_network_*`: Deprecate `farm_id` attribute as it conflicts with other resources. ([#144](https://github.com/ovh/terraform-provider-ovh/issues/144))

FEATURES:

* __New Datasource:__ `ovh_me_ipxe_script` ([#141](https://github.com/ovh/terraform-provider-ovh/pull/141))
* __New Datasource:__ `ovh_me_ipxe_scripts` ([#141](https://github.com/ovh/terraform-provider-ovh/pull/141))
* __New Resource:__ `ovh_me_ipxe_script` ([#141](https://github.com/ovh/terraform-provider-ovh/pull/141))

IMPROVEMENTS:

* Stop failing sweepers if mandatory env vars are missing. ([#142](https://github.com/ovh/terraform-provider-ovh/pull/142))
* r/iploadbalancing_*: Add importers and tests([#140](https://github.com/ovh/terraform-provider-ovh/pull/140))
* r/iploadbalancing_tcp_farm, r/iploadbalancing_http_farm: Extend read function to get all values ([#140](https://github.com/ovh/terraform-provider-ovh/pull/140))
* r/iploadbalancing_http_route: Extend read function to get action value ([#140](https://github.com/ovh/terraform-provider-ovh/pull/140))
* r/iploadbalancing_http_route_rule: Read all values ([#140](https://github.com/ovh/terraform-provider-ovh/pull/140))
* r/iploadbalancing_tcp_frontend, r/iploadbalancing_http_frontend: Some code refactoring according to linter ([#140](https://github.com/ovh/terraform-provider-ovh/pull/140))

BUG FIXES:

* r/iploadbalancing_vrack_network: fix sweepers ([#142](https://github.com/ovh/terraform-provider-ovh/pull/142))
* r/iploadbalancing_tcp_farm, r/iploadbalancing_http_farm: Fix typo in 'oco' probe type ([#140](https://github.com/ovh/terraform-provider-ovh/pull/140))
* r/iploadbalancing_tcp_farm_server, r/iploadbalancing_http_farm_server: Allow port to have a nil value  ([#140](https://github.com/ovh/terraform-provider-ovh/pull/140))

## 0.7.0 (March 02, 2020)

FEATURES:

* __New Datasource:__ `ovh_vps` ([#126](https://github.com/ovh/terraform-provider-ovh/pull/126))

IMPROVEMENTS:

* r/iploadbalancing_http_farm: add cookie stickiness ([#133](https://github.com/ovh/terraform-provider-ovh/pull/133))
* r/dedicated_server_reboot_task, r/dedicated_server_install_task: retry task on 500/404 errors due to API unstability ([#134](https://github.com/ovh/terraform-provider-ovh/pull/134))

BUG FIXES:

* r/dedicated_server_reboot_task, r/dedicated_server_install_task: fix missing ForcesNew attributes ([#135](https://github.com/ovh/terraform-provider-ovh/pull/135))
* r/me_ssh_key: fix missing ForcesNew attribute ([#136](https://github.com/ovh/terraform-provider-ovh/pull/136))
* r/domain_zone_record, domain_zone_redirection: don't fail sweepers on missing OVH_ZONE env var. ([#138](https://github.com/ovh/terraform-provider-ovh/pull/138))
* r/cloud_network_private: fix sweeper. ([#138](https://github.com/ovh/terraform-provider-ovh/pull/138))

## 0.6.0 (January 15, 2020)

FEATURES:

* __New Datasource:__ `ovh_dedicated_server` ([#100](https://github.com/ovh/terraform-provider-ovh/pull/100))
* __New Datasource:__ `ovh_dedicated_servers` ([#100](https://github.com/ovh/terraform-provider-ovh/pull/100))
* __New Datasource:__ `ovh_dedicated_server_boots` ([#105](https://github.com/ovh/terraform-provider-ovh/pull/105))
* __New Datasource:__ `ovh_dedicated_server_boots` ([#105](https://github.com/ovh/terraform-provider-ovh/pull/105))
* __New Datasource:__ `ovh_dedicated_server_installation_templates` ([#101](https://github.com/ovh/terraform-provider-ovh/pull/101))
* __New Datasource:__ `ovh_iploadbalancing_vrack_network` ([#127](https://github.com/ovh/terraform-provider-ovh/pull/127))
* __New Datasource:__ `ovh_iploadbalancing_vrack_networks` ([#127](https://github.com/ovh/terraform-provider-ovh/pull/127))
* __New Datasource:__ `ovh_me_installation_template` ([#103](https://github.com/ovh/terraform-provider-ovh/pull/103))
* __New Datasource:__ `ovh_me_installation_templates` ([#103](https://github.com/ovh/terraform-provider-ovh/pull/103))
* __New Datasource:__ `ovh_me_ssh_key` ([#93](https://github.com/ovh/terraform-provider-ovh/pull/93))
* __New Datasource:__ `ovh_me_ssh_keys` ([#93](https://github.com/ovh/terraform-provider-ovh/pull/93))
* __New Datasource:__ `ovh_vracks` ([#114](https://github.com/ovh/terraform-provider-ovh/pull/114))
* __New Resource:__ `ovh_dedicated_server_install_task` ([#117](https://github.com/ovh/terraform-provider-ovh/pull/117))
* __New Resource:__ `ovh_dedicated_server_reboot_task` ([#116](https://github.com/ovh/terraform-provider-ovh/pull/116))
* __New Resource:__ `ovh_dedicated_server_update` ([#116](https://github.com/ovh/terraform-provider-ovh/pull/116))
* __New Resource:__ `ovh_iploadbalancing_vrack_network` ([#127](https://github.com/ovh/terraform-provider-ovh/pull/127),[#129](https://github.com/ovh/terraform-provider-ovh/pull/129))
* __New Resource:__ `ovh_me_installation_template` ([#103](https://github.com/ovh/terraform-provider-ovh/pull/103))
* __New Resource:__ `ovh_me_installation_template_partition_scheme` ([#103](https://github.com/ovh/terraform-provider-ovh/pull/103))
* __New Resource:__ `ovh_me_installation_template_partition_scheme_partition` ([#103](https://github.com/ovh/terraform-provider-ovh/pull/103))
* __New Resource:__ `ovh_me_installation_template_partition_scheme_hardware_raid` ([#103](https://github.com/ovh/terraform-provider-ovh/pull/103))
* __New Resource:__ `ovh_me_ssh_key` ([#93](https://github.com/ovh/terraform-provider-ovh/pull/93))
* __New Resource:__ `ovh_vrack_dedicated_server` ([#115](https://github.com/ovh/terraform-provider-ovh/pull/115))
* __New Resource:__ `ovh_vrack_dedicated_server_interface` ([#115](https://github.com/ovh/terraform-provider-ovh/pull/115))
* __New Resource:__ `ovh_vrack_iploadbalancing` ([#127](https://github.com/ovh/terraform-provider-ovh/pull/127))

IMPROVEMENTS:

* provider: bump to go 1.13 ([#104](https://github.com/ovh/terraform-provider-ovh/pull/104),[#118](https://github.com/ovh/terraform-provider-ovh/pull/118))
* provider: migrate to terraform-plugin-sdk ([#98](https://github.com/ovh/terraform-provider-ovh/pull/98))
* provider: skip testacc if required env vars are missing ([#106](https://github.com/ovh/terraform-provider-ovh/pull/106))
* d/cloud_regions: add "has_services_up" filter ([#112](https://github.com/ovh/terraform-provider-ovh/pull/112))
* r/ip_reverse: add sweeper (([#99](https://github.com/ovh/terraform-provider-ovh/pull/99), [#102](https://github.com/ovh/terraform-provider-ovh/pull/102))
* acceptance tests: Add PreCheck for HTTP Loadbalancing ([#94](https://github.com/ovh/terraform-provider-ovh/pull/94))
* r/cloud_network_private_subnet: add importer ([#124](https://github.com/ovh/terraform-provider-ovh/pull/124))

BUG FIXES:

* helpers: Fix nil pointer funcs which return wrong golang values in case of HCL nil values ([#120](https://github.com/ovh/terraform-provider-ovh/pull/120))
* r/cloud_network_private, r/cloud_network_private_subnet: fix acctest & rework ([#113](https://github.com/ovh/terraform-provider-ovh/pull/113))
* handle record id bigger than 32bits ([#109](https://github.com/ovh/terraform-provider-ovh/pull/109))
* docs: Correct variable escaping in ovh_iploadbalancing_http_route example ([#97](https://github.com/ovh/terraform-provider-ovh/pull/97))
* docs: Add "ovh_iploadbalancing_refresh" to the website sidebar ([#96](https://github.com/ovh/terraform-provider-ovh/pull/96))

## 0.5.0 (May 22, 2019)

NOTES:

* provider: This release includes only a Terraform SDK upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and this update should have no significant changes in behavior for the provider. ([#86](https://github.com/ovh/terraform-provider-ovh/issues/86))

## 0.4.0 (May 22, 2019)

NOTES/DEPRECATIONS:

* `*/ovh_publiccloud_*`: Deprecate `publiccloud` data sources & resources (to issue warnings when used) ([#81](https://github.com/ovh/terraform-provider-ovh/issues/81))

FEATURES:

* __New Resource:__ `ovh_iploadbalancing_tcp_frontend` ([#58](https://github.com/ovh/terraform-provider-ovh/issues/58))

IMPROVEMENTS:

* provider: Enable request/response logging in `>=DEBUG` mode ([#77](https://github.com/ovh/terraform-provider-ovh/issues/77))
* provider: Make homedir detection more robust ([#82](https://github.com/ovh/terraform-provider-ovh/issues/82))

BUG FIXES:

* resource/ovh_domain_zone_record: Attempt retries to avoid errors caused by eventual consistency after creation/update/deletion ([#77](https://github.com/ovh/terraform-provider-ovh/issues/77))
* resource/ovh_domain_zone_record: Make fieldtype non-updatable ([#84](https://github.com/ovh/terraform-provider-ovh/issues/84))
* resource/ovh_domain_zone_redirection: Return errors from refreshing after creation/update/deletion ([#77](https://github.com/ovh/terraform-provider-ovh/issues/77))

## 0.3.0 (July 11, 2018)

DEPRECATIONS / NOTES:

Resources and datasources names now reflects the OVH API endpoints. As such,
resources & datasources that doesn't comply are now deprecated and will be removed
in next release.

* data-source/ovh_publiccloud_region: Deprecated in favor of data-source/ovh_cloud_region ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* data-source/ovh_publiccloud_regions: Deprecated in favor of data-source/ovh_cloud_regions ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* resource/ovh_publiccloud_private_network: Deprecated in favor of resource/ovh_cloud_network_private ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* resource/ovh_publiccloud_private_network_subnet: Deprecated in favor of resource/ovh_cloud_network_private_subnet ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* resource/ovh_publiccloud_user: Deprecated in favor of resource/ovh_cloud_user ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* resource/ovh_vrack_publiccloud_attachment: Deprecated in favor of resource/ovh_vrack_cloudproject ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))

FEATURES

* __New Datasource:__ `ovh_cloud_region` ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* __New Datasource:__ `ovh_cloud_regions` ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* __New Datasource:__ `ovh_domain_zone` ([#39](https://github.com/ovh/terraform-provider-ovh/pull/39))
* __New Datasource:__ `ovh_iploadbalancing` ([#40](https://github.com/ovh/terraform-provider-ovh/pull/40))
* __New Datasource:__ `ovh_me_paymentmean_creditcard` ([#34](https://github.com/ovh/terraform-provider-ovh/pull/34),[#52](https://github.com/ovh/terraform-provider-ovh/pull/52))
* __New Datasource:__ `ovh_me_paymentmean_bankaccount` ([#34](https://github.com/ovh/terraform-provider-ovh/pull/34),[#52](https://github.com/ovh/terraform-provider-ovh/pull/52))
* __New Resource:__ `ovh_iploadbalancing_tcp_farm` ([#32](https://github.com/ovh/terraform-provider-ovh/pull/32))
* __New Resource:__ `ovh_iploadbalancing_tcp_farm_server` ([#33](https://github.com/ovh/terraform-provider-ovh/pull/33))
* __New Resource:__ `ovh_iploadbalancing_http_route` ([#35](https://github.com/ovh/terraform-provider-ovh/pull/35))
* __New Resource:__ `ovh_iploadbalancing_http_route_rule` ([#35](https://github.com/ovh/terraform-provider-ovh/pull/35))
* __New Resource:__ `ovh_domain_zone_redirection` ([#36](https://github.com/ovh/terraform-provider-ovh/pull/36))
* __New Resource:__ `ovh_cloud_network_private` ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* __New Resource:__ `ovh_cloud_network_private_subnet` ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* __New Resource:__ `ovh_cloud_user` ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))
* __New Resource:__ `ovh_vrack_cloudproject` ([#31](https://github.com/ovh/terraform-provider-ovh/pull/31))

IMPROVEMENTS

* Various doc improvements ([#14](https://github.com/ovh/terraform-provider-ovh/pull/14),[#16](https://github.com/ovh/terraform-provider-ovh/pull/16),[#17](https://github.com/ovh/terraform-provider-ovh/pull/17),[#26](https://github.com/ovh/terraform-provider-ovh/pull/26),[#53](https://github.com/ovh/terraform-provider-ovh/pull/51),[#51](https://github.com/ovh/terraform-provider-ovh/pull/53))
* provider: Fallback to get current home directory ([#19](https://github.com/ovh/terraform-provider-ovh/pull/19))
* provider: bump to terraform v0.10.8 ([#49](https://github.com/ovh/terraform-provider-ovh/pull/49))
* r/ovh_domain_zone_record: add sweeper ([#50](https://github.com/ovh/terraform-provider-ovh/pull/50))
* r/ovh_domain_zone_redirection: add sweeper ([#50](https://github.com/ovh/terraform-provider-ovh/pull/50))


BUG FIXES:

* resource/ovh_domain_zone_record: Fixes [[#25](https://github.com/ovh/terraform-provider-ovh/issues/25)] by id removal, cleans up naming and struct repetition for domain zone record resource
* provider: Fixes [[#27](https://github.com/ovh/terraform-provider-ovh/issues/27)] by switching to `/auth/currentCredential` for client validation

## 0.2.0 (January 10, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* d/ovh_publiccloud_region: Deprecated fields which don't comply
  with lowercase & underscore convention (`continentCode`, `datacenterLocation`).
  Use `continent_code` and `datacenter_location` instead. ([#4](https://github.com/ovh/terraform-provider-ovh/issues/4))

FEATURES

* __New Resource:__ `ovh_domain_zone_record` ([#3](https://github.com/ovh/terraform-provider-ovh/issues/3))

IMPROVEMENTS

* The provider config can now source its credentials from `~/.ovh.conf` ([#10](https://github.com/ovh/terraform-provider-ovh/issues/10))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
