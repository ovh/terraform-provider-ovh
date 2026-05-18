package ovh

// Common attributes shared across kube resources
const (
	kubeServiceNameKey = "service_name"
	kubeKubeIdKey      = "kube_id"
	kubeRegionKey      = "region"
	kubeNameKey        = "name"
	kubeVersionKey     = "version"
	kubeStatusKey      = "status"
	kubeCreatedAtKey   = "created_at"
	kubeUpdatedAtKey   = "updated_at"
	kubeProjectIdKey   = "project_id"
	kubeNodesKey       = "nodes"
	kubeFlavorKey      = "flavor"
)

// Cluster attributes
const (
	kubeClusterLoadBalancersSubnetIdKey       = "load_balancers_subnet_id"
	kubeClusterNodesSubnetIdKey               = "nodes_subnet_id"
	kubeClusterPrivateNetworkIDKey            = "private_network_id"
	kubeClusterPrivateNetworkConfigurationKey = "private_network_configuration"
	kubeClusterUpdatePolicyKey                = "update_policy"
	kubeClusterPlanKey                        = "plan"
	kubeClusterProxyModeKey                   = "kube_proxy_mode"
	kubeClusterIsUpToDateKey                  = "is_up_to_date"
	kubeClusterControlPlaneIsUpToDateKey      = "control_plane_is_up_to_date"
	kubeClusterNextUpgradeVersionsKey         = "next_upgrade_versions"
	kubeClusterNodesUrlKey                    = "nodes_url"
	kubeClusterUrlKey                         = "url"
	kubeClusterKubeconfigKey                  = "kubeconfig"
	kubeClusterKubeconfigAttributesKey        = "kubeconfig_attributes"
	kubeClusterIpAllocationPolicyKey          = "ip_allocation_policy"
	kubeClusterPodsIpv4CidrKey                = "pods_ipv4_cidr"
	kubeClusterServicesIpv4CidrKey            = "services_ipv4_cidr"
	kubeClusterDefaultVrackGatewayKey         = "default_vrack_gateway"
	kubeClusterPrivateNetworkRoutingAsDefault = "private_network_routing_as_default"

	// Deprecated
	kubeClusterCustomization = "customization"

	kubeClusterCustomizationApiServerKey = "customization_apiserver"
	kubeClusterCustomizationKubeProxyKey = "customization_kube_proxy"
	kubeClusterCustomizationCiliumKey    = "customization_cilium"

	// customization_apiserver / customization sub-attributes
	kubeClusterCustomizationAdmissionPluginsKey = "admissionplugins"
	kubeClusterCustomizationEnabledKey          = "enabled"
	kubeClusterCustomizationDisabledKey         = "disabled"
	kubeClusterCustomizationApiServerNestedKey  = "apiserver"

	// customization_kube_proxy sub-attributes
	kubeClusterCustomizationIptablesKey      = "iptables"
	kubeClusterCustomizationIpvsKey          = "ipvs"
	kubeClusterCustomizationMinSyncPeriodKey = "min_sync_period"
	kubeClusterCustomizationSyncPeriodKey    = "sync_period"
	kubeClusterCustomizationSchedulerKey     = "scheduler"
	kubeClusterCustomizationTcpFinTimeoutKey = "tcp_fin_timeout"
	kubeClusterCustomizationTcpTimeoutKey    = "tcp_timeout"
	kubeClusterCustomizationUdpTimeoutKey    = "udp_timeout"

	// customization_cilium sub-attributes
	kubeClusterCiliumClusterMeshKey    = "cluster_mesh"
	kubeClusterCiliumClusterID         = "cluster_id"
	kubeClusterCiliumHubbleKey         = "hubble"
	kubeClusterCiliumApiServerKey      = "api_server"
	kubeClusterCiliumNodePortKey       = "node_port"
	kubeClusterCiliumServiceTypeKey    = "service_type"
	kubeClusterCiliumRelayKey          = "relay"
	kubeClusterCiliumUiKey             = "ui"
	kubeClusterCiliumBackendResources  = "backend_resources"
	kubeClusterCiliumFrontendResources = "frontend_resources"
	kubeClusterCiliumLimitsKey         = "limits"
	kubeClusterCiliumRequestsKey       = "requests"
	kubeClusterCiliumCpuKey            = "cpu"
	kubeClusterCiliumMemoryKey         = "memory"

	// kubeconfig_attributes sub-attributes
	kubeClusterKubeconfigHostKey                 = "host"
	kubeClusterKubeconfigClusterCaCertificateKey = "cluster_ca_certificate"
	kubeClusterKubeconfigClientCertificateKey    = "client_certificate"
	kubeClusterKubeconfigClientKeyKey            = "client_key"
)

