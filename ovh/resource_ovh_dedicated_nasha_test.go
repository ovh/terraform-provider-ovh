package ovh

import (
	"bytes"
	"fmt"
)

// var TestAccDedicatedNASHAPlan = [][]map[string]interface{}{
// 	{
// 		{"Status": "active", "Address": "10.0.0.11", "Port": 80, "Weight": 3, "DisplayName": "testBackendA"},
// 		{"Port": 8080, "Probe": true, "Backup": true},
// 		{"Port": 8080, "Probe": false, "Backup": false, "Weight": 2, "DisplayName": "testBackendB"},
// 	},
// 	{
// 		{"Status": "inactive", "Address": "10.0.0.12", "Port": 80},
// 		{"Port": 8080, "ProxyProtocolVersion": "v2", "Ssl": true},
// 		{"Port": 8080, "ProxyProtocolVersion": "v1", "Ssl": true, "Backup": false},
// 		{"Port": 8080, "ProxyProtocolVersion": nil, "Ssl": true, "Backup": true, "Status": "active"},
// 	},
// }

// type TestAccDedicatedNASHA struct {
// 	ServiceName          string
// 	ServerId             int     `json:"serverId"`
// 	BackendId            int     `json:"backendId"`
// 	FarmId               int     `json:"farmId"`
// 	DisplayName          *string `json:"displayName"`
// 	Address              *string `json:"address"`
// 	Cookie               *string `json:"cookie"`
// 	Port                 *int    `json:"port"`
// 	ProxyProtocolVersion *string `json:"proxyProtocolVersion"`
// 	Chain                *string `json:"chain"`
// 	Weight               *int    `json:"weight"`
// 	Probe                *bool   `json:"probe"`
// 	Ssl                  *bool   `json:"ssl"`
// 	Backup               *bool   `json:"backup"`
// 	Status               *string `json:"status"`
// }

// type TestAccDedicatedNASHAWrapper struct {
// 	Expected *TestAccDedicatedNASHA
// }

