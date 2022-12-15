package ovh

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

var TestAccDedicatedNASHAPlan = [][]map[string]interface{}{
	{
		{"Size": 10, "Protocol": "NFS", "IP": "149.202.78.2/32", "Type": "readonly"},
		// {"Size": 20, "Protocol": "CIFS", "IP": "149.202.78.2/32", "Type": "readonly"},
	},
}

type TestAccDedicatedNASHAPartition struct {
	ServiceName     string
	Name            string `json:"partitionName"`
	Protocol        string `json:"protocol"`
	Size            int    `json:"size"`
	Capacity        int    `json:"partitionCapacity"`
	UsedBySnapshots int    `json:"usedBySnapshots"`
}

type TestAccDedicatedNASHAPartitionAccess struct {
	ServiceName   string  `json:"serviceName"`
	PartitionName string  `json:"partitionName"`
	IP            string  `json:"ip"`
	Type          *string `json:"displayName"`
}

type TestAccDedicatedNASHAWrapper struct {
	RandomName              string
	ExpectedPartition       *TestAccDedicatedNASHAPartition
	ExpectedPartitionAccess *TestAccDedicatedNASHAPartitionAccess
}

func (w *TestAccDedicatedNASHAWrapper) Config() string {
	var config bytes.Buffer
	config.WriteString(fmt.Sprintf(`
	resource "ovh_dedicated_nasha_partition" "testacc" {
		service_name = "%s"
		name = "%s"
		protocol = "%s"
		size = %d
	}
	resource "ovh_dedicated_nasha_partition_access" "testacc" {
		service_name = "${ovh_dedicated_nasha_partition.testacc.service_name}"
		partition_name = "${ovh_dedicated_nasha_partition.testacc.name}"
		ip = "%s"
	`,
		w.ExpectedPartition.ServiceName,
		w.ExpectedPartition.Name,
		w.ExpectedPartition.Protocol,
		w.ExpectedPartition.Size,
		w.ExpectedPartitionAccess.IP))
	helpers.ConditionalAttributeString(&config, "type", w.ExpectedPartitionAccess.Type)
	config.WriteString(`}`)
	return config.String()
}

func (original *TestAccDedicatedNASHAPartition) MustEqual(compared *TestAccDedicatedNASHAPartition) error {
	if original.Name != compared.Name {
		return fmt.Errorf("Name differs ('%s' vs '%s')", original.Name, compared.Name)
	}
	if original.Protocol != compared.Protocol {
		return fmt.Errorf("Protocol differs ('%s' vs '%s')", original.Protocol, compared.Protocol)
	}
	if original.Size != compared.Size {
		return fmt.Errorf("Size differs ('%d' vs '%d')", original.Size, compared.Size)
	}
	if original.Capacity != compared.Capacity {
		return fmt.Errorf("Capacity differs ('%d' vs '%d')", original.Capacity, compared.Capacity)
	}
	if original.UsedBySnapshots != compared.UsedBySnapshots {
		return fmt.Errorf("UsedBySnapshots differs ('%d' vs '%d')", original.UsedBySnapshots, compared.UsedBySnapshots)
	}
	return nil
}

func (original *TestAccDedicatedNASHAPartitionAccess) MustEqual(compared *TestAccDedicatedNASHAPartitionAccess) error {
	if original.ServiceName != compared.ServiceName {
		return fmt.Errorf("ServiceName differs ('%s' vs '%s')", original.ServiceName, compared.ServiceName)
	}
	if original.PartitionName != compared.PartitionName {
		return fmt.Errorf("PartitionName differs ('%s' vs '%s')", original.PartitionName, compared.PartitionName)
	}
	if original.IP != compared.IP {
		return fmt.Errorf("ServiceName differs ('%s' vs '%s')", original.IP, compared.IP)
	}
	if !reflect.DeepEqual(original.Type, compared.Type) {
		return fmt.Errorf("Type differs ('%s' vs '%s')", *original.Type, *compared.Type)
	}
	return nil
}

