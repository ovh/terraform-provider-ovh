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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("ovh_cloud_project_kube_nodepool", &resource.Sweeper{
		Name: "ovh_cloud_project_kube_nodepool",
		F:    testSweepCloudProjectKubeNodePool,
	})
}

func testSweepCloudProjectKubeNodePool(region string) error {
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

		pools := make([]CloudProjectKubeNodePoolResponse, 0)
		if err := client.Get(fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", serviceName, kubeId), &pools); err != nil {
			return fmt.Errorf("Error calling GET /cloud/project/%s/kube/%s/nodepool:\n\t %q", serviceName, kubeId, err)
		}

		if len(pools) == 0 {
			log.Print("[DEBUG] No pool to sweep")
			return nil
		}

		for _, p := range pools {
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				if err := client.Delete(fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, p.Id), nil); err != nil {
					return resource.RetryableError(err)
				}
				// Successful delete
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var testAccCloudProjectKubeNodePoolConfig = `
resource "ovh_cloud_project_kube" "cluster" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  version      = "%s"
}

resource "ovh_cloud_project_kube_nodepool" "pool" {
  service_name  = ovh_cloud_project_kube.cluster.service_name
  kube_id       = ovh_cloud_project_kube.cluster.id
  name          = ovh_cloud_project_kube.cluster.name
  flavor_name   = "b2-7"
  desired_nodes = 1
  min_nodes     = 0
  max_nodes     = 1
  template {
    metadata {
      annotations = {
        a1 = "av1"
      }
      finalizers = ["finalizer.extensions/v1beta1"]
      labels = {
        l1 = "lv1"
      }
    }
    spec {
      unschedulable = false
      taints = [
        {
          effect = "PreferNoSchedule"
          key    = "t1"
          value  = "tv1"
        }
      ]
    }
  }
}

`

var testAccCloudProjectKubeNodePoolConfigUpdated = `
resource "ovh_cloud_project_kube" "cluster" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  version      = "%s"
}

resource "ovh_cloud_project_kube_nodepool" "pool" {
  service_name  = ovh_cloud_project_kube.cluster.service_name
  kube_id       = ovh_cloud_project_kube.cluster.id
  name          = ovh_cloud_project_kube.cluster.name
  flavor_name   = "b2-7"
  desired_nodes = 2
  min_nodes     = 0
  max_nodes     = 2
  template {
    metadata {
      annotations = {
        a2 = "av2"
      }
      finalizers = []
      labels = {
        l2 = "lv2"
      }
    }
    spec {
      unschedulable = false
      taints = []
    }
  }
}

`

var testAccCloudProjectKubeNodePoolConfigUpdatedScaleToZero = `
resource "ovh_cloud_project_kube" "cluster" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  version      = "%s"
}

resource "ovh_cloud_project_kube_nodepool" "pool" {
  service_name  = ovh_cloud_project_kube.cluster.service_name
  kube_id       = ovh_cloud_project_kube.cluster.id
  name          = ovh_cloud_project_kube.cluster.name
  flavor_name   = "b2-7"
  desired_nodes = 0
  min_nodes     = 0
  max_nodes     = 2
  template {
    metadata {
      annotations = {
        a2 = "av2"
      }
      finalizers = []
      labels = {
        l2 = "lv2"
      }
    }
    spec {
      unschedulable = false
      taints = []
    }
  }
}

`

func TestAccCloudProjectKubeNodePool(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")
	version := os.Getenv("OVH_CLOUD_PROJECT_KUBE_VERSION_TEST")
	resourceName := "ovh_cloud_project_kube_nodepool.pool"
	config := fmt.Sprintf(
		testAccCloudProjectKubeNodePoolConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		version,
	)
	configUpdated := fmt.Sprintf(
		testAccCloudProjectKubeNodePoolConfigUpdated,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
		version,
	)
	configUpdatedScaleToZero := fmt.Sprintf(
		testAccCloudProjectKubeNodePoolConfigUpdatedScaleToZero,
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
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "version", version),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "flavor_name", "b2-7"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "desired_nodes", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "min_nodes", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "max_nodes", "1"),

					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.metadata.0.annotations.a1", "av1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.metadata.0.finalizers.0", "finalizer.extensions/v1beta1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.metadata.0.labels.l1", "lv1"),

					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.spec.0.taints.0.effect", "PreferNoSchedule"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.spec.0.taints.0.key", "t1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.spec.0.taints.0.value", "tv1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.spec.0.unschedulable", "false"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "version", version),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "flavor_name", "b2-7"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "desired_nodes", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "min_nodes", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "max_nodes", "2"),

					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.metadata.0.annotations.a2", "av2"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.metadata.0.finalizers.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.metadata.0.labels.l2", "lv2"),

					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.spec.0.taints.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.spec.0.unschedulable", "false"),
				),
			},
			{
				Config: configUpdatedScaleToZero,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_kube.cluster", "kubeconfig"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube.cluster", "version", version),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "flavor_name", "b2-7"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "desired_nodes", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "min_nodes", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "max_nodes", "2"),

					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.metadata.0.annotations.a2", "av2"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.metadata.0.finalizers.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.metadata.0.labels.l2", "lv2"),

					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.spec.0.taints.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_project_kube_nodepool.pool", "template.0.spec.0.unschedulable", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					poolId := state.RootModule().Resources[resourceName].Primary.ID
					kubernetesClusterID := state.RootModule().Resources["ovh_cloud_project_kube.cluster"].Primary.ID
					return fmt.Sprintf("%s/%s/%s", os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), kubernetesClusterID, poolId), nil
				},
			},
		},
	})
}
