package ovh

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var TestAccDedicatedNASHAPlan = [][]map[string]interface{}{
	{
		{"Name": "testacc1", "Size": 10, "Protocol": "NFS", "IP": "10.10.10.10", "Mode": "Readonly"},
	},
	{
		{"Name": "testacc2"},
	},
}

type TestAccDedicatedNASHAPartition struct {
	ServiceName string
	Name        string
	Protocol    string
	Size        int
}

type TestAccDedicatedNASHAPartitionAccess struct {
	ServiceName   string
	PartitionName string
	IP            string
	Mode          *string `json:"displayName"`
}

type TestAccDedicatedNASHAWrapper struct {
	ExpectedPartition       *TestAccDedicatedNASHAPartition
	ExpectedPartitionAccess *TestAccDedicatedNASHAPartitionAccess
}

func (w *TestAccDedicatedNASHAWrapper) Config() string {
	var config bytes.Buffer
	config.WriteString(fmt.Sprintf(`
    resource "ovh_dedicated_nasha_partition" "testacc" {
	  service_name = "%s"
	  name = "testacc"
	  protocol = "NFS"
	  size = 10
	}

	resource "ovh_dedicated_nasha_partition_access" "testacc" {
	  service_name = "%s"
	  partition_name = ""
	  ip = ""`, w.ExpectedPartition.ServiceName, w.ExpectedPartitionAccess.PartitionName))
	conditionalAttributeString(&config, "mode", w.ExpectedPartitionAccess.Mode)
	config.WriteString(`}`)
	return config.String()
}

func (original *TestAccDedicatedNASHAPartition) MustEqual(compared *TestAccDedicatedNASHAPartition) error {
	if !reflect.DeepEqual(original.Name, compared.Name) {
		return fmt.Errorf("Name differs")
	}
	return nil
}

func (original *TestAccDedicatedNASHAPartitionAccess) MustEqual(compared *TestAccDedicatedNASHAPartitionAccess) error {
	if !reflect.DeepEqual(original.ServiceName, compared.ServiceName) {
		return fmt.Errorf("ServiceName differs")
	}
	return nil
}

// type TestAccDedicatedNASHAStep struct {
// 	Response *TestAccDedicatedNASHA
// 	Expected *TestAccDedicatedNASHA
// }

// Returns assembled terraform TestStep from test wrapper
func (w *TestAccDedicatedNASHAWrapper) TestStep(c map[string]interface{}) resource.TestStep {
	// if val, ok := c["Name"]; ok {
	// 	w.ExpectedPartition.Name = getNilStringPointer(val)
	// }
	// if val, ok := c["Address"]; ok {
	// 	w.ExpectedPartition.Address = getNilStringPointer(val)
	// }
	// if val, ok := c["Port"]; ok {
	// 	w.ExpectedPartitionAccess.IP = getNilIntPointer(val)
	// }

	// expected := *w.Expected

	// set OVH API defaults instead of nil before checking
	// if expected.Probe == nil {
	// 	val := false
	// 	expected.Probe = &val
	// }
	// if expected.Ssl == nil {
	// 	val := false
	// 	expected.Ssl = &val
	// }
	// if expected.Backup == nil {
	// 	val := false
	// 	expected.Backup = &val
	// }
	// if expected.Weight == nil {
	// 	val := 1
	// 	expected.Weight = &val
	// }

	return resource.TestStep{
		Config: w.Config(),
		Check: resource.ComposeTestCheckFunc(
			w.TestCheck(w.ExpectedPartition, w.ExpectedPartitionAccess),
		),
	}
}

func (w *TestAccDedicatedNASHAWrapper) TestCheck(expectedPartition *TestAccDedicatedNASHAPartition, expectedPartitionAccess *TestAccDedicatedNASHAPartitionAccess) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// Check if resource states exists in local state at all
		name := "ovh_dedicated_nasha_partition.testacc"
		partitionResource, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		// Check if nasha partition state fetched from the server reflects the expected state
		partitionResponse := &TestAccDedicatedNASHAPartition{}
		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s", os.Getenv("OVH_NASHA_SERVICE"), partitionResource.Primary.Attributes["name"])
		err := config.OVHClient.Get(endpoint, partitionResponse)
		if err != nil {
			return fmt.Errorf("calling GET %s :\n\t %s", endpoint, err.Error())
		}
		err = expectedPartition.MustEqual(partitionResponse)
		if err != nil {
			return fmt.Errorf("%s %s state differs from expected : %s", name, partitionResource.Primary.Attributes["name"], err.Error())
		}

		// Check if resource states exists in local state at all
		name = "ovh_dedicated_nasha_partition_access.testacc"
		partitionAccessResource, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		// Check if nasha_partition_access state fetched from the server reflects the expected state
		partitionAccessResponse := &TestAccDedicatedNASHAPartitionAccess{}
		config = testAccProvider.Meta().(*Config)
		endpoint = fmt.Sprintf("/dedicated/nasha/%s/partition/%s", os.Getenv("OVH_NASHA_SERVICE"), partitionAccessResource.Primary.Attributes["name"])
		err = config.OVHClient.Get(endpoint, partitionResponse)
		if err != nil {
			return fmt.Errorf("calling GET %s :\n\t %s", endpoint, err.Error())
		}
		err = expectedPartitionAccess.MustEqual(partitionAccessResponse)
		if err != nil {
			return fmt.Errorf("%s %s state differs from expected : %s", name, partitionAccessResource.Primary.Attributes["name"], err.Error())
		}

		return nil
	}
}

func (w *TestAccDedicatedNASHAWrapper) TestDestroy(state *terraform.State) error {
	leftovers := false
	// Removing partition will also remove all its defined access rules
	for _, resource := range state.RootModule().Resources {
		if resource.Type == "ovh_dedicated_nasha_partition" {
			config := testAccProvider.Meta().(*Config)
			service := os.Getenv("OVH_NASHA_SERVICE")
			partition := w.ExpectedPartition.Name
			ip := url.PathEscape(resource.Primary.Attributes["ip"])
			endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/access/%s", service, partition, ip)
			err := config.OVHClient.Get(endpoint, nil)
			if err == nil {
				leftovers = true
			}
		}
	}
	if leftovers {
		return fmt.Errorf("ovh_dedicated_nasha_partition still exists")
	}
	return nil
}

func newTestAccDedicatedNASHAWrapper() *TestAccDedicatedNASHAWrapper {
	return &TestAccDedicatedNASHAWrapper{
		ExpectedPartition:       &TestAccDedicatedNASHAPartition{ServiceName: os.Getenv("OVH_NASHA_SERVICE")},
		ExpectedPartitionAccess: &TestAccDedicatedNASHAPartitionAccess{ServiceName: os.Getenv("OVH_NASHA_SERVICE")},
	}
}

func TestAccDedicatedNASHABasicCreate(t *testing.T) {
	for _, plan := range TestAccDedicatedNASHAPlan {
		w := newTestAccDedicatedNASHAWrapper()
		var steps []resource.TestStep
		for _, tcase := range plan {
			steps = append(steps, w.TestStep(tcase))
		}
		resource.Test(t, resource.TestCase{
			Providers:    testAccProviders,
			CheckDestroy: w.TestDestroy,
			Steps:        steps,
		})
	}
}