// Returns assembled terraform TestStep from test wrapper
func (w *TestAccDedicatedNASHAWrapper) TestStep(c map[string]interface{}) resource.TestStep {
	w.ExpectedPartition.ServiceName = os.Getenv("OVH_NASHA_SERVICE")
	w.ExpectedPartitionAccess.ServiceName = os.Getenv("OVH_NASHA_SERVICE")

	if val, ok := c["Protocol"]; ok {
		w.ExpectedPartition.Protocol = val.(string)
	}
	if val, ok := c["Size"]; ok {
		w.ExpectedPartition.Size = val.(int)
	}
	if val, ok := c["IP"]; ok {
		w.ExpectedPartitionAccess.IP = val.(string)
	}
	if val, ok := c["Type"]; ok {
		w.ExpectedPartitionAccess.Type = helpers.GetNilStringPointer(val)
	}

	// set OVH API defaults instead of nil before checking
	if w.ExpectedPartitionAccess.Type == nil {
		val := "readwrite"
		w.ExpectedPartitionAccess.Type = &val
	}

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

		// // Check if resource states exists in local state at all
		// name = "ovh_dedicated_nasha_partition_access.testacc"
		// partitionAccessResource, ok := state.RootModule().Resources[name]
		// if !ok {
		// 	return fmt.Errorf("Not found: %s", name)
		// }
		// // Check if nasha_partition_access state fetched from the server reflects the expected state
		// partitionAccessResponse := &TestAccDedicatedNASHAPartitionAccess{}
		// config = testAccProvider.Meta().(*Config)
		// endpoint = fmt.Sprintf("/dedicated/nasha/%s/partition/%s", os.Getenv("OVH_NASHA_SERVICE"), partitionAccessResource.Primary.Attributes["name"])
		// err = config.OVHClient.Get(endpoint, partitionResponse)
		// if err != nil {
		// 	return fmt.Errorf("calling GET %s :\n\t %s", endpoint, err.Error())
		// }
		// err = expectedPartitionAccess.MustEqual(partitionAccessResponse)
		// if err != nil {
		// 	return fmt.Errorf("%s %s state differs from expected : %s", name, partitionAccessResource.Primary.Attributes["name"], err.Error())
		// }

		return nil
	}
}

func (w *TestAccDedicatedNASHAWrapper) TestDestroy(state *terraform.State) error {
	leftovers := false

	for _, resource := range state.RootModule().Resources {
		if resource.Type == "ovh_dedicated_nasha_partition_access" {
			config := testAccProvider.Meta().(*Config)
			service := os.Getenv("OVH_NASHA_SERVICE")
			partition := w.ExpectedPartitionAccess.PartitionName
			ip := url.PathEscape(resource.Primary.Attributes["ip"])
			endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/access/%s", service, partition, ip)
			err := config.OVHClient.Get(endpoint, nil)
			if err == nil {
				leftovers = true
			}
		}
	}

	for _, resource := range state.RootModule().Resources {
		if resource.Type == "ovh_dedicated_nasha_partition" {
			config := testAccProvider.Meta().(*Config)
			service := os.Getenv("OVH_NASHA_SERVICE")
			partition := w.ExpectedPartition.Name
			// ip := url.PathEscape(resource.Primary.Attributes["ip"])
			endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s", service, partition)
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
	w := &TestAccDedicatedNASHAWrapper{}
	hash := md5.Sum([]byte(time.Now().String()))
	w.RandomName = "testacc_" + hex.EncodeToString(hash[:])
	w.ExpectedPartition = &TestAccDedicatedNASHAPartition{ServiceName: os.Getenv("OVH_NASHA_SERVICE"), Name: w.RandomName}
	w.ExpectedPartitionAccess = &TestAccDedicatedNASHAPartitionAccess{ServiceName: os.Getenv("OVH_NASHA_SERVICE"), PartitionName: w.RandomName}
	return w
}

func TestAccDedicatedNASHA(t *testing.T) {
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
