package ovh

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

const testCloudProjectRegionStoragePresign = `
resource "ovh_cloud_project_region_storage_presign" "presign_url" {
  service_name = "%s"
  region_name  = "%s"
  name         = "%s"
  expire       = 3600
  method       = "GET"
  object       = "%s"
}
`

func testAccPreCheckCloudRegionStorage(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_BUCKET_NAME_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_OBJECT_TEST")
}

func testAccCheckCloudRegionStorage(t *testing.T) {

	type cloudProjectRegionStorageResponse struct {
		Name   string `json:"name"`
		Region string `json:"region"`
	}

	r := cloudProjectRegionStorageResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s/region/%s/storage/%s",
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"),
		os.Getenv("OVH_CLOUD_PROJECT_STORAGE_BUCKET_NAME_TEST"))

	err := testAccOVHClient.Get(endpoint, &r)
	if err != nil {
		t.Fatalf("Error: %q\n", err)
	}
	t.Logf("Read Storage Container %s -> name: '%s', region: '%s'", endpoint, r.Name, r.Region)
}

func TestCloudProjectRegionStoragePresign(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regionName := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")
	name := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_BUCKET_NAME_TEST")
	object := os.Getenv("OVH_CLOUD_PROJECT_STORAGE_OBJECT_TEST")

	config := fmt.Sprintf(testCloudProjectRegionStoragePresign, serviceName, regionName, name, object)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudRegionStorage(t)
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccCheckCloudRegionStorage(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_region_storage_presign.presign_url", "url"),
				),
			},
		},
	})
}
