package ovh

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ovh/go-ovh/ovh"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccOVHClient *ovh.Client

func init() {
	log.SetOutput(os.Stdout)
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"ovh": testAccProvider,
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
	checkEnvOrFail(t, "OVH_ENDPOINT")
	checkEnvOrFail(t, "OVH_APPLICATION_KEY")
	checkEnvOrFail(t, "OVH_APPLICATION_SECRET")
	checkEnvOrFail(t, "OVH_CONSUMER_KEY")

	if testAccOVHClient == nil {
		config := Config{
			Endpoint:          os.Getenv("OVH_ENDPOINT"),
			ApplicationKey:    os.Getenv("OVH_APPLICATION_KEY"),
			ApplicationSecret: os.Getenv("OVH_APPLICATION_SECRET"),
			ConsumerKey:       os.Getenv("OVH_CONSUMER_KEY"),
		}

		if err := config.loadAndValidate(); err != nil {
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
	checkEnvOrSkip(t, "OVH_IP")
	checkEnvOrSkip(t, "OVH_IP_BLOCK")
	checkEnvOrSkip(t, "OVH_IP_REVERSE")
}

// Checks that the environment variables needed for the /domain acceptance tests
// are set.
func testAccPreCheckDomain(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_ZONE")
}

// Checks that the environment variables needed for the /cloud acceptance tests
// are set.
func testAccPreCheckCloud(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_PUBLIC_CLOUD")
}

// Checks that the environment variables needed for the /ipLoadbalacing acceptance tests
// are set.
func testAccPreCheckIpLoadbalancing(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_IPLB_SERVICE")
}

// Checks that the environment variables needed for the /vrack acceptance tests
// are set.
func testAccPreCheckVRack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VRACK")
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

func testAccPreCheckVPS(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_VPS")
}

func testAccCheckVRackExists(t *testing.T) {
	type vrackResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	r := vrackResponse{}

	endpoint := fmt.Sprintf("/vrack/%s", os.Getenv("OVH_VRACK"))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}
}

func testAccCheckCloudExists(t *testing.T) {
	type cloudProjectResponse struct {
		ID          string `json:"project_id"`
		Status      string `json:"status"`
		Description string `json:"description"`
	}

	r := cloudProjectResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s", os.Getenv("OVH_PUBLIC_CLOUD"))

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

	endpoint := fmt.Sprintf("/ipLoadbalancing/%s", os.Getenv("OVH_IPLB_SERVICE"))

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

	endpoint := fmt.Sprintf("/domain/zone/%s", os.Getenv("OVH_ZONE"))

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
