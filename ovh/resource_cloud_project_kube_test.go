package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_cloud_project_kube", &resource.Sweeper{
		Name: "ovh_cloud_project_kube",
		Dependencies: []string{
			"ovh_cloud_project_kube_nodepool",
		},
		F: testSweepCloudProjectKube,
	})
}

func testSweepCloudProjectKube(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST is not set. No kube to sweep")
		return nil
	}

	kubeIds := make([]string, 0)
	if err := client.Get(fmt.Sprintf("/cloud/project/%s/kube", serviceName), &kubeIds); err != nil {
		return fmt.Errorf("Error calling GET /cloud/project/%s/kube:\n\t %q", serviceName, err)
	}
	for _, kubeId := range kubeIds {
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, kubeId)
		res := &CloudProjectKubeResponse{}
		log.Printf("[DEBUG] read kube %s from project: %s", kubeId, serviceName)
		if err := client.Get(endpoint, res); err != nil {
			return err
		}
		if !strings.HasPrefix(res.Name, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, kubeId), nil); err != nil {
				return resource.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}

	}
	return nil
}

var testAccCloudProjectKubeConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
    name          = "%s"
	region        = "%s"
	version = "%s"
}
`

var testAccCloudProjectKubeEmptyVersionConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
    name          = "%s"
	region        = "%s"
}
`

func TestAccCloudProjectKube_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_KUBE_VERSION_TEST")
	config := fmt.Sprintf(
		testAccCloudProjectKubeConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		version,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "version", version),
				),
			},
		},
	})
}

// TestAccCloudProjectKubeEmptyVersion_basic
// create a public cluster
// check some properties
// update cluster name
// check some properties && cluster updated name
func TestAccCloudProjectKubeEmptyVersion_basic(t *testing.T) {
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")

	name := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccCloudProjectKubeEmptyVersionConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
	)

	updatedName := acctest.RandomWithPrefix(test_prefix)
	updatedConfig := fmt.Sprintf(
		testAccCloudProjectKubeEmptyVersionConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		updatedName,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", kubeClusterNameKey, name),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "version"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_kube.cluster", kubeClusterNameKey, updatedName),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_kube.cluster", "version"),
				),
			},
		},
	})
}
