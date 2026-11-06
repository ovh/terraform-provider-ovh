package ovh

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

	// BGP range in cidr, past serviceRangeCidr; its first and last addresses are reserved
	bgpRangeCidr      = "192.168.0.16/28"
	bgpClientPeer     = "192.168.0.21"
	clientASN         = "64600"
	infrastructureASN = "65000"
)

// testAccVrackServicesBGPInfraPeers returns one infrastructure peer per edge device of the region:
// 4, or 2 on the LABEU lab platform.
func testAccVrackServicesBGPInfraPeers(region string) []string {
	peers := []string{"192.168.0.17", "192.168.0.18", "192.168.0.19", "192.168.0.20"}
	if strings.HasPrefix(region, "labeu-") {
		return peers[:2]
	}
	return peers
}

// testAccHCLStringList renders values as an HCL list of strings.
func testAccHCLStringList(values []string) string {
	return `["` + strings.Join(values, `", "`) + `"]`
}

const testAccStorageEfsConfig = `
	data "ovh_storage_efs" "tf-acc-efs" {
		service_name = "%s"
	}
`

// #1 - region
// #2 - plan_option and target_spec
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

const testAccVrackServicesBGPPlanOption = `
		plan_option = [
			{
				plan_code    = "bgp-beta-vrack-services"
				duration     = "P1M"
				pricing_mode = "default"
				quantity     = 1
			}
		]
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
	subnetAZ := os.Getenv("OVH_VRACK_SERVICES_SUBNET_AZ")
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
				Config: fmt.Sprintf(testAccVrackServicesConfig, region, `
						target_spec = {
							subnets = []
						}
					`),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "0"),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.id"),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.client_asn"),
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
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "DRAFT"),
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
				// Step 5 - update resource - add subnet with an empty service range
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az           = "%s"
									cidr         = "%s"
									display_name = "%s"
									vlan         = %s
								},
							]
							services = {
								service_range = [
									{
										cidr              = "%s"
										subnet_cidr       = "%s"
										service_endpoints = []
									},
								]
							}
						}
					`, subnetAZ, cidr, subnetDisplayName, vlan, serviceRangeCidr, cidr),
				) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.vrack_id", vrackServiceName),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.az", subnetAZ),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayName),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan", vlan),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.cidr", serviceRangeCidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.subnet_cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.service_endpoints.#", "0"),
				),
			},
			{
				// Step 6 - update resource - update subnet displayName, vlan, cidr and service range cidr
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az           = "%s"
									cidr         = "%s"
									display_name = "%s"
								},
							]
							services = {
								service_range = [
									{
										cidr              = "%s"
										subnet_cidr       = "%s"
										service_endpoints = []
									},
								]
							}
						}
					`, subnetAZ, cidrUpdated, subnetDisplayNameUpdated, serviceRangeCidrUpdated, cidrUpdated),
				) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.az", subnetAZ),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayNameUpdated),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.cidr", serviceRangeCidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.subnet_cidr", cidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.service_endpoints.#", "0"),
				),
			},
			{
				// Step 7 - update resource - add service endpoint
				Config: fmt.Sprintf(testAccStorageEfsConfig, storageEfsServiceName) +
					fmt.Sprintf(testAccVrackServicesConfig, region,
						fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az           = "%s"
									cidr         = "%s"
									display_name = "%s"
								},
							]
							services = {
								service_range = [
									{
										cidr              = "%s"
										subnet_cidr       = "%s"
										service_endpoints = [
											{
												managed_service_urn = data.ovh_storage_efs.tf-acc-efs.iam.urn
											}
										]
									},
								]
							}
						}
					`, subnetAZ, cidrUpdated, subnetDisplayNameUpdated, serviceRangeCidrUpdated, cidrUpdated),
					) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.az", subnetAZ),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayNameUpdated),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.cidr", serviceRangeCidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.subnet_cidr", cidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.service_endpoints.#", "1"),
					resource.TestCheckResourceAttrPair(
						"ovh_vrackservices.tf-acc-vrackservices",
						"target_spec.services.service_range.0.service_endpoints.0.managed_service_urn",
						"data.ovh_storage_efs.tf-acc-efs", "iam.urn",
					),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.services.service_range.0.service_endpoints.#", "1"),
				),
			},
			{
				// Step 8 - update resource - delete service endpoint
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az           = "%s"
									cidr         = "%s"
									display_name = "%s"
								},
							]
							services = {
								service_range = [
									{
										cidr              = "%s"
										subnet_cidr       = "%s"
										service_endpoints = []
									},
								]
							}
						}
					`, subnetAZ, cidrUpdated, subnetDisplayNameUpdated, serviceRangeCidrUpdated, cidrUpdated),
				) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.az", subnetAZ),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayNameUpdated),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.cidr", serviceRangeCidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.subnet_cidr", cidrUpdated),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.service_endpoints.#", "0"),
				),
			},
			{
				// Step 9 - update resource - delete subnet and service range
				Config: fmt.Sprintf(testAccVrackServicesConfig, region, `
						target_spec = {
							subnets = []
						}
					`) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "0"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.subnets.#", "0"),
				),
			},
			{
				// Step 10 - update resource - add subnet with a service endpoint
				Config: fmt.Sprintf(testAccStorageEfsConfig, storageEfsServiceName) +
					fmt.Sprintf(testAccVrackServicesConfig, region,
						fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az           = "%s"
									cidr         = "%s"
									display_name = "%s"
								},
							]
							services = {
								service_range = [
									{
										cidr              = "%s"
										subnet_cidr       = "%s"
										service_endpoints = [
											{
												managed_service_urn = data.ovh_storage_efs.tf-acc-efs.iam.urn
											}
										]
									},
								]
							}
						}
					`, subnetAZ, cidr, subnetDisplayName, serviceRangeCidr, cidr),
					) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.az", subnetAZ),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayName),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.cidr", serviceRangeCidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.subnet_cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.service_endpoints.#", "1"),
					resource.TestCheckResourceAttrPair(
						"ovh_vrackservices.tf-acc-vrackservices",
						"target_spec.services.service_range.0.service_endpoints.0.managed_service_urn",
						"data.ovh_storage_efs.tf-acc-efs", "iam.urn",
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
									az           = "%s"
									cidr         = "%s"
									display_name = "%s"
								},
							]
							services = {
								service_range = [
									{
										cidr              = "%s"
										subnet_cidr       = "%s"
										service_endpoints = [
											{
												managed_service_urn = data.ovh_storage_efs.tf-acc-efs.iam.urn
											}
										]
									},
								]
							}
						}
					`, subnetAZ, cidr, subnetDisplayName, serviceRangeCidr, cidr),
					),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					// VrackServices resource state not updated yet, need a refresh
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.az", subnetAZ),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.display_name", subnetDisplayName),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.cidr", serviceRangeCidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.subnet_cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.services.service_range.0.service_endpoints.#", "1"),
					resource.TestCheckResourceAttrPair(
						"ovh_vrackservices.tf-acc-vrackservices",
						"target_spec.services.service_range.0.service_endpoints.0.managed_service_urn",
						"data.ovh_storage_efs.tf-acc-efs", "iam.urn",
					),
				),
			},
			{
				// Step 12 - refresh to force state update on VrackService resource (current_state.configuration_status)
				PlanOnly:     true,
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.reason", "MISSING_VRACK_ASSOCIATION"),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "current_state.vrack_id"),
				),
			},
			{
				// Step 13 - update resource - delete subnet
				Config: fmt.Sprintf(testAccVrackServicesConfig, region, `
						target_spec = {
							subnets = []
						}
					`) +
					fmt.Sprintf(testAccVrackVrackServicesBinding, vrackServiceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.region", region),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.configuration_status.status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "0"),
				),
			},
		},
	})
}

// TestAccResourceVrackServices_bgpOption orders the BGP option on an existing vRack Services,
// configures it, then terminates it.
func TestAccResourceVrackServices_bgpOption(t *testing.T) {
	region := os.Getenv("OVH_VRACK_SERVICES_REGION")
	subnetAZ := os.Getenv("OVH_VRACK_SERVICES_SUBNET_AZ")
	infraPeers := testAccVrackServicesBGPInfraPeers(region)
	allowedPrefix := os.Getenv("OVH_VRACK_SERVICES_BGP_ALLOWED_PREFIX_TEST")
	skipAllowedPrefix := func() (bool, error) { return allowedPrefix == "", nil }

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckVrackServicesBGP(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1 - order without BGP option
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az   = "%s"
									cidr = "%s"
								},
							]
						}
					`, subnetAZ, cidr),
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.id"),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.client_asn"),
				),
			},
			{
				// Step 2 - order BGP option on the existing vRack Services
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					testAccVrackServicesBGPPlanOption+fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az   = "%s"
									cidr = "%s"
								},
							]
						}
					`, subnetAZ, cidr),
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestMatchResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.id", regexp.MustCompile(`^bgp-[a-z0-9-]+$`)),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.configuration_status.status", "DRAFT"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.client_asn", "64512"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.infrastructure_asn", "65534"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.#", "0"),
				),
			},
			{
				// Step 3 - configure BGP
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					testAccVrackServicesBGPPlanOption+fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az   = "%s"
									cidr = "%s"
								},
							]
							bgp = {
								client_asn         = %s
								infrastructure_asn = %s
								bgp_range = [
									{
										cidr                 = "%s"
										subnet_cidr          = "%s"
										client_peers         = ["%s"]
										infrastructure_peers = %s
									},
								]
							}
						}
					`, subnetAZ, cidr, clientASN, infrastructureASN, bgpRangeCidr, cidr, bgpClientPeer, testAccHCLStringList(infraPeers)),
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestMatchResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.id", regexp.MustCompile(`^bgp-[a-z0-9-]+$`)),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.client_asn", clientASN),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.infrastructure_asn", infrastructureASN),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.client_asn", clientASN),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.infrastructure_asn", infrastructureASN),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.0.cidr", bgpRangeCidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.0.subnet_cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.0.client_peers.0", bgpClientPeer),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.0.infrastructure_peers.#", strconv.Itoa(len(infraPeers))),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.bgp_range.#", "1"),
				),
			},
			{
				// Step 4 - add an allowed prefix: a block of the account parked and eligible for BGP,
				// routed to the BGP by the API
				SkipFunc: skipAllowedPrefix,
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					testAccVrackServicesBGPPlanOption+fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az   = "%s"
									cidr = "%s"
								},
							]
							bgp = {
								client_asn         = %s
								infrastructure_asn = %s
								allowed_prefixes = {
									ipv4 = ["%s"]
									ipv6 = []
								}
								bgp_range = [
									{
										cidr                 = "%s"
										subnet_cidr          = "%s"
										client_peers         = ["%s"]
										infrastructure_peers = %s
									},
								]
							}
						}
					`, subnetAZ, cidr, clientASN, infrastructureASN, allowedPrefix, bgpRangeCidr, cidr, bgpClientPeer, testAccHCLStringList(infraPeers)),
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.allowed_prefixes.ipv4.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.allowed_prefixes.ipv4.0", allowedPrefix),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.allowed_prefixes.ipv4.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.allowed_prefixes.ipv4.0", allowedPrefix),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.configuration_status.status", "READY"),
				),
			},
			{
				// Step 5 - remove the allowed prefix, which unroutes the block; bgp_range stays in the config,
				// as omitting a non-empty one plans its removal
				SkipFunc: skipAllowedPrefix,
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					testAccVrackServicesBGPPlanOption+fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az   = "%s"
									cidr = "%s"
								},
							]
							bgp = {
								client_asn         = %s
								infrastructure_asn = %s
								allowed_prefixes = {
									ipv4 = []
									ipv6 = []
								}
								bgp_range = [
									{
										cidr                 = "%s"
										subnet_cidr          = "%s"
										client_peers         = ["%s"]
										infrastructure_peers = %s
									},
								]
							}
						}
					`, subnetAZ, cidr, clientASN, infrastructureASN, bgpRangeCidr, cidr, bgpClientPeer, testAccHCLStringList(infraPeers)),
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.allowed_prefixes.ipv4.#", "0"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.allowed_prefixes.ipv4.#", "0"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.configuration_status.reason", "MISSING_ALLOWED_PREFIX"),
				),
			},
			{
				// Step 6 - import resource to state
				ResourceName:            "ovh_vrackservices.tf-acc-vrackservices",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"plan", "plan_option", "ovh_subsidiary", "order"},
			},
			{
				// Step 7 - reset BGP ranges
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					testAccVrackServicesBGPPlanOption+fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az   = "%s"
									cidr = "%s"
								},
							]
							bgp = {
								bgp_range = []
							}
						}
					`, subnetAZ, cidr),
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.#", "0"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.bgp_range.#", "0"),
				),
			},
			{
				// Step 8 - terminate BGP option: the termination is only requested, Agora suspends then
				// deletes the BGP later, so it is still there once applied
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az   = "%s"
									cidr = "%s"
								},
							]
						}
					`, subnetAZ, cidr),
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "plan_option.#", "0"),
					resource.TestMatchResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.id", regexp.MustCompile(`^bgp-[a-z0-9-]+$`)),
				),
			},
		},
	})
}

