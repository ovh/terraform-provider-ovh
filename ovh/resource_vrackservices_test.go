package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	subnetDisplayName = "tf.test.subnet"
	cidr              = "192.168.0.0/24"
	serviceRangeCidr  = "192.168.0.0/29"
	vlan              = "30"

	subnetDisplayNameUpdated = "tf.test.subnet.updated"
	cidrUpdated              = "10.120.0.0/24"
	serviceRangeCidrUpdated  = "10.120.0.0/29"
)

const testAccStorageEfsConfig = `
	data "ovh_storage_efs" "tf-acc-efs" {
		service_name = "%s"
	}
`

// #1 - resourceName
// #2 - region
// #3 - targetSpec
const testAccVrackServicesConfig = `
	data "ovh_me" "myaccount" {}

	resource "ovh_vrackservices" "tf-acc-vrackservices" {
		ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
		plan = [
			{
				plan_code = "vrack-services"
				duration = "P1M"
				pricing_mode = "default"

				configuration = [
					{
						label = "region_name"
						value = "%s"
					}
				]
			}
		]
		
		%s
	}
`

// #1 - vrackServiceName
// #2 - vrackServicesServiceName
const testAccVrackVrackServicesBinding = `
	resource "ovh_vrack_vrackservices" "vrack-vrackservices-binding" {
		service_name   = "%s"
		vrack_services = ovh_vrackservices.tf-acc-vrackservices.id
	}
`

func testCheckAttrNull(resource, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("resource not found")
		}

		if _, ok := rs.Primary.Attributes[attr]; ok {
			return fmt.Errorf("expected %s to be null", attr)
		}
		return nil
	}
}

