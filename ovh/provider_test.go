package ovh

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/ovh/go-ovh/ovh"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccOVHClient *ovh.Client

func init() {
	log.SetOutput(os.Stdout)
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ovh": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func checkEnv(t *testing.T, e string) {
	if os.Getenv(e) == "" {
		t.Fatalf("%s must be set for acceptance tests", e)
	}
}

// Checks that the environment variables needed to create the OVH API client
// are set and create the client right away.
func testAccPreCheckCredentials(t *testing.T) {
	checkEnv(t, "OVH_ENDPOINT")
	checkEnv(t, "OVH_APPLICATION_KEY")
	checkEnv(t, "OVH_APPLICATION_SECRET")
	checkEnv(t, "OVH_CONSUMER_KEY")

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
	checkEnv(t, "OVH_IP")
	checkEnv(t, "OVH_IP_BLOCK")
	checkEnv(t, "OVH_IP_REVERSE")
}

// Checks that the environment variables needed for the /domain acceptance tests
// are set.
func testAccPreCheckDomain(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnv(t, "OVH_ZONE")
}

// Checks that the environment variables needed for the /cloud acceptance tests
// are set.
func testAccPreCheckPublicCloud(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnv(t, "OVH_PUBLIC_CLOUD")
}

// Checks that the environment variables needed for the /ipLoadbalacing acceptance tests
// are set.
func testAccPreCheckIpLoadbalancing(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnv(t, "OVH_IPLB_SERVICE")
}

// Checks that the environment variables needed for the /vrack acceptance tests
// are set.
func testAccPreCheckVRack(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnv(t, "OVH_VRACK")
}

// Checks that the environment variables needed for the /me/paymentMean acceptance tests
// are set.
func testAccPreCheckMePaymentMean(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnv(t, "OVH_TEST_BANKACCOUNT")
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
	t.Logf("Read VRack %s -> name:'%s', desc:'%s' ", endpoint, r.Name, r.Description)

}

func testAccCheckPublicCloudExists(t *testing.T) {
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
