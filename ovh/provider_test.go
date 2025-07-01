package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

var testAccProviders map[string]*schema.Provider
var testAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)

var testAccProvider *schema.Provider
var testAccOVHClient *ovhwrap.Client

func init() {
	log.SetOutput(os.Stdout)
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"ovh": testAccProvider,
	}

	// Instantiate a composite provider
	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		context.Background(),
		Provider().GRPCProvider, // Provider using terraform-plugin-sdk
	)
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(&OvhProvider{}), // Provider using terraform-plugin-framework
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatal(err)
	}

	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"ovh": func() (tfprotov6.ProviderServer, error) {
			return muxServer, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func checkEnvOrFail(t *testing.T, e string) {
	if os.Getenv(e) == "" {
		t.Fatalf("%s must be set for acceptance tests", e)
	}
}

func checkEnvOrSkip(t *testing.T, e string) {
	if os.Getenv(e) == "" {
		t.Skipf("[WARN] %s must be set for acceptance tests. Skipping.", e)
	}
}

// Checks that the environment variables needed to create the OVH API client
// are set and create the client right away.
func testAccPreCheckCredentials(t *testing.T) {
	if testAccOVHClient == nil {
		config := Config{
			Endpoint:          os.Getenv("OVH_ENDPOINT"),
			ApplicationKey:    os.Getenv("OVH_APPLICATION_KEY"),
			ApplicationSecret: os.Getenv("OVH_APPLICATION_SECRET"),
			ConsumerKey:       os.Getenv("OVH_CONSUMER_KEY"),
			ClientID:          os.Getenv("OVH_CLIENT_ID"),
			ClientSecret:      os.Getenv("OVH_CLIENT_SECRET"),
			ApiRateLimit:      ratelimit.NewUnlimited(),
		}

		if err := config.load(); err != nil {
			t.Fatalf("Couldn't load OVH Client: %s", err)
		} else {
			testAccOVHClient = config.OVHClient
		}
	}
}

// Checks that the environment variables needed for the /ip acceptance tests
// are set.
func testAccPreCheckIp(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_IP_TEST")
	checkEnvOrSkip(t, "OVH_IP_BLOCK_TEST")
	checkEnvOrSkip(t, "OVH_IP_REVERSE_TEST")
}

// Checks that the environment variables needed for the /ip/move acceptance tests
// are set.
func testAccPreCheckIpMove(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_IP_MOVE_SERVICE_NAME_TEST")
}

// Checks that the environment variables needed to order /ip/service for acceptance tests
// are set.
func testAccPreCheckOrderIpService(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_TESTACC_ORDER_IP_SERVICE")
}

// Checks that the environment variables needed to order /cloud/project for acceptance tests
// are set.
func testAccPreCheckOrderCloudProject(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_TESTACC_ORDER_CLOUD_PROJECT")
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
}

// Checks that the environment variables needed for the /domain acceptance tests
// are set.
func testAccPreCheckDomain(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_ZONE_TEST")
}

// Checks that the environment variables needed to order /domain for acceptance tests
// are set.
func testAccPreCheckOrderDomain(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_TESTACC_ORDER_DOMAIN")
}

// Checks that the environment variables needed for the /hosting/privatedatabase acceptance tests
// are set.
func testAccPreCheckHostingPrivateDatabase(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
}

// Checks that the environment variables needed for the /hosting/privatedatabase acceptance tests
// are set.
func testAccPreCheckHostingPrivateDatabaseDatabase(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")
}

// Checks that the environment variables needed for the /hosting/privatedatabase acceptance tests
// are set.
func testAccPreCheckHostingPrivateDatabaseUser(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_USER_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_PASSWORD_TEST")
}

// Checks that the environment variables needed for the /hosting/privatedatabase acceptance tests
// are set.
func testAccPreCheckHostingPrivateDatabaseUserGrant(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_NAME_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_USER_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_PASSWORD_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_GRANT_TEST")
}

// Checks that the environment variables needed for the /hosting/privatedatabase acceptance tests
// are set.
func testAccPreCheckHostingPrivateDatabaseWhitelist(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_ENGINE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_DC_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_WHITELIST_IP_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_WHITELIST_NAME_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_HOSTING_PRIVATEDATABASE_WHITELIST_SFTP_TEST")
}

// Checks that the environment variables needed for the /dbaas acceptance tests
// are set.
func testAccPreCheckDbaasLogs(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_SERVICE_TEST")
}

func testAccPreCheckDbaasLogsInput(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_LOGSTASH_VERSION_TEST")
}

func testAccPreCheckDbaasLogsCluster(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_CLUSTER_ID")
}

func testAccPreCheckDbaasLogsClusterRetention(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_CLUSTER_ID")
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_CLUSTER_RETENTION_ID")
}

// Checks that the environment variables needed for the /cloud acceptance tests
// are set.
func testAccPreCheckCloud(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SERVICE_TEST")
}

// Checks that the environment variables needed for the /cloud/region/storage acceptance tests
// are set.
func testAccPreCheckCloudStorage(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_BUCKET_NAME_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_OBJECT_TEST")
}

// Checks that the environment variables needed for the /publicCloud/project/{id}/rancher acceptance tests
// are set.
func testAccPreCheckRancher(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_RANCHER_TEST")
}

// Checks that the environment variables needed for the /cloud/{cloudId}/containerregistry acceptance tests
// are set.
func testAccPreCheckContainerRegistry(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudProjectExists(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_CONTAINERREGISTRY_REGION_TEST")
}

// Checks that the environment variables needed for the /cloud/{cloudId}/containerregistry/{registryID}/openIdConnect acceptance tests
// are set.
func testAccPreCheckContainerRegistryOIDC(t *testing.T) {
	testAccPreCheckContainerRegistry(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_CONTAINERREGISTRY_OIDC_ENDPOINT_TEST")
}

// Checks that the environment variables needed for the /cloud/{cloudId}/containerregistry/{registryID}/iam acceptance tests
// are set.
func testAccPreCheckContainerRegistryIAM(t *testing.T) {
	testAccPreCheckContainerRegistry(t)
}

// Checks that the environment variables needed for the /cloud/project/{projectId}/network/private/ acceptance tests are set.
func testAccPreCheckCloudNetworkPrivate(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudProjectExists(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_PRIVATE_NETWORK_TEST")
}

// Checks that the environment variables needed for the /cloud/project/{projectId}/database/ acceptance tests are set.
func testAccPreCheckCloudDatabase(t *testing.T) {
	testAccPreCheckCloudDatabaseNoEngine(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_DATABASE_ENGINE_TEST")
}

// Checks that the environment variables needed for the /cloud/project/{projectId}/database/<engine>/ acceptance tests are set.
func testAccPreCheckCloudDatabaseNoEngine(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudProjectExists(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_DATABASE_VERSION_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_DATABASE_REGION_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_DATABASE_FLAVOR_TEST")
}

// Checks that the environment variables needed for the /cloud/project/{projectId}/database/mongodb/ acceptance tests are set.
func testAccPreCheckCloudDatabaseMongoDBNoEngine(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudProjectExists(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_DATABASE_MONGODB_VERSION_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_DATABASE_MONGODB_REGION_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_DATABASE_MONGODB_FLAVOR_TEST")
}

// Checks that the environment variables needed for the /cloud/project/{projectId}/database/{engine}/{clusterId}/ipRestriction/ acceptance tests are set.
func testAccPreCheckCloudDatabaseIpRestriction(t *testing.T) {
	testAccPreCheckCloudDatabase(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_DATABASE_IP_RESTRICTION_IP_TEST")
}

func testAccPreCheckCloudRegion(t *testing.T) {
	testAccPreCheckCloud(t)
	testAccCheckCloudProjectExists(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_REGION_TEST")
}

func testAccPreCheckCloudRegionLoadbalancer(t *testing.T) {
	testAccPreCheckCloudRegion(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_LOADBALANCER_TEST")
}

// Checks that the environment variables needed for the /cloud/project/{projectId}/ip/failover acceptance tests
// are set.
func testAccPreCheckFailoverIpAttach(t *testing.T) {
	testAccPreCheckCredentials(t)
	testAccPreCheckCloud(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FAILOVER_IP_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_1_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FAILOVER_IP_ROUTED_TO_2_TEST")
}

// Checks that the environment variables needed for the /cloud/{cloudId}/kube acceptance tests
// are set.
func testAccPreCheckKubernetes(t *testing.T) {
	testAccPreCheckCredentials(t)
	testAccPreCheckCloud(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_KUBE_VERSION_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_KUBE_PREV_VERSION_TEST")
}

// Checks that the environment variables needed for the /vrack/{service}/cloudProject acceptance tests
// are set.
func testAccPreCheckKubernetesVRack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
}

// Checks that the environment variables needed for the /vrack/{service}/ovhCloudConnect/{ovhCloudConnect} acceptance tests
// are set.
func testAccPreCheckOCCVRack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_VRACK_OVH_CLOUD_CONNECT")
}

// Checks that the environment variables needed for the /vrack/{service}/dedicatedCloud/{dedicatedCloud} acceptance tests
// are set.
func testAccPreCheckDedicatedCloud(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_VRACK_DEDICATED_CLOUD")
}

// Checks that the environment variables needed for the /vrack/{service}/dedicatedCloudDatacenter/{datacenter} acceptance tests
// are set.
func testAccPreCheckDedicatedCloudDatacenter(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_VRACK_DEDICATED_CLOUD_DATACENTER")
	checkEnvOrSkip(t, "OVH_VRACK_TARGET_SERVICE_TEST")
}

// Checks that the environment variables needed for the /vrack/{service}/ipv6/{ipv6} acceptance tests
// are set.
func testAccPreCheckIPv6VRack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_IP_V6_BLOCK_TEST")
}

// Checks that the environment variables needed for the /vrack/{service}/ipv6/{ipv6} import acceptance tests
// are set.
func testAccPreCheckIPv6ImportVRack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_IP_V6_BLOCK_IMPORT_TEST")
}

// Checks that the environment variables needed for the /vrack/{service}/ipv6/{ipv6}/routedSubrange{subrange} acceptance tests
// are set.
func testAccPreCheckIPv6RoutedSubrangeVRack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_IP_V6_BLOCK_TEST")
	checkEnvOrSkip(t, "OVH_IP_V6_ROUTED_SUBRANGE_TEST")
	checkEnvOrSkip(t, "OVH_IP_V6_ROUTED_SUBRANGE_NEXTHOP_TEST")
}

// Checks that the environment variables needed for the /vrack/{service}/ipv6/{ipv6}/routedSubrange{subrange} import acceptance tests
// are set.
func testAccPreCheckIPv6RoutedSubrangeImportVRack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_IP_V6_BLOCK_TEST")
	checkEnvOrSkip(t, "OVH_IP_V6_ROUTED_SUBRANGE_IMPORT_TEST")
	checkEnvOrSkip(t, "OVH_IP_V6_ROUTED_SUBRANGE_NEXTHOP_TEST")
}

// Checks that the environment variables needed for the /ipLoadbalacing acceptance tests
// are set.
func testAccPreCheckIpLoadbalancing(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_IPLB_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_IPLB_IPFO_TEST")
}

// Checks that the environment variables needed to order /ipLoadbalacing for acceptance tests
// are set.
func testAccPreCheckOrderIpLoadbalancing(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_TESTACC_ORDER_IPLOADBALANCING")
}

// Checks that the environment variables needed to order /vps for acceptance tests
// are set.
func testAccPreCheckOrderVPS(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_TESTACC_ORDER_VPS")
}

// Checks that the environment variables needed to order /vrack for acceptance tests
// are set.
func testAccPreCheckOrderVrack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_TESTACC_ORDER_VRACK")
}

// Checks that the environment variables needed for the /vrack acceptance tests
// are set.
func testAccPreCheckVRack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
}

// Checks that the environment variables needed for the /vrack/{serviceName}/ip acceptance tests
// are set, and variable used for block with region.
func testAccPreCheckVRackIpWithRegion(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_VRACK_IP_BLOCK")
	checkEnvOrSkip(t, "OVH_VRACK_IP_REGION")
}

// Checks that the environment variables needed for the /me/paymentMean acceptance tests
// are set.
func testAccPreCheckMePaymentMean(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_TEST_BANKACCOUNT")
}

func testAccPreCheckDedicatedServer(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_DEDICATED_SERVER")
}

func testAccPreCheckOrderDedicatedServer(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_TESTACC_ORDER_DEDICATED_SERVER")
}

func testAccPreCheckVPS(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VPS")
	checkEnvOrSkip(t, "OVH_VPS_IMAGE_ID")
}

func testAccPreCheckOkms(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_OKMS")
}

func testAccPreCheckOkmsCredential(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_OKMS")
	checkEnvOrSkip(t, "OVH_OKMS_CREDENTIAL")
}

func testAccPreCheckCloudInstance(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_NETWORK_PRIVATE_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_NETWORK_PRIVATE_SUBNET_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_FLOATING_IP_ID")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_GATEWAY_ID")
}

func testAccCheckVRackExists(t *testing.T) {
	type vrackResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	r := vrackResponse{}

	endpoint := fmt.Sprintf("/vrack/%s", url.PathEscape(os.Getenv("OVH_VRACK_SERVICE_TEST")))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}
}

func testAccCheckCloudProjectExists(t *testing.T) {
	type cloudProjectResponse struct {
		ID          string `json:"project_id"`
		Status      string `json:"status"`
		Description string `json:"description"`
	}

	r := cloudProjectResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s", url.PathEscape(os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}
	t.Logf("Read Cloud Project %s -> status: '%s', desc: '%s'", endpoint, r.Status, r.Description)
}

func testAccCheckIpLoadbalancingExists(t *testing.T) {
	type iplbResponse struct {
		ServiceName string `json:"serviceName"`
		State       string `json:"state"`
	}

	r := iplbResponse{}

	endpoint := fmt.Sprintf("/ipLoadbalancing/%s", url.PathEscape(os.Getenv("OVH_IPLB_SERVICE_TEST")))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}
	t.Logf("Read IPLB service %s -> state: '%s', serviceName: '%s'", endpoint, r.State, r.ServiceName)
}

func testAccCheckDomainZoneExists(t *testing.T) {
	type domainZoneResponse struct {
		NameServers []string `json:"nameServers"`
	}

	r := domainZoneResponse{}

	endpoint := fmt.Sprintf("/domain/zone/%s", url.PathEscape(os.Getenv("OVH_ZONE_TEST")))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}

	t.Logf("Read Domain Zone %s -> nameservers: '%v'", endpoint, r.NameServers)
}