// NodePool attributes
const (
	kubeNodePoolFlavorNameKey                               = "flavor_name"
	kubeNodePoolDesiredNodesKey                             = "desired_nodes"
	kubeNodePoolMinNodesKey                                 = "min_nodes"
	kubeNodePoolMaxNodesKey                                 = "max_nodes"
	kubeNodePoolAutoscaleKey                                = "autoscale"
	kubeNodePoolAutoscalingScaleDownUnneededTimeSecondsKey  = "autoscaling_scale_down_unneeded_time_seconds"
	kubeNodePoolAutoscalingScaleDownUnreadyTimeSecondsKey   = "autoscaling_scale_down_unready_time_seconds"
	kubeNodePoolAutoscalingScaleDownUtilizationThresholdKey = "autoscaling_scale_down_utilization_threshold"
	kubeNodePoolAntiAffinityKey                             = "anti_affinity"
	kubeNodePoolMonthlyBilledKey                            = "monthly_billed"
	kubeNodePoolAvailableNodesKey                           = "available_nodes"
	kubeNodePoolCurrentNodesKey                             = "current_nodes"
	kubeNodePoolUpToDateNodesKey                            = "up_to_date_nodes"
	kubeNodePoolSizeStatusKey                               = "size_status"
	kubeNodePoolTemplateKey                                 = "template"
	kubeNodePoolAvailabilityZonesKey                        = "availability_zones"
	kubeNodePoolAttachFloatingIpsKey                        = "attach_floating_ips"

	// template sub-attributes
	kubeNodePoolTemplateMetadataKey      = "metadata"
	kubeNodePoolTemplateSpecKey          = "spec"
	kubeNodePoolTemplateFinalizersKey    = "finalizers"
	kubeNodePoolTemplateLabelsKey        = "labels"
	kubeNodePoolTemplateAnnotationsKey   = "annotations"
	kubeNodePoolTemplateUnschedulableKey = "unschedulable"
	kubeNodePoolTemplateTaintsKey        = "taints"
	kubeNodePoolTemplateTaintKeyKey      = "key"
	kubeNodePoolTemplateTaintEffectKey   = "effect"
	kubeNodePoolTemplateTaintValueKey    = "value"
)

// OIDC attributes
const (
	kubeOidcClientIdKey       = "client_id"
	kubeOidcIssuerUrlKey      = "issuer_url"
	kubeOidcUsernameClaimKey  = "oidc_username_claim"
	kubeOidcUsernamePrefixKey = "oidc_username_prefix"
	kubeOidcGroupsClaimKey    = "oidc_groups_claim"
	kubeOidcGroupsPrefixKey   = "oidc_groups_prefix"
	kubeOidcRequiredClaimKey  = "oidc_required_claim"
	kubeOidcSigningAlgsKey    = "oidc_signing_algs"
	kubeOidcCaContentKey      = "oidc_ca_content"
)

// Log Subscription attributes
const (
	kubeLogSubscriptionIdKey           = "subscription_id"
	kubeLogSubscriptionKindKey         = "kind"
	kubeLogSubscriptionStreamIdKey     = "stream_id"
	kubeLogSubscriptionResourceKey     = "resource"
	kubeLogSubscriptionResourceNameKey = "name"
	kubeLogSubscriptionResourceTypeKey = "type"
)

// IP Restrictions attributes
const (
	kubeIpRestrictionsIpsKey = "ips"
)

// Node attributes (used in nodes data sources)
const (
	kubeNodeDeployedAtKey = "deployed_at"
	kubeNodeIdKey         = "id"
	kubeNodeInstanceIdKey = "instance_id"
	kubeNodeIsUpToDateKey = "is_up_to_date"
	kubeNodePoolIdKey     = "node_pool_id"
)