func (w *TestAccDedicatedNASHAWrapper) Config() string {
	var config bytes.Buffer
	config.WriteString(fmt.Sprintf(`
    resource "ovh_dedicated_nasha_partition" "testacc" {
	  service_name = "%s"
	  name = "testacc"
	}

	resource "ovh_dedicated_nasha_partition_access" "testacc" {
	  service_name = "%s"
	  partition_name = ""
	  farm_id = "${ovh_iploadbalancing_tcp_farm.testacc.id}"
	  address = "%s"
	  status = "%s"
	`, w.Expected.ServiceName,
		w.Expected.ServiceName,
		*w.Expected.Address,
		*w.Expected.Status))
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

// func (server *TestAccDedicatedNASHA) MustEqual(compared *TestAccDedicatedNASHA) error {
// 	if !reflect.DeepEqual(server.DisplayName, compared.DisplayName) {
// 		return fmt.Errorf("DisplayName differs")
// 	}
// 	if !reflect.DeepEqual(server.Address, compared.Address) {
// 		return fmt.Errorf("Address differs")
// 	}
// 	if !reflect.DeepEqual(server.Port, compared.Port) {
// 		return fmt.Errorf("Port differs")
// 	}
// 	if !reflect.DeepEqual(server.ProxyProtocolVersion, compared.ProxyProtocolVersion) {
// 		return fmt.Errorf("ProxyProtocolVersion differs")
// 	}
// 	if !reflect.DeepEqual(server.Chain, compared.Chain) {
// 		return fmt.Errorf("Chain differs")
// 	}
// 	if !reflect.DeepEqual(server.Weight, compared.Weight) {
// 		return fmt.Errorf("Weight differs")
// 	}
// 	if !reflect.DeepEqual(server.Probe, compared.Probe) {
// 		return fmt.Errorf("Probe differs")
// 	}
// 	if !reflect.DeepEqual(server.Ssl, compared.Ssl) {
// 		return fmt.Errorf("Ssl differs")
// 	}
// 	if !reflect.DeepEqual(server.Backup, compared.Backup) {
// 		return fmt.Errorf("Backup differs")
// 	}
// 	if !reflect.DeepEqual(server.Status, compared.Status) {
// 		return fmt.Errorf("Status differs")
// 	}
// 	return nil
// }

// type TestAccDedicatedNASHAStep struct {
// 	Response *TestAccDedicatedNASHA
// 	Expected *TestAccDedicatedNASHA
// }

// func (w *TestAccDedicatedNASHAWrapper) TestStep(c map[string]interface{}) resource.TestStep {
// 	if val, ok := c["DisplayName"]; ok {
// 		w.Expected.DisplayName = getNilStringPointer(val)
// 	}
// 	if val, ok := c["Address"]; ok {
// 		w.Expected.Address = getNilStringPointer(val)
// 	}
// 	if val, ok := c["Port"]; ok {
// 		w.Expected.Port = getNilIntPointer(val)
// 	}
// 	if val, ok := c["ProxyProtocolVersion"]; ok {
// 		w.Expected.ProxyProtocolVersion = getNilStringPointer(val)
// 	}
// 	if val, ok := c["Chain"]; ok {
// 		w.Expected.Chain = getNilStringPointer(val)
// 	}
// 	if val, ok := c["Weight"]; ok {
// 		w.Expected.Weight = getNilIntPointer(val)
// 	}
// 	if val, ok := c["Probe"]; ok {
// 		w.Expected.Probe = getNilBoolPointer(val)
// 	}
// 	if val, ok := c["Ssl"]; ok {
// 		w.Expected.Ssl = getNilBoolPointer(val)
// 	}
// 	if val, ok := c["Backup"]; ok {
// 		w.Expected.Backup = getNilBoolPointer(val)
// 	}
// 	if val, ok := c["Status"]; ok {
// 		w.Expected.Status = getNilStringPointer(val)
// 	}

// 	expected := *w.Expected

// 	// set OVH API defaults instead of nil before checking
// 	if expected.Probe == nil {
// 		val := false
// 		expected.Probe = &val
// 	}
// 	if expected.Ssl == nil {
// 		val := false
// 		expected.Ssl = &val
// 	}
// 	if expected.Backup == nil {
// 		val := false
// 		expected.Backup = &val
// 	}
// 	if expected.Weight == nil {
// 		val := 1
// 		expected.Weight = &val
// 	}

// 	return resource.TestStep{
// 		Config: w.Config(),
// 		Check: resource.ComposeTestCheckFunc(
// 			w.TestCheck(expected),
// 		),
// 	}
// }

// func (w *TestAccDedicatedNASHAWrapper) TestCheck(expected TestAccDedicatedNASHA) resource.TestCheckFunc {
// 	return func(state *terraform.State) error {
// 		response := &TestAccDedicatedNASHA{}
// 		name := "ovh_iploadbalancing_tcp_farm_server.testacc"
// 		resource, ok := state.RootModule().Resources[name]
// 		if !ok {
// 			return fmt.Errorf("Not found: %s", name)
// 		}
// 		config := testAccProvider.Meta().(*Config)
// 		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%s/server/%s", os.Getenv("OVH_IPLB_SERVICE"), resource.Primary.Attributes["farm_id"], resource.Primary.ID)
// 		err := config.OVHClient.Get(endpoint, response)
// 		if err != nil {
// 			return fmt.Errorf("calling GET %s :\n\t %s", endpoint, err.Error())
// 		}

// 		err = expected.MustEqual(response)
// 		if err != nil {
// 			return fmt.Errorf("%s %s state differs from expected : %s", name, resource.Primary.ID, err.Error())
// 		}
// 		return nil
// 	}
// }

// func (w *TestAccDedicatedNASHAWrapper) TestDestroy(state *terraform.State) error {
// 	leftovers := false
// 	for _, resource := range state.RootModule().Resources {
// 		if resource.Type != "ovh_iploadbalancing_tcp_farm_server" {
// 			continue
// 		}

// 		config := testAccProvider.Meta().(*Config)
// 		endpoint := fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d/server/%s", os.Getenv("OVH_IPLB_SERVICE"), w.Expected.FarmId, resource.Primary.ID)
// 		err := config.OVHClient.Get(endpoint, nil)
// 		if err == nil {
// 			leftovers = true
// 		}
// 	}
// 	if leftovers {
// 		return fmt.Errorf("IpLoadbalancing farm still exists")
// 	}
// 	return nil
// }

// func newTestAccDedicatedNASHAWrapper() *TestAccDedicatedNASHAWrapper {
// 	return &TestAccDedicatedNASHAWrapper{
// 		Expected: &TestAccDedicatedNASHA{ServiceName: os.Getenv("OVH_IPLB_SERVICE")},
// 	}
// }

// func TestAccDedicatedNASHABasicCreate(t *testing.T) {
// 	for _, plan := range TestAccDedicatedNASHAPlan {
// 		w := newTestAccDedicatedNASHAWrapper()
// 		var steps []resource.TestStep
// 		for _, tcase := range plan {
// 			steps = append(steps, w.TestStep(tcase))
// 		}
// 		resource.Test(t, resource.TestCase{
// 			Providers:    testAccProviders,
// 			CheckDestroy: w.TestDestroy,
// 			Steps:        steps,
// 		})
// 	}
// }