func testAccPreCheckDedicatedCeph(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_DEDICATED_CEPH")
}

func testAccPreCheckWorkflowBackup(t *testing.T) {
	checkEnvOrSkip(t, WORKFLOW_BACKUP_TEST_INSTANCE_ID_ENV_VAR)
	checkEnvOrSkip(t, WORKFLOW_BACKUP_TEST_REGION_ENV_VAR)
}

// This variable shall be defined to run the test because it targets an internal route that shall be authorized per user
func testAccPreCheckDedicatedServerNetworking(t *testing.T) {
	checkEnvOrSkip(t, "TEST_DEDICATED_SERVER_NETWORKING")
}

func testAccPreCheckIamResourceGroup(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_VPS")
}

// Checks that the environment variables needed ofr the /storage/netapp acceptance tests
// are set.
func testAccPreCheckStorageEfs(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_STORAGE_EFS_SERVICE_TEST")
}

func testAccCheckStorageEfsExists(t *testing.T) {
	type efsServiceResponse struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	r := efsServiceResponse{}

	endpoint := fmt.Sprintf("/storage/netapp/%s", url.PathEscape(os.Getenv("OVH_STORAGE_EFS_SERVICE_TEST")))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}
	t.Logf("Read Storage EFS service %s -> state: '%s', serviceName: %s", endpoint, r.Status, r.ID)
}
