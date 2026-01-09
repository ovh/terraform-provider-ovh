package ovh

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"go.uber.org/ratelimit"
)

func getFlavorAndImage(project, region string) (string, string, error) {
	client, err := clientDefault(&Config{
		ApiRateLimit: ratelimit.NewUnlimited(),
	})
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

func TestAccCloudProjectInstance_privateNetworkCreate(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	networkName := acctest.RandomWithPrefix(test_prefix)
	vlanID := rand.Intn(4000)

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
							vlan_id = %d
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
		networkName,
		vlanID)

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

func TestAccCloudProjectInstance_privateNetworkAlreadyExists(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	networkID := os.Getenv("OVH_CLOUD_PROJECT_NETWORK_PRIVATE_TEST")
	subnetID := os.Getenv("OVH_CLOUD_PROJECT_NETWORK_PRIVATE_SUBNET_TEST")
	floatingIpID := os.Getenv("OVH_CLOUD_PROJECT_FLOATING_IP_ID")
	gatewayID := os.Getenv("OVH_CLOUD_PROJECT_GATEWAY_ID")

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
						floating_ip {
							id = "%s"
						}
						network {
							id = "%s"
							subnet_id = "%s"
						}
						gateway {
							id = "%s"
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
		floatingIpID, networkID, subnetID, gatewayID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckCloudInstance(t)
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

func TestAccCloudProjectInstance_multiNIC(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	rName := acctest.RandomWithPrefix(test_prefix)
	netName1 := fmt.Sprintf("%s_net1", rName)
	netName2 := fmt.Sprintf("%s_net2", rName)

	vlanID1 := rand.Intn(2000) + 2
	vlanID2 := vlanID1 + 1

	var testCreateInstanceMultiNIC = fmt.Sprintf(`
		resource "private" "net1" {
			service_name = "%s"
			name         = "%s"
			regions      = ["%s"]
			vlan_id      = %d
		}

		resource "ovh_cloud_project_subnet" "sub1" {
			service_name = "%s"
			network_id   = private.net1.id
			region       = "%s"
			network      = "192.168.1.0/24"
			dhcp         = true
			start        = "192.168.1.100"
			end          = "192.168.1.200"
			no_gateway   = false
		}

		resource "private" "net2" {
			service_name = "%s"
			name         = "%s"
			regions      = ["%s"]
			vlan_id      = %d
		}

		resource "ovh_cloud_project_subnet" "sub2" {
			service_name = "%s"
			network_id   = private.net2.id
			region       = "%s"
			network      = "192.168.2.0/24"
			dhcp         = true
			start        = "192.168.2.100"
			end          = "192.168.2.200"
			no_gateway   = false
		}

		resource "ovh_cloud_project_instance" "multi_nic" {
			service_name = "%s"
			region       = "%s"
			name         = "%s"
            billing_period = "hourly"
			
			flavor {
				flavor_id = "%s"
			}
            boot_from {
                image_id = "%s"
            }

			network {
                # Interface Publique
				public = true

                # Interface Privée 1
				private {
					network {
						id        = private.net1.id
						subnet_id = ovh_cloud_project_subnet.sub1.id
					}
				}

                # Interface Privée 2
				private {
					network {
						id        = private.net2.id
						subnet_id = ovh_cloud_project_subnet.sub2.id
					}
				}
			}
		}
	`,
		serviceName, netName1, region, vlanID1,
		serviceName, region,
		serviceName, netName2, region, vlanID2,
		serviceName, region,
		serviceName, region, rName, flavor, image,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateInstanceMultiNIC,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.multi_nic", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_project_instance.multi_nic", "network.0.private.#", "2"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.multi_nic", "addresses.0.ip"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.multi_nic", "addresses.1.ip"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_instance.multi_nic", "addresses.2.ip"),
				),
			},
		},
	})
}
