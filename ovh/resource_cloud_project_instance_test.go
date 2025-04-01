package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func getFlavorAndImage(project, region string) (string, string, error) {
	client, err := clientDefault(&Config{})
	if err != nil {
		return "", "", fmt.Errorf("error getting client: %w", err)
	}

	type ResponseStruct struct {
		ID     string `json:"id"`
		Type   string `json:"type"`
		OSType string `json:"osType"`
		Name   string `json:"name"`
	}

	endpoint := fmt.Sprintf("/cloud/project/%s/flavor?region=%s", url.PathEscape(project), url.QueryEscape(region))

	var response []*ResponseStruct
	if err := client.Get(endpoint, &response); err != nil {
		return "", "", fmt.Errorf("failed to get flavors: %w", err)
	}

	for _, flav := range response {
		if flav.Name != "b2-7" {
			continue
		}

		endpoint = fmt.Sprintf("/cloud/project/%s/image?region=%s&osType=%s&flavorType=%s",
			url.PathEscape(project),
			url.QueryEscape(region),
			url.QueryEscape(flav.OSType),
			url.QueryEscape(flav.Type),
		)

		var images []*ResponseStruct
		if err := client.Get(endpoint, &images); err != nil {
			return "", "", fmt.Errorf("failed to get images: %w", err)
		}

		if len(images) > 0 {
			return flav.ID, images[0].ID, nil
		}
	}

	return "", "", fmt.Errorf("found no flavor and image for project %s and region %s", project, region)
}

func TestAccCloudProjectInstance_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	var testCreateInstance = fmt.Sprintf(`
			resource "ovh_cloud_project_instance" "instance" {
				service_name = "%s"
				region = "%s"
				billing_period = "hourly"
				boot_from {
					image_id = "%s"
				}
				flavor {
					flavor_id = "%s"
				}
				name = "TestInstance"
				ssh_key {
					name = "%s"
				}
				network {
					public = true
				}
			}
		`,
		serviceName,
		region,
		image,
		flavor,
		os.Getenv("OVH_CLOUD_PROJECT_SSH_NAME_TEST"))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.instance", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_name", "b2-7"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_id", flavor),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "image_id", image),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "name", "TestInstance"),
				),
			},
		},
	})
}

func TestAccCloudProjectInstance3AZ_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_3AZ_REGION_TEST")
	az := os.Getenv("OVH_CLOUD_PROJECT_AVAILABILITY_ZONE_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	var testCreateInstance = fmt.Sprintf(`
			resource "ovh_cloud_project_instance" "instance" {
				service_name = "%s"
				region = "%s"
				billing_period = "hourly"
				boot_from {
					image_id = "%s"
				}
				flavor {
					flavor_id = "%s"
				}
				name = "TestInstance"
				ssh_key {
					name = "%s"
				}
				network {
					public = true
				}
				availability_zone = "%s"
			}
		`,
		serviceName,
		region,
		image,
		flavor,
		os.Getenv("OVH_CLOUD_PROJECT_SSH_NAME_TEST"),
		az)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.instance", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_name", "b3-8"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_id", flavor),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "image_id", image),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "name", "TestInstance"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "availability_zone", az),
				),
			},
		},
	})
}

func TestAccCloudProjectInstance_withSSHKeyCreate(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	sshKeyName := acctest.RandomWithPrefix(test_prefix)

	var testCreateInstance = fmt.Sprintf(`
			resource "ovh_cloud_project_instance" "instance" {
				service_name = "%s"
				region = "%s"
				billing_period = "hourly"
				boot_from {
					image_id = "%s"
				}
				flavor {
					flavor_id = "%s"
				}
				name = "TestInstance"
				ssh_key_create {
					name = "%s"
					public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC9xPpdqP3sx2H+gcBm65tJEaUbuifQ1uGkgrWtNY0PRKNNPdy+3yoVOtxk6Vjo4YZ0EU/JhmQfnrK7X7Q5vhqYxmozi0LiTRt0BxgqHJ+4hWTWMIOgr+C2jLx7ZsCReRk+fy5AHr6h0PHQEuXVLXeUy/TDyuY2JPtUZ5jcqvLYgQ== my-key"
				}
				network {
					public = true
				}
			}
		`,
		serviceName,
		region,
		image,
		flavor,
		sshKeyName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.instance", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_name", "b2-7"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_id", flavor),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "image_id", image),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "name", "TestInstance"),
				),
			},
		},
	})
}

func TestAccCloudProjectInstance_privateNetwork(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	networkName := acctest.RandomWithPrefix(test_prefix)

	var testCreateInstance = fmt.Sprintf(`
			resource "ovh_cloud_project_instance" "instance" {
				service_name = "%s"
				region = "%s"
				billing_period = "hourly"
				boot_from {
					image_id = "%s"
				}
				flavor {
					flavor_id = "%s"
				}
				name = "TestInstance"
				ssh_key {
					name = "%s"
				}
				network {
					private {
						network_create {
							name = "%s"
							vlan_id = 1237
							subnet{
								ip_version = 4
								cidr = "10.0.0.1/20"
								enable_dhcp = true
							}
						}
					}
				}
			}
		`,
		serviceName,
		region,
		image,
		flavor,
		os.Getenv("OVH_CLOUD_PROJECT_SSH_NAME_TEST"),
		networkName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.instance", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_name", "b2-7"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "flavor_id", flavor),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "image_id", image),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.instance", "name", "TestInstance"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.instance", "addresses.0.ip"),
				),
			},
		},
	})
}