func TestAccResourceVrackServices_basic(t *testing.T) {
	region := os.Getenv("OVH_VRACK_SERVICES_REGION")
	vrackServiceName := os.Getenv("OVH_VRACK_SERVICE_TEST")
	storageEfsServiceName := os.Getenv("OVH_STORAGE_EFS_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckVrackServices(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1 - input validation
				Config: fmt.Sprintf(testAccStorageEfsConfig, storageEfsServiceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_storage_efs.tf-acc-efs", "region", region),
				),
			},
			{
				// Step 2 - order
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					`
						target_spec = {
							subnets = []
						}
					`,
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "0"),
				),
			},
			{
				// Step 3 - import resource to state
				ResourceName:            "ovh_vrackservices.tf-acc-vrackservices",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"plan", "ovh_subsidiary", "order"},

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "0"),
					resource.TestMatchResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "id", regexp.MustCompile(`^vrs-[a-z0-9-]+$`)),
				),
			},
			{
				// Step 4 - associate to a vRack
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					`
						target_spec = {
							subnets = []
						}
					`) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrack_vrackservices.vrack-vrackservices-binding", "service_name", vrackServiceName),
				),
			},
			{
				// Step 5 - update resource - add empty subnet
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									cidr         = "%s"
									display_name = "%s"
									service_range = {
										cidr = "%s"
									}
									service_endpoints = []
									vlan = %s
								},
							]
						}
					`, cidr, subnetDisplayName, serviceRangeCidr, vlan)) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayName),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_range.cidr", serviceRangeCidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan", vlan),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_endpoints.#", "0"),
				),
			},
			{
				// Step 6 - update resource - update subnet displayName, vlan, cidr and serviceRangeCidr
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									cidr         = "%s"
									display_name = "%s"
									service_range = {
										cidr = "%s"
									}
									service_endpoints = []
								},
							]
						}
					`, cidrUpdated, subnetDisplayNameUpdated, serviceRangeCidrUpdated)) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayNameUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_range.cidr", serviceRangeCidrUpdated),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_endpoints.#", "0"),
				),
			},
			{
				// Step 7 - update resource - add subnet.service_endpoint
				Config: fmt.Sprintf(testAccStorageEfsConfig, storageEfsServiceName) +
					fmt.Sprintf(testAccVrackServicesConfig, region,
						fmt.Sprintf(`
							target_spec = {
								subnets = [
									{
										cidr         = "%s"
										display_name = "%s"
										service_range = {
											cidr = "%s"
										}
										service_endpoints = [
											{
												managed_service_urn = data.ovh_storage_efs.tf-acc-efs.iam.urn
											}
										]
									},
								]
							}
						`, cidrUpdated, subnetDisplayNameUpdated, serviceRangeCidrUpdated),
					) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "ACTIVE"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayNameUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_range.cidr", serviceRangeCidrUpdated),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_endpoints.#", "1"),
					resource.TestCheckResourceAttr(
						"ovh_vrackservices.tf-acc-vrackservices",
						"target_spec.subnets.0.service_endpoints.0.managed_service_urn",
						fmt.Sprintf("urn:v1:eu:resource:storageNetApp:%s", storageEfsServiceName),
					),
				),
			},
			{
				// Step 8 - update resource - delete subnet.service_endpoint
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									cidr         = "%s"
									display_name = "%s"
									service_range = {
										cidr = "%s"
									}
									service_endpoints = []
								},
							]
						}
					`, cidrUpdated, subnetDisplayNameUpdated, serviceRangeCidrUpdated)) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayNameUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_range.cidr", serviceRangeCidrUpdated),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_endpoints.#", "0"),
				),
			},
			{
				// Step 9 - update resource - delete empty subnet
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					`
						target_spec = {
							subnets = []
						}
					`) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "0"),
				),
			},
			{
				// Step 10 - update resource - add subnet with service_endpoint
				Config: fmt.Sprintf(testAccStorageEfsConfig, storageEfsServiceName) +
					fmt.Sprintf(testAccVrackServicesConfig, region,
						fmt.Sprintf(`
							target_spec = {
								subnets = [
									{
										cidr         = "%s"
										display_name = "%s"
										service_range = {
											cidr = "%s"
										}
										service_endpoints = [
											{
												managed_service_urn = data.ovh_storage_efs.tf-acc-efs.iam.urn
											}
										]
									},
								]
							}
						`, cidr, subnetDisplayName, serviceRangeCidr),
					) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "ACTIVE"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayName),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_range.cidr", serviceRangeCidr),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_endpoints.#", "1"),
					resource.TestCheckResourceAttr(
						"ovh_vrackservices.tf-acc-vrackservices",
						"target_spec.subnets.0.service_endpoints.0.managed_service_urn",
						fmt.Sprintf("urn:v1:eu:resource:storageNetApp:%s", storageEfsServiceName),
					),
				),
			},
			{
				// Step 11 - dissociate from vRack
				Config: fmt.Sprintf(testAccStorageEfsConfig, storageEfsServiceName) +
					fmt.Sprintf(testAccVrackServicesConfig, region,
						fmt.Sprintf(`
							target_spec = {
								subnets = [
									{
										cidr         = "%s"
										display_name = "%s"
										service_range = {
											cidr = "%s"
										}
										service_endpoints = [
											{
												managed_service_urn = data.ovh_storage_efs.tf-acc-efs.iam.urn
											}
										]
									},
								]
							}
						`, cidr, subnetDisplayName, serviceRangeCidr),
					),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					// VrackServices resource state not updated yet, need a refresh
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "ACTIVE"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayName),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_range.cidr", serviceRangeCidr),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.service_endpoints.#", "1"),
					resource.TestCheckResourceAttr(
						"ovh_vrackservices.tf-acc-vrackservices",
						"target_spec.subnets.0.service_endpoints.0.managed_service_urn",
						fmt.Sprintf("urn:v1:eu:resource:storageNetApp:%s", storageEfsServiceName),
					),
				),
			},
			{
				// Step 12 - refresh to force state update on VrackService resource (current_state.product_status)
				PlanOnly:     true,
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "DRAFT"),
				),
			},
			{
				// Step 13 - update resource - delete subnet
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					`
						target_spec = {
							subnets = []
						}
					`) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.product_status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "0"),
				),
			},
		},
	})
}
