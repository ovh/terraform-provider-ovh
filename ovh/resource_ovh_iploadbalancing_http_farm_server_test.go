package ovh

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var TestAccIpLoadbalancingHttpFarmServerPlan = [][]map[string]interface{}{
	{
		{
			"Status":      "active",
			"Address":     "10.0.0.11",
			"Port":        80,
			"Weight":      3,
			"DisplayName": "testBackendA",
		},
		{
			"Status":      "active",
			"Address":     "10.0.0.11",
			"Port":        8080,
			"Weight":      3,
			"DisplayName": "testBackendA",
			"Probe":       true,
			"Backup":      true,
		},
		{
			"Status":      "active",
			"Address":     "10.0.0.11",
			"Port":        8080,
			"Weight":      2,
			"DisplayName": "testBackendB",
			"Probe":       false,
			"Backup":      false,
		},
	},
	{
		{
			"Status":  "inactive",
			"Address": "10.0.0.12",
			"Port":    80,
		},
		{
			"Status":               "active",
			"Address":              "10.0.0.11",
			"Port":                 8080,
			"ProxyProtocolVersion": "v2",
			"Ssl":                  true,
		},
		{
			"Status":               "active",
			"Address":              "10.0.0.11",
			"Port":                 8080,
			"ProxyProtocolVersion": "v1",
			"Ssl":                  true,
			"Backup":               false,
		},
		{
			"Status":               "active",
			"Address":              "10.0.0.11",
			"Port":                 8080,
			"ProxyProtocolVersion": nil,
			"Ssl":                  true,
			"Backup":               true,
		},
	},
}

type TestAccIpLoadbalancingHttpFarmServer struct {
	ServiceName          string
	ServerId             int     `json:"serverId"`
	BackendId            int     `json:"backendId"`
	FarmId               int     `json:"farmId"`
	DisplayName          *string `json:"displayName"`
	Address              string  `json:"address"`
	Cookie               *string `json:"cookie"`
	Port                 *int    `json:"port"`
	ProxyProtocolVersion *string `json:"proxyProtocolVersion"`
	Chain                *string `json:"chain"`
	Weight               *int    `json:"weight"`
	Probe                *bool   `json:"probe"`
	Ssl                  *bool   `json:"ssl"`
	Backup               *bool   `json:"backup"`
	Status               string  `json:"status"`
}

type TestAccIpLoadbalancingHttpFarmServerWrapper struct {
	Expected *TestAccIpLoadbalancingHttpFarmServer
}

func (w *TestAccIpLoadbalancingHttpFarmServerWrapper) Config() string {
	var config bytes.Buffer

	config.WriteString(fmt.Sprintf(`
    resource "ovh_iploadbalancing_http_farm" "testacc" {
	  service_name = "%s"
	  display_name = "testacc"
	  port = 8080
	  zone = "all"
	  probe {
	    port = 8080
	    interval = 30
	    type = "http"
	  }
	}
	resource "ovh_iploadbalancing_http_farm_server" "testacc" {
	  service_name = "%s"
	  farm_id = "${ovh_iploadbalancing_http_farm.testacc.id}"
	  address = "%s"
	  status = "%s"
	`, w.Expected.ServiceName,
		w.Expected.ServiceName,
		w.Expected.Address,
		w.Expected.Status,
	))

	conditionalAttributeString(&config, "display_name", w.Expected.DisplayName)
	conditionalAttributeInt(&config, "port", w.Expected.Port)
	conditionalAttributeString(&config, "proxy_protocol_version", w.Expected.ProxyProtocolVersion)
	conditionalAttributeInt(&config, "weight", w.Expected.Weight)
	conditionalAttributeBool(&config, "probe", w.Expected.Probe)
	conditionalAttributeBool(&config, "ssl", w.Expected.Ssl)
	conditionalAttributeBool(&config, "backup", w.Expected.Backup)
	config.WriteString(`}`)
	return config.String()
}

func (server *TestAccIpLoadbalancingHttpFarmServer) MustEqual(compared *TestAccIpLoadbalancingHttpFarmServer) error {
	if !reflect.DeepEqual(server.DisplayName, compared.DisplayName) {
		return fmt.Errorf("DisplayName differs")
	}
	if !reflect.DeepEqual(server.Address, compared.Address) {
		return fmt.Errorf("Address differs")
	}
	if !reflect.DeepEqual(server.Port, compared.Port) {
		return fmt.Errorf("Port differs")
	}
	if !reflect.DeepEqual(server.ProxyProtocolVersion, compared.ProxyProtocolVersion) {
		return fmt.Errorf("ProxyProtocolVersion differs")
	}
	if !reflect.DeepEqual(server.Chain, compared.Chain) {
		return fmt.Errorf("Chain differs")
	}
	if !reflect.DeepEqual(server.Weight, compared.Weight) {
		return fmt.Errorf("Weight differs")
	}
	if !reflect.DeepEqual(server.Probe, compared.Probe) {
		return fmt.Errorf("Probe differs")
	}
	if !reflect.DeepEqual(server.Ssl, compared.Ssl) {
		return fmt.Errorf("Ssl differs")
	}
	if !reflect.DeepEqual(server.Backup, compared.Backup) {
		return fmt.Errorf("Backup differs")
	}
	if !reflect.DeepEqual(server.Status, compared.Status) {
		return fmt.Errorf("Status differs")
	}
	return nil
}

