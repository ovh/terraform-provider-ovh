package ovh

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	SERVICE_NAME_ATT = "service_name"
)

var (
	// descriptions contains the descriptions for the fields of the providers' schemas.
	// The descriptions are grouped here because we need to have the exact same description
	// in each provider used by the MuxServer, else it doesn't boot.
	descriptions = map[string]string{
		"endpoint": "The OVH API endpoint to target (ex: \"ovh-eu\")",

		// Authentication via app key / app secret / comsumer key
		"application_key":    "The OVH API Application Key",
		"application_secret": "The OVH API Application Secret",
		"consumer_key":       "The OVH API Consumer Key",

		// Authentication via oAuth2
		"client_id":     "OAuth 2.0 application's ID",
		"client_secret": "OAuth 2.0 application's secret",
	}
)

// Provider returns a *schema.Provider for OVH.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["endpoint"],
			},
			"application_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["application_key"],
			},
			"application_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["application_secret"],
			},
			"consumer_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["consumer_key"],
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["client_id"],
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["client_secret"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ovh_cloud_project_capabilities_containerregistry":               dataSourceCloudProjectCapabilitiesContainerRegistry(),
			"ovh_cloud_project_capabilities_containerregistry_filter":        dataSourceCloudProjectCapabilitiesContainerRegistryFilter(),
			"ovh_cloud_project_containerregistry_ip_restrictions_management": dataSourceCloudProjectContainerRegistryIPRestrictionsManagement(),
			"ovh_cloud_project_containerregistry_ip_restrictions_registry":   dataSourceCloudProjectContainerRegistryIPRestrictionsRegistry(),
			"ovh_cloud_project_containerregistries":                          dataSourceCloudProjectContainerRegistries(),
			"ovh_cloud_project_containerregistry":                            dataSourceCloudProjectContainerRegistry(),
			"ovh_cloud_project_containerregistry_oidc":                       dataSourceCloudProjectContainerRegistryOIDC(),
			"ovh_cloud_project_containerregistry_users":                      dataSourceCloudProjectContainerRegistryUsers(),
			"ovh_cloud_project_database":                                     dataSourceCloudProjectDatabase(),
			"ovh_cloud_project_databases":                                    dataSourceCloudProjectDatabases(),
			"ovh_cloud_project_database_capabilities":                        dataSourceCloudProjectDatabaseCapabilities(),
			"ovh_cloud_project_database_certificates":                        dataSourceCloudProjectDatabaseCertificates(),
			"ovh_cloud_project_database_database":                            dataSourceCloudProjectDatabaseDatabase(),
			"ovh_cloud_project_database_databases":                           dataSourceCloudProjectDatabaseDatabases(),
			"ovh_cloud_project_database_integration":                         dataSourceCloudProjectDatabaseIntegration(),
			"ovh_cloud_project_database_integrations":                        dataSourceCloudProjectDatabaseIntegrations(),
			"ovh_cloud_project_database_kafka_acl":                           dataSourceCloudProjectDatabaseKafkaACL(),
			"ovh_cloud_project_database_kafka_acls":                          dataSourceCloudProjectDatabaseKafkaAcls(),
			"ovh_cloud_project_database_kafka_schemaregistryacl":             dataSourceCloudProjectDatabaseKafkaSchemaRegistryAcl(),
			"ovh_cloud_project_database_kafka_schemaregistryacls":            dataSourceCloudProjectDatabaseKafkaSchemaRegistryAcls(),
			"ovh_cloud_project_database_kafka_topic":                         dataSourceCloudProjectDatabaseKafkaTopic(),
			"ovh_cloud_project_database_kafka_topics":                        dataSourceCloudProjectDatabaseKafkaTopics(),
			"ovh_cloud_project_database_kafka_user_access":                   dataSourceCloudProjectDatabaseKafkaUserAccess(),
			"ovh_cloud_project_database_log_subscription":                    dataSourceCloudProjectDatabaseLogSubscription(),
			"ovh_cloud_project_database_log_subscriptions":                   dataSourceCloudProjectDatabaseLogSubscriptions(),
			"ovh_cloud_project_database_m3db_namespace":                      dataSourceCloudProjectDatabaseM3dbNamespace(),
			"ovh_cloud_project_database_m3db_namespaces":                     dataSourceCloudProjectDatabaseM3dbNamespaces(),
			"ovh_cloud_project_database_m3db_user":                           dataSourceCloudProjectDatabaseM3dbUser(),
			"ovh_cloud_project_database_mongodb_user":                        dataSourceCloudProjectDatabaseMongodbUser(),
			"ovh_cloud_project_database_opensearch_pattern":                  dataSourceCloudProjectDatabaseOpensearchPattern(),
			"ovh_cloud_project_database_opensearch_patterns":                 dataSourceCloudProjectDatabaseOpensearchPatterns(),
			"ovh_cloud_project_database_opensearch_user":                     dataSourceCloudProjectDatabaseOpensearchUser(),
			"ovh_cloud_project_database_postgresql_user":                     dataSourceCloudProjectDatabasePostgresqlUser(),
			"ovh_cloud_project_database_postgresql_connection_pool":          dataSourceCloudProjectDatabasePostgresqlConnectionPool(),
			"ovh_cloud_project_database_postgresql_connection_pools":         dataSourceCloudProjectDatabasePostgresqlConnectionPools(),
			"ovh_cloud_project_database_redis_user":                          dataSourceCloudProjectDatabaseRedisUser(),
			"ovh_cloud_project_database_user":                                dataSourceCloudProjectDatabaseUser(),
			"ovh_cloud_project_database_users":                               dataSourceCloudProjectDatabaseUsers(),
			"ovh_cloud_project_failover_ip_attach":                           dataSourceCloudProjectFailoverIpAttach(),
			"ovh_cloud_project_kube":                                         dataSourceCloudProjectKube(),
			"ovh_cloud_project_kube_iprestrictions":                          dataSourceCloudProjectKubeIPRestrictions(),
			"ovh_cloud_project_kube_nodepool_nodes":                          dataSourceCloudProjectKubeNodepoolNodes(),
			"ovh_cloud_project_kube_oidc":                                    dataSourceCloudProjectKubeOIDC(),
			"ovh_cloud_project_kube_nodepool":                                dataSourceCloudProjectKubeNodepool(),
			"ovh_cloud_project_kube_nodes":                                   dataSourceCloudProjectKubeNodes(),
			"ovh_cloud_project_region":                                       dataSourceCloudProjectRegion(),
			"ovh_cloud_project_regions":                                      dataSourceCloudProjectRegions(),
			"ovh_cloud_project_user":                                         datasourceCloudProjectUser(),
			"ovh_cloud_project_user_s3_credential":                           dataCloudProjectUserS3Credential(),
			"ovh_cloud_project_user_s3_credentials":                          dataCloudProjectUserS3Credentials(),
			"ovh_cloud_project_user_s3_policy":                               dataCloudProjectUserS3Policy(),
			"ovh_cloud_project_users":                                        datasourceCloudProjectUsers(),
			"ovh_cloud_project_vrack":                                        dataSourceCloudProjectVrack(),
			"ovh_dbaas_logs_cluster":                                         dataSourceDbaasLogsCluster(),
			"ovh_dbaas_logs_clusters":                                        dataSourceDbaasLogsClusters(),
			"ovh_dbaas_logs_input_engine":                                    dataSourceDbaasLogsInputEngine(),
			"ovh_dbaas_logs_output_graylog_stream":                           dataSourceDbaasLogsOutputGraylogStream(),
			"ovh_dedicated_ceph":                                             dataSourceDedicatedCeph(),
			"ovh_dedicated_installation_templates":                           dataSourceDedicatedInstallationTemplates(),
			"ovh_dedicated_installation_template":                            dataSourceDedicatedInstallationTemplate(),
			"ovh_dedicated_nasha":                                            dataSourceDedicatedNasha(),
			"ovh_dedicated_server":                                           dataSourceDedicatedServer(),
			"ovh_dedicated_server_boots":                                     dataSourceDedicatedServerBoots(),
			"ovh_dedicated_servers":                                          dataSourceDedicatedServers(),
			"ovh_hosting_privatedatabase":                                    dataSourceHostingPrivateDatabase(),
			"ovh_hosting_privatedatabase_database":                           dataSourceHostingPrivateDatabaseDatabase(),
			"ovh_hosting_privatedatabase_user":                               dataSourceHostingPrivateDatabaseUser(),
			"ovh_hosting_privatedatabase_user_grant":                         dataSourceHostingPrivateDatabaseUserGrant(),
			"ovh_hosting_privatedatabase_whitelist":                          dataSourceHostingPrivateDatabaseWhitelist(),
			"ovh_domain_zone":                                                dataSourceDomainZone(),
			"ovh_iam_permissions_group":                                      dataSourceIamPermissionsGroup(),
			"ovh_iam_permissions_groups":                                     dataSourceIamPermissionsGroups(),
			"ovh_iam_policies":                                               dataSourceIamPolicies(),
			"ovh_iam_policy":                                                 dataSourceIamPolicy(),
			"ovh_iam_reference_actions":                                      dataSourceIamReferenceActions(),
			"ovh_iam_reference_resource_type":                                dataSourceIamReferenceResourceType(),
			"ovh_iam_resource_group":                                         dataSourceIamResourceGroup(),
			"ovh_iam_resource_groups":                                        dataSourceIamResourceGroups(),
			"ovh_ip_service":                                                 dataSourceIpService(),
			"ovh_iploadbalancing":                                            dataSourceIpLoadbalancing(),
			"ovh_iploadbalancing_vrack_network":                              dataSourceIpLoadbalancingVrackNetwork(),
			"ovh_iploadbalancing_vrack_networks":                             dataSourceIpLoadbalancingVrackNetworks(),
			"ovh_me":                                                         dataSourceMe(),
			"ovh_me_api_oauth2_client":                                       dataSourceMeApiOauth2Client(),
			"ovh_me_api_oauth2_clients":                                      dataSourceMeApiOauth2Clients(),
			"ovh_me_identity_group":                                          dataSourceMeIdentityGroup(),
			"ovh_me_identity_groups":                                         dataSourceMeIdentityGroups(),
			"ovh_me_identity_user":                                           dataSourceMeIdentityUser(),
			"ovh_me_identity_users":                                          dataSourceMeIdentityUsers(),
			"ovh_me_installation_template":                                   dataSourceMeInstallationTemplate(),
			"ovh_me_installation_templates":                                  dataSourceMeInstallationTemplates(),
			"ovh_me_paymentmean_bankaccount":                                 dataSourceMePaymentmeanBankaccount(),
			"ovh_me_paymentmean_creditcard":                                  dataSourceMePaymentmeanCreditcard(),
			"ovh_order_cart":                                                 dataSourceOrderCart(),
			"ovh_order_cart_product":                                         dataSourceOrderCartProduct(),
			"ovh_order_cart_product_options":                                 dataSourceOrderCartProductOptions(),
			"ovh_order_cart_product_options_plan":                            dataSourceOrderCartProductOptionsPlan(),
			"ovh_order_cart_product_plan":                                    dataSourceOrderCartProductPlan(),
			"ovh_vps":                                                        dataSourceVPS(),
			"ovh_vpss":                                                       dataSourceVPSs(),
			"ovh_vracks":                                                     dataSourceVracks(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ovh_cloud_project":                                              resourceCloudProject(),
			"ovh_cloud_project_containerregistry":                            resourceCloudProjectContainerRegistry(),
			"ovh_cloud_project_containerregistry_oidc":                       resourceCloudProjectContainerRegistryOIDC(),
			"ovh_cloud_project_containerregistry_user":                       resourceCloudProjectContainerRegistryUser(),
			"ovh_cloud_project_containerregistry_ip_restrictions_management": resourceCloudProjectContainerRegistryIPRestrictionsManagement(),
			"ovh_cloud_project_containerregistry_ip_restrictions_registry":   resourceCloudProjectContainerRegistryIPRestrictionsRegistry(),
			"ovh_cloud_project_database":                                     resourceCloudProjectDatabase(),
			"ovh_cloud_project_database_database":                            resourceCloudProjectDatabaseDatabase(),
			"ovh_cloud_project_database_integration":                         resourceCloudProjectDatabaseIntegration(),
			"ovh_cloud_project_database_ip_restriction":                      resourceCloudProjectDatabaseIpRestriction(),
			"ovh_cloud_project_database_kafka_acl":                           resourceCloudProjectDatabaseKafkaAcl(),
			"ovh_cloud_project_database_kafka_schemaregistryacl":             resourceCloudProjectDatabaseKafkaSchemaRegistryAcl(),
			"ovh_cloud_project_database_kafka_topic":                         resourceCloudProjectDatabaseKafkaTopic(),
			"ovh_cloud_project_database_log_subscription":                    resourceCloudProjectDatabaseLogSubscription(),
			"ovh_cloud_project_database_m3db_namespace":                      resourceCloudProjectDatabaseM3dbNamespace(),
			"ovh_cloud_project_database_m3db_user":                           resourceCloudProjectDatabaseM3dbUser(),
			"ovh_cloud_project_database_mongodb_user":                        resourceCloudProjectDatabaseMongodbUser(),
			"ovh_cloud_project_database_opensearch_pattern":                  resourceCloudProjectDatabaseOpensearchPattern(),
			"ovh_cloud_project_database_opensearch_user":                     resourceCloudProjectDatabaseOpensearchUser(),
			"ovh_cloud_project_database_postgresql_user":                     resourceCloudProjectDatabasePostgresqlUser(),
			"ovh_cloud_project_database_postgresql_connection_pool":          resourceCloudProjectDatabasePostgresqlConnectionPool(),
			"ovh_cloud_project_database_redis_user":                          resourceCloudProjectDatabaseRedisUser(),
			"ovh_cloud_project_database_user":                                resourceCloudProjectDatabaseUser(),
			"ovh_cloud_project_failover_ip_attach":                           resourceCloudProjectFailoverIpAttach(),
			"ovh_cloud_project_gateway":                                      resourceCloudProjectGateway(),
			"ovh_cloud_project_kube":                                         resourceCloudProjectKube(),
			"ovh_cloud_project_kube_nodepool":                                resourceCloudProjectKubeNodePool(),
			"ovh_cloud_project_kube_oidc":                                    resourceCloudProjectKubeOIDC(),
			"ovh_cloud_project_kube_iprestrictions":                          resourceCloudProjectKubeIpRestrictions(),
			"ovh_cloud_project_network_private":                              resourceCloudProjectNetworkPrivate(),
			"ovh_cloud_project_network_private_subnet":                       resourceCloudProjectNetworkPrivateSubnet(),
			"ovh_cloud_project_region_storage_presign":                       resourceCloudProjectRegionStoragePresign(),
			"ovh_cloud_project_user":                                         resourceCloudProjectUser(),
			"ovh_cloud_project_user_s3_credential":                           resourceCloudProjectUserS3Credential(),
			"ovh_cloud_project_user_s3_policy":                               resourceCloudProjectUserS3Policy(),
			"ovh_cloud_project_workflow_backup":                              resourceCloudProjectWorkflowBackup(),
			"ovh_dbaas_logs_cluster":                                         resourceDbaasLogsCluster(),
			"ovh_dbaas_logs_input":                                           resourceDbaasLogsInput(),
			"ovh_dbaas_logs_output_graylog_stream":                           resourceDbaasLogsOutputGraylogStream(),
			"ovh_dedicated_ceph_acl":                                         resourceDedicatedCephACL(),
			"ovh_dedicated_nasha_partition":                                  resourceDedicatedNASHAPartition(),
			"ovh_dedicated_nasha_partition_access":                           resourceDedicatedNASHAPartitionAccess(),
			"ovh_dedicated_nasha_partition_snapshot":                         resourceDedicatedNASHAPartitionSnapshot(),
			"ovh_dedicated_server_install_task":                              resourceDedicatedServerInstallTask(),
			"ovh_dedicated_server_reboot_task":                               resourceDedicatedServerRebootTask(),
			"ovh_dedicated_server_update":                                    resourceDedicatedServerUpdate(),
			"ovh_dedicated_server_networking":                                resourceDedicatedServerNetworking(),
			"ovh_domain_zone":                                                resourceDomainZone(),
			"ovh_domain_zone_record":                                         resourceOvhDomainZoneRecord(),
			"ovh_domain_zone_redirection":                                    resourceOvhDomainZoneRedirection(),
			"ovh_hosting_privatedatabase":                                    resourceHostingPrivateDatabase(),
			"ovh_hosting_privatedatabase_database":                           resourceHostingPrivateDatabaseDatabase(),
			"ovh_hosting_privatedatabase_user":                               resourceHostingPrivateDatabaseUser(),
			"ovh_hosting_privatedatabase_user_grant":                         resourceHostingPrivateDatabaseUserGrant(),
			"ovh_hosting_privatedatabase_whitelist":                          resourceHostingPrivateDatabaseWhitelist(),
			"ovh_iam_policy":                                                 resourceIamPolicy(),
			"ovh_iam_permissions_group":                                      resourceIamPermissionsGroup(),
			"ovh_iam_resource_group":                                         resourceIamResourceGroup(),
			"ovh_ip_reverse":                                                 resourceIpReverse(),
			"ovh_ip_service":                                                 resourceIpService(),
			"ovh_ip_move":                                                    resourceIpServiceMove(),
			"ovh_iploadbalancing":                                            resourceIpLoadbalancing(),
			"ovh_iploadbalancing_http_farm":                                  resourceIpLoadbalancingHttpFarm(),
			"ovh_iploadbalancing_http_farm_server":                           resourceIpLoadbalancingHttpFarmServer(),
			"ovh_iploadbalancing_http_frontend":                              resourceIpLoadbalancingHttpFrontend(),
			"ovh_iploadbalancing_http_route":                                 resourceIPLoadbalancingHttpRoute(),
			"ovh_iploadbalancing_http_route_rule":                            resourceIPLoadbalancingHttpRouteRule(),
			"ovh_iploadbalancing_refresh":                                    resourceIPLoadbalancingRefresh(),
			"ovh_iploadbalancing_tcp_farm":                                   resourceIpLoadbalancingTcpFarm(),
			"ovh_iploadbalancing_tcp_farm_server":                            resourceIpLoadbalancingTcpFarmServer(),
			"ovh_iploadbalancing_tcp_frontend":                               resourceIpLoadbalancingTcpFrontend(),
			"ovh_iploadbalancing_tcp_route":                                  resourceIPLoadbalancingTcpRoute(),
			"ovh_iploadbalancing_tcp_route_rule":                             resourceIPLoadbalancingTcpRouteRule(),
			"ovh_iploadbalancing_vrack_network":                              resourceIPLoadbalancingVrackNetwork(),
			"ovh_me_identity_group":                                          resourceMeIdentityGroup(),
			"ovh_me_api_oauth2_client":                                       resourceApiOauth2Client(),
			"ovh_me_identity_user":                                           resourceMeIdentityUser(),
			"ovh_me_installation_template":                                   resourceMeInstallationTemplate(),
			"ovh_me_installation_template_partition_scheme":                  resourceMeInstallationTemplatePartitionScheme(),
			"ovh_me_installation_template_partition_scheme_hardware_raid":    resourceMeInstallationTemplatePartitionSchemeHardwareRaid(),
			"ovh_me_installation_template_partition_scheme_partition":        resourceMeInstallationTemplatePartitionSchemePartition(),
			"ovh_vrack":                            resourceVrack(),
			"ovh_vrack_cloudproject":               resourceVrackCloudProject(),
			"ovh_vrack_dedicated_server":           resourceVrackDedicatedServer(),
			"ovh_vrack_dedicated_server_interface": resourceVrackDedicatedServerInterface(),
			"ovh_vrack_ip":                         resourceVrackIp(),
			"ovh_vrack_iploadbalancing":            resourceVrackIpLoadbalancing(),
		},

		ConfigureContextFunc: ConfigureContextFunc,
	}
}

func ConfigureContextFunc(context context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		lockAuth: &sync.Mutex{},
	}

	if v, ok := d.GetOk("endpoint"); ok {
		config.Endpoint = v.(string)
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
	if v, ok := d.GetOk("client_id"); ok {
		config.ClientID = v.(string)
	}
	if v, ok := d.GetOk("client_secret"); ok {
		config.ClientSecret = v.(string)
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}

	return &config, nil
}