// TestAccResourceVrackServices_bgpDoubleOrder orders a vRack Services along with its BGP option.
func TestAccResourceVrackServices_bgpDoubleOrder(t *testing.T) {
	region := os.Getenv("OVH_VRACK_SERVICES_REGION")
	subnetAZ := os.Getenv("OVH_VRACK_SERVICES_SUBNET_AZ")
	infraPeers := testAccVrackServicesBGPInfraPeers(region)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckVrackServicesBGP(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1 - order vRack Services and BGP option at once, with a BGP configuration
				Config: fmt.Sprintf(testAccVrackServicesConfig, region,
					testAccVrackServicesBGPPlanOption+fmt.Sprintf(`
						target_spec = {
							subnets = [
								{
									az   = "%s"
									cidr = "%s"
								},
							]
							bgp = {
								client_asn         = %s
								infrastructure_asn = %s
								bgp_range = [
									{
										cidr                 = "%s"
										subnet_cidr          = "%s"
										client_peers         = ["%s"]
										infrastructure_peers = %s
									},
								]
							}
						}
					`, subnetAZ, cidr, clientASN, infrastructureASN, bgpRangeCidr, cidr, bgpClientPeer, testAccHCLStringList(infraPeers)),
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "resource_status", "READY"),
					resource.TestMatchResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.id", regexp.MustCompile(`^bgp-[a-z0-9-]+$`)),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.#", "1"),
					testCheckAttrNull("ovh_vrackservices.tf-acc-vrackservices", "target_spec.subnets.0.vlan"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.client_asn", clientASN),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.infrastructure_asn", infrastructureASN),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.client_asn", clientASN),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.infrastructure_asn", infrastructureASN),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.#", "1"),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.0.cidr", bgpRangeCidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.0.subnet_cidr", cidr),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.0.client_peers.0", bgpClientPeer),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "target_spec.bgp.bgp_range.0.infrastructure_peers.#", strconv.Itoa(len(infraPeers))),
					resource.TestCheckResourceAttr("ovh_vrackservices.tf-acc-vrackservices", "current_state.bgp.bgp_range.#", "1"),
				),
			},
			{
				// Step 2 - import resource to state
				ResourceName:            "ovh_vrackservices.tf-acc-vrackservices",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"plan", "plan_option", "ovh_subsidiary", "order"},
			},
		},
	})
}