type TestAccIpLoadbalancingHttpFarmServerStep struct {
	Response *TestAccIpLoadbalancingHttpFarmServer
	Expected *TestAccIpLoadbalancingHttpFarmServer
}

func (w *TestAccIpLoadbalancingHttpFarmServerWrapper) TestStep(c map[string]interface{}) resource.TestStep {
	w.Expected.DisplayName = getNilStringPointerFromData(c, "DisplayName")
	w.Expected.Address = c["Address"].(string)
	w.Expected.Port = getNilIntPointerFromData(c, "Port")
	w.Expected.ProxyProtocolVersion = getNilStringPointerFromData(c, "ProxyProtocolVersion")
	w.Expected.Chain = getNilStringPointerFromData(c, "Chain")
	w.Expected.Weight = getNilIntPointerFromData(c, "Weight")
	w.Expected.Probe = getNilBoolPointerFromData(c, "Probe")
	w.Expected.Ssl = getNilBoolPointerFromData(c, "Ssl")
	w.Expected.Backup = getNilBoolPointerFromData(c, "Backup")
	w.Expected.Status = c["Status"].(string)

	expected := *w.Expected

	// set OVH API defaults instead of nil before checking
	if expected.Probe == nil {
		val := false
		expected.Probe = &val
	}
	if expected.Ssl == nil {
		val := false
		expected.Ssl = &val
	}
	if expected.Backup == nil {
		val := false
		expected.Backup = &val
	}
	if expected.Weight == nil {
		val := 1
		expected.Weight = &val
	}

	return resource.TestStep{
		Config: w.Config(),
		Check: resource.ComposeTestCheckFunc(
			w.TestCheck(expected),
		),
	}
}

func (w *TestAccIpLoadbalancingHttpFarmServerWrapper) TestCheck(expected TestAccIpLoadbalancingHttpFarmServer) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		response := &TestAccIpLoadbalancingHttpFarmServer{}
		name := "ovh_iploadbalancing_http_farm_server.testacc"
		resource, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/farm/%s/server/%s", os.Getenv("OVH_IPLB_SERVICE"), resource.Primary.Attributes["farm_id"], resource.Primary.ID)
		err := config.OVHClient.Get(endpoint, response)
		if err != nil {
			return fmt.Errorf("calling GET %s :\n\t %s", endpoint, err.Error())
		}

		err = expected.MustEqual(response)
		if err != nil {
			return fmt.Errorf("%s %s state differs from expected : %s", name, resource.Primary.ID, err.Error())
		}
		return nil
	}
}

func (w *TestAccIpLoadbalancingHttpFarmServerWrapper) TestDestroy(state *terraform.State) error {
	leftovers := false
	for _, resource := range state.RootModule().Resources {
		if resource.Type != "ovh_iploadbalancing_http_farm_server" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/http/farm/%d/server/%s", os.Getenv("OVH_IPLB_SERVICE"), w.Expected.FarmId, resource.Primary.ID)
		err := config.OVHClient.Get(endpoint, nil)
		if err == nil {
			leftovers = true
		}
	}
	if leftovers {
		return fmt.Errorf("IpLoadbalancing farm still exists")
	}
	return nil
}

func newTestAccIpLoadbalancingHttpFarmServerWrapper() *TestAccIpLoadbalancingHttpFarmServerWrapper {
	return &TestAccIpLoadbalancingHttpFarmServerWrapper{
		Expected: &TestAccIpLoadbalancingHttpFarmServer{ServiceName: os.Getenv("OVH_IPLB_SERVICE")},
	}
}

func TestAccIpLoadbalancingHttpFarmServerBasicCreate(t *testing.T) {
	for _, plan := range TestAccIpLoadbalancingHttpFarmServerPlan {
		w := newTestAccIpLoadbalancingHttpFarmServerWrapper()
		var steps []resource.TestStep
		for _, tcase := range plan {
			steps = append(steps, w.TestStep(tcase))
		}
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheckIpLoadbalancing(t) },
			Providers:    testAccProviders,
			CheckDestroy: w.TestDestroy,
			Steps:        steps,
		})
	}
}
