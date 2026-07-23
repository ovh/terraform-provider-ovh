package ovh

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// instanceCreationCase describes one point in the instance creation-parameter
// matrix: a distinct combination of the settable creation attributes.
type instanceCreationCase struct {
	name       string
	powerState string // "" => omit (server defaults ACTIVE)
	withSSHKey bool
	withAZ     bool // requires OVH_INSTANCE_AZ_TEST; case is skipped when unset
	sgCount    int  // number of ovh_cloud_security_group to create + attach
	withVolume bool // create + attach one CLASSIC block volume
	networks   string
}

// buildInstanceCreationConfig renders the HCL for one creation case, inlining
// any companion resources (ssh key, security groups, block volume) the case needs.
func buildInstanceCreationConfig(tc instanceCreationCase, serviceName, region, flavorID, imageID, az, suffix string) string {
	var b strings.Builder

	sshKeyLine := ""
	if tc.withSSHKey {
		fmt.Fprintf(&b, `
resource "ovh_cloud_ssh_key" "key" {
  service_name = %q
  name         = "test-key-%s"
  public_key   = %q
}
`, serviceName, suffix, testAccCloudSshKeyPublicKeyA)
		sshKeyLine = "ssh_key_name = ovh_cloud_ssh_key.key.name"
	}

	sgLine := ""
	if tc.sgCount > 0 {
		refs := make([]string, tc.sgCount)
		for i := 0; i < tc.sgCount; i++ {
			fmt.Fprintf(&b, `
resource "ovh_cloud_security_group" "sg%d" {
  service_name = %q
  region       = %q
  name         = "test-sg-%s-%d"
}
`, i, serviceName, region, suffix, i)
			refs[i] = fmt.Sprintf("ovh_cloud_security_group.sg%d.id", i)
		}
		sgLine = fmt.Sprintf("security_group_ids = [%s]", strings.Join(refs, ", "))
	}

	volLine := ""
	if tc.withVolume {
		fmt.Fprintf(&b, `
resource "ovh_cloud_storage_block_volume" "vol" {
  service_name = %q
  region       = %q
  name         = "test-vol-%s"
  size         = 10
  volume_type  = "CLASSIC"
}
`, serviceName, region, suffix)
		volLine = "volume_ids = [ovh_cloud_storage_block_volume.vol.id]"
	}

	azLine := ""
	if tc.withAZ && az != "" {
		azLine = fmt.Sprintf("availability_zone = %q", az)
	}

	powerLine := ""
	if tc.powerState != "" {
		powerLine = fmt.Sprintf("power_state = %q", tc.powerState)
	}

	netLine := ""
	if tc.networks != "" {
		netLine = "networks = " + tc.networks
	}

	fmt.Fprintf(&b, `
resource "ovh_cloud_instance" "test" {
  service_name = %q
  region       = %q
  name         = "test-inst-%s"
  flavor_id    = %q
  image_id     = %q
  %s
  %s
  %s
  %s
  %s
  %s
}
`, serviceName, region, suffix, flavorID, imageID, powerLine, azLine, sshKeyLine, sgLine, volLine, netLine)

	return b.String()
}

// TestAccCloudInstance_creationMatrix exercises the CREATE path across a
// deterministic matrix of parameter combinations. Each case creates an instance
// with its combination, asserts it converges to READY with the expected observed
// state, proves a second apply is a no-op (idempotency / no spurious diff), and
// round-trips through import.
func TestAccCloudInstance_creationMatrix(t *testing.T) {
	const rn = "ovh_cloud_instance.test"

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorID := resolveInstanceFlavorID(t, serviceName, region, testAccInstanceFlavorName)
	imageID := resolveInstanceImageID(t, serviceName, region, testAccInstanceImageName)
	az := os.Getenv("OVH_INSTANCE_AZ_TEST")

	publicNet := `[{ public = true }]`

	cases := []instanceCreationCase{
		{name: "minimal_public", networks: publicNet},
		{name: "no_network"},
		{name: "power_active", powerState: "ACTIVE", networks: publicNet},
		{name: "power_shutoff", powerState: "SHUTOFF", networks: publicNet},
		{name: "power_shelved", powerState: "SHELVED", networks: publicNet},
		{name: "with_ssh_key", withSSHKey: true, networks: publicNet},
		{name: "with_availability_zone", withAZ: true, networks: publicNet},
		{name: "one_security_group", sgCount: 1, networks: publicNet},
		{name: "two_security_groups", sgCount: 2, networks: publicNet},
		{name: "one_volume", withVolume: true, networks: publicNet},
		{name: "ssh_and_two_sgs", withSSHKey: true, sgCount: 2, networks: publicNet},
		{name: "kitchen_sink_shutoff", powerState: "SHUTOFF", withSSHKey: true, sgCount: 1, withVolume: true, networks: publicNet},
		{name: "shutoff_no_network", powerState: "SHUTOFF"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			suffix := acctest.RandomWithPrefix("m")
			config := buildInstanceCreationConfig(tc, serviceName, region, flavorID, imageID, az, suffix)

			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(rn, "resource_status", "READY"),
				resource.TestCheckResourceAttrSet(rn, "id"),
				resource.TestCheckResourceAttrSet(rn, "checksum"),
				resource.TestCheckResourceAttrSet(rn, "created_at"),
			}
			if tc.powerState != "" {
				checks = append(checks,
					resource.TestCheckResourceAttr(rn, "power_state", tc.powerState),
					resource.TestCheckResourceAttr(rn, "current_state.power_state", tc.powerState),
				)
			}
			if tc.withSSHKey {
				checks = append(checks,
					resource.TestCheckResourceAttrPair(rn, "ssh_key_name", "ovh_cloud_ssh_key.key", "name"),
					resource.TestCheckResourceAttr(rn, "current_state.ssh_key_name", "test-key-"+suffix),
				)
			}
			if tc.sgCount > 0 {
				checks = append(checks,
					resource.TestCheckResourceAttr(rn, "security_group_ids.#", strconv.Itoa(tc.sgCount)),
					resource.TestCheckResourceAttr(rn, "current_state.security_groups.#", strconv.Itoa(tc.sgCount)),
				)
			}
			if tc.withVolume {
				checks = append(checks,
					resource.TestCheckResourceAttr(rn, "volume_ids.#", "1"),
					resource.TestCheckResourceAttrSet(rn, "current_state.volumes.0.id"),
					resource.TestCheckResourceAttrSet(rn, "current_state.volumes.0.size"),
				)
			}
			if strings.Contains(tc.networks, "public") {
				checks = append(checks,
					resource.TestCheckResourceAttr(rn, "networks.0.public", "true"),
					resource.TestCheckResourceAttrSet(rn, "current_state.networks.0.addresses.0.ip"),
				)
			}
			if tc.withAZ {
				checks = append(checks,
					resource.TestCheckResourceAttr(rn, "availability_zone", az),
					resource.TestCheckResourceAttr(rn, "current_state.location.availability_zone", az),
				)
			}

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheckCloudInstance(t)
					if tc.withAZ && az == "" {
						t.Skip("OVH_INSTANCE_AZ_TEST must be set for the availability_zone creation case")
					}
				},
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: config,
						Check:  resource.ComposeAggregateTestCheckFunc(checks...),
					},
					{
						// Idempotency: a second apply of the same config must be a no-op.
						Config: config,
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
						},
					},
					{
						ResourceName:            rn,
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateIdFunc:       testAccCloudInstanceImportStateIdFunc(rn),
						ImportStateVerifyIgnore: []string{"checksum"},
					},
				},
			})
		})
	}
}
