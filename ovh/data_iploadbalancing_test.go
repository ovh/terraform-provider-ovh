package ovh

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpLoadbalancingDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE_TEST")
	config := fmt.Sprintf(testAccIpLoadbalancingDatasourceConfig_Basic, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "id", serviceName),
					resource.TestCheckResourceAttrSet(
						"data.ovh_iploadbalancing.iplb", "urn"),
				),
			},
		},
	})
}

func TestAccIpLoadbalancingDataSource_statevrack(t *testing.T) {
	serviceName := os.Getenv("OVH_IPLB_SERVICE_TEST")
	config := fmt.Sprintf(testAccIpLoadbalancingDatasourceConfig_StateAndVrack, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "id", serviceName),
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "state", "ok"),
					resource.TestCheckResourceAttr(
						"data.ovh_iploadbalancing.iplb", "vrack_eligibility", "true"),
					resource.TestCheckResourceAttrSet(
						"data.ovh_iploadbalancing.iplb", "urn"),
				),
			},
		},
	})
}

const testAccIpLoadbalancingDatasourceConfig_Basic = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
}
`
const testAccIpLoadbalancingDatasourceConfig_StateAndVrack = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
  state = "ok"
  vrack_eligibility = true
}
`

func TestIpLoadbalancingDataSourceRead_zoneFilter(t *testing.T) {
	const serviceName = "loadbalancer-test"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/time":
			fmt.Fprint(w, "1700000000")
		case "/ipLoadbalancing/" + serviceName:
			fmt.Fprint(w, `{"serviceName":"`+serviceName+`","ipLoadbalancing":"192.0.2.10","offer":"lb1","state":"ok","zone":["gra","rbx"]}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	tests := []struct {
		name    string
		zones   []interface{}
		wantErr bool
	}{
		{name: "same zones", zones: []interface{}{"gra", "rbx"}},
		{name: "same zones in another order", zones: []interface{}{"rbx", "gra"}},
		{name: "subset of zones", zones: []interface{}{"gra"}, wantErr: true},
		{name: "other zone", zones: []interface{}{"sbg"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, dataSourceIpLoadbalancing().Schema, map[string]interface{}{
				"service_name": serviceName,
				"zone":         tt.zones,
			})

			err := dataSourceIpLoadbalancingRead(d, newRefreshTestConfig(t, server.URL))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected no result for a zone filter that does not match")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if d.Id() != serviceName {
				t.Fatalf("expected id %q, got %q", serviceName, d.Id())
			}
		})
	}
}
