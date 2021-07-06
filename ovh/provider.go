package ovh

import (
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
	ini "gopkg.in/ini.v1"
)

// Provider returns a *schema.Provider for OVH.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_ENDPOINT", nil),
				Description: descriptions["endpoint"],
			},
			"application_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_APPLICATION_KEY", ""),
				Description: descriptions["application_key"],
			},
			"application_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_APPLICATION_SECRET", ""),
				Description: descriptions["application_secret"],
			},
			"consumer_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CONSUMER_KEY", ""),
				Description: descriptions["consumer_key"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ovh_cloud_project_capabilities_containerregistry":        dataSourceCloudProjectCapabilitiesContainerRegistry(),
			"ovh_cloud_project_capabilities_containerregistry_filter": dataSourceCloudProjectCapabilitiesContainerRegistryFilter(),
			"ovh_cloud_project_containerregistries":                   dataSourceCloudProjectContainerRegistries(),
			"ovh_cloud_project_containerregistry":                     dataSourceCloudProjectContainerRegistry(),
			"ovh_cloud_project_containerregistry_users":               dataSourceCloudProjectContainerRegistryUsers(),
			"ovh_cloud_project_kube":                                  dataSourceCloudProjectKube(),
			"ovh_cloud_project_region":                                dataSourceCloudProjectRegion(),
			"ovh_cloud_project_regions":                               dataSourceCloudProjectRegions(),
			"ovh_dbaas_logs_input_engine":                             dataSourceDbaasLogsInputEngine(),
			"ovh_dbaas_logs_output_graylog_stream":                    dataSourceDbaasLogsOutputGraylogStream(),
			"ovh_dedicated_ceph":                                      dataSourceDedicatedCeph(),
			"ovh_dedicated_installation_templates":                    dataSourceDedicatedInstallationTemplates(),
			"ovh_dedicated_server":                                    dataSourceDedicatedServer(),
			"ovh_dedicated_server_boots":                              dataSourceDedicatedServerBoots(),
			"ovh_dedicated_servers":                                   dataSourceDedicatedServers(),
			"ovh_domain_zone":                                         dataSourceDomainZone(),
			"ovh_ip_service":                                          dataSourceIpService(),
			"ovh_iploadbalancing":                                     dataSourceIpLoadbalancing(),
			"ovh_iploadbalancing_vrack_network":                       dataSourceIpLoadbalancingVrackNetwork(),
			"ovh_iploadbalancing_vrack_networks":                      dataSourceIpLoadbalancingVrackNetworks(),
			"ovh_me_identity_user":                                    dataSourceMeIdentityUser(),
			"ovh_me_identity_users":                                   dataSourceMeIdentityUsers(),
			"ovh_me_installation_template":                            dataSourceMeInstallationTemplate(),
			"ovh_me_installation_templates":                           dataSourceMeInstallationTemplates(),
			"ovh_me_ipxe_script":                                      dataSourceMeIpxeScript(),
			"ovh_me_ipxe_scripts":                                     dataSourceMeIpxeScripts(),
			"ovh_me_paymentmean_bankaccount":                          dataSourceMePaymentmeanBankaccount(),
			"ovh_me_paymentmean_creditcard":                           dataSourceMePaymentmeanCreditcard(),
			"ovh_me_ssh_key":                                          dataSourceMeSshKey(),
			"ovh_me_ssh_keys":                                         dataSourceMeSshKeys(),
			"ovh_order_cart":                                          dataSourceOrderCart(),
			"ovh_order_cart_product":                                  dataSourceOrderCartProduct(),
			"ovh_order_cart_product_options":                          dataSourceOrderCartProductOptions(),
			"ovh_order_cart_product_options_plan":                     dataSourceOrderCartProductOptionsPlan(),
			"ovh_order_cart_product_plan":                             dataSourceOrderCartProductPlan(),
			"ovh_vps":                                                 dataSourceVPS(),
			"ovh_vracks":                                              dataSourceVracks(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ovh_cloud_project":                                           resourceCloudProject(),
			"ovh_cloud_project_containerregistry":                         resourceCloudProjectContainerRegistry(),
			"ovh_cloud_project_containerregistry_user":                    resourceCloudProjectContainerRegistryUser(),
			"ovh_cloud_project_kube":                                      resourceCloudProjectKube(),
			"ovh_cloud_project_kube_nodepool":                             resourceCloudProjectKubeNodePool(),
			"ovh_cloud_project_network_private":                           resourceCloudProjectNetworkPrivate(),
			"ovh_cloud_project_network_private_subnet":                    resourceCloudProjectNetworkPrivateSubnet(),
			"ovh_cloud_project_user":                                      resourceCloudProjectUser(),
			"ovh_dbaas_logs_input":                                        resourceDbaasLogsInput(),
			"ovh_dbaas_logs_output_graylog_stream":                        resourceDbaasLogsOutputGraylogStream(),
			"ovh_dedicated_ceph_acl":                                      resourceDedicatedCephACL(),
			"ovh_dedicated_server_install_task":                           resourceDedicatedServerInstallTask(),
			"ovh_dedicated_server_reboot_task":                            resourceDedicatedServerRebootTask(),
			"ovh_dedicated_server_update":                                 resourceDedicatedServerUpdate(),
			"ovh_domain_zone":                                             resourceDomainZone(),
			"ovh_domain_zone_record":                                      resourceOvhDomainZoneRecord(),
			"ovh_domain_zone_redirection":                                 resourceOvhDomainZoneRedirection(),
			"ovh_ip_reverse":                                              resourceIpReverse(),
			"ovh_ip_service":                                              resourceIpService(),
			"ovh_iploadbalancing":                                         resourceIpLoadbalancing(),
			"ovh_iploadbalancing_http_farm":                               resourceIpLoadbalancingHttpFarm(),
			"ovh_iploadbalancing_http_farm_server":                        resourceIpLoadbalancingHttpFarmServer(),
			"ovh_iploadbalancing_http_frontend":                           resourceIpLoadbalancingHttpFrontend(),
			"ovh_iploadbalancing_http_route":                              resourceIPLoadbalancingRouteHTTP(),
			"ovh_iploadbalancing_http_route_rule":                         resourceIPLoadbalancingRouteHTTPRule(),
			"ovh_iploadbalancing_refresh":                                 resourceIPLoadbalancingRefresh(),
			"ovh_iploadbalancing_tcp_farm":                                resourceIpLoadbalancingTcpFarm(),
			"ovh_iploadbalancing_tcp_farm_server":                         resourceIpLoadbalancingTcpFarmServer(),
			"ovh_iploadbalancing_tcp_frontend":                            resourceIpLoadbalancingTcpFrontend(),
			"ovh_iploadbalancing_vrack_network":                           resourceIPLoadbalancingVrackNetwork(),
			"ovh_me_identity_user":                                        resourceMeIdentityUser(),
			"ovh_me_installation_template":                                resourceMeInstallationTemplate(),
			"ovh_me_installation_template_partition_scheme":               resourceMeInstallationTemplatePartitionScheme(),
			"ovh_me_installation_template_partition_scheme_hardware_raid": resourceMeInstallationTemplatePartitionSchemeHardwareRaid(),
			"ovh_me_installation_template_partition_scheme_partition":     resourceMeInstallationTemplatePartitionSchemePartition(),
			"ovh_me_ipxe_script":                                          resourceMeIpxeScript(),
			"ovh_me_ssh_key":                                              resourceMeSshKey(),
			"ovh_vrack":                                                   resourceVrack(),
			"ovh_vrack_cloudproject":                                      resourceVrackCloudProject(),
			"ovh_vrack_dedicated_server":                                  resourceVrackDedicatedServer(),
			"ovh_vrack_dedicated_server_interface":                        resourceVrackDedicatedServerInterface(),
			"ovh_vrack_ip":                                                resourceVrackIp(),
			"ovh_vrack_iploadbalancing":                                   resourceVrackIpLoadbalancing(),
		},

		ConfigureFunc: configureProvider,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"endpoint": "The OVH API endpoint to target (ex: \"ovh-eu\").",

		"application_key": "The OVH API Application Key.",

		"application_secret": "The OVH API Application Secret.",
		"consumer_key":       "The OVH API Consumer key.",
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Endpoint: d.Get("endpoint").(string),
		lockAuth: &sync.Mutex{},
	}

	rawPath := "~/.ovh.conf"
	configPath, err := homedir.Expand(rawPath)
	if err != nil {
		return &config, fmt.Errorf("Failed to expand config path %q: %s", rawPath, err)
	}

	if _, err := os.Stat(configPath); err == nil {
		c, err := ini.Load(configPath)
		if err != nil {
			return nil, err
		}

		section, err := c.GetSection(d.Get("endpoint").(string))
		if err != nil {
			return nil, err
		}
		config.ApplicationKey = section.Key("application_key").String()
		config.ApplicationSecret = section.Key("application_secret").String()
		config.ConsumerKey = section.Key("consumer_key").String()
	}

	if v, ok := d.GetOk("application_key"); ok {
		config.ApplicationKey = v.(string)
	}
	if v, ok := d.GetOk("application_secret"); ok {
		config.ApplicationSecret = v.(string)
	}
	if v, ok := d.GetOk("consumer_key"); ok {
		config.ConsumerKey = v.(string)
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func deprecated(r *schema.Resource, msg string) *schema.Resource {
	r.DeprecationMessage = msg
	return r
}
