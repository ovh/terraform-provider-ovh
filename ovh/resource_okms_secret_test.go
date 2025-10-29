package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Step 1: create secret v1
const testAccOkmsSecretResourceConfigV1 = `
resource "ovh_okms_secret" "test" {
  okms_id = "%s"
  path    = "%s"
  version = {
    data = jsonencode({ initial = "v1" })
  }
}
`

// Step 2: update secret with cas=1 to create v2
const testAccOkmsSecretResourceConfigV2 = `
resource "ovh_okms_secret" "test" {
  okms_id = "%s"
  path    = "%s"
  cas     = 1
  version = {
    data = jsonencode({ initial = "v1", second = "v2" })
  }
}
`

// Step 3: update again with cas=2 and change metadata (max_versions) to ensure metadata update and version 3 creation.
const testAccOkmsSecretResourceConfigV3 = `
resource "ovh_okms_secret" "test" {
  okms_id = "%s"
  path    = "%s"
  cas     = 2
  metadata = {
    max_versions = 5
  }
  version = {
    data = jsonencode({ initial = "v1", second = "v2", third = "v3" })
  }
}
`

// Update attempt with an incorrect CAS value (expected to fail)
const testAccOkmsSecretResourceConfigWrongCas = `
resource "ovh_okms_secret" "test" {
	okms_id = "%s"
	path    = "%s"
	cas     = 9999
	version = {
		data = jsonencode({ initial = "v1", wrong = "cas" })
	}
}
`

// Update to create v2 with cas=1 and set metadata.cas_required = true
const testAccOkmsSecretResourceConfigV2CasRequired = `
resource "ovh_okms_secret" "test" {
	okms_id = "%s"
	path    = "%s"
	cas     = 1
	metadata = {
		cas_required = true
	}
	version = {
		data = jsonencode({ initial = "v1", second = "v2" })
	}
}
`

// Metadata-only update (same version data) should NOT create a new version (version.id stays constant)
const testAccOkmsSecretResourceConfigMetadataOnlyNoNewVersion = `
resource "ovh_okms_secret" "test" {
	okms_id = "%s"
	path    = "%s"
	version = {
		data = jsonencode({ initial = "v1" })
	}
}
`

const testAccOkmsSecretResourceConfigMetadataOnlyNoNewVersionStep2 = `
resource "ovh_okms_secret" "test" {
	okms_id = "%s"
	path    = "%s"
	cas     = 1
	metadata = {
		max_versions = 9
	}
	// SAME DATA -> should not create a new version
	version = {
		data = jsonencode({ initial = "v1" })
	}
}
`

// Create secret with initial metadata to test MetadataValue.ToCreate implementation
const testAccOkmsSecretResourceConfigCreateWithMetadata = `
resource "ovh_okms_secret" "test" {
	okms_id = "%s"
	path    = "%s"
	metadata = {
		cas_required = true
		max_versions = 7
		deactivate_version_after = "0s"
		custom_metadata = {
			env = "acc"
		}
	}
	version = {
		data = jsonencode({ initial = "v1" })
	}
}
`

func TestAccOkmsSecretResource_basicLifecycle(t *testing.T) {
	okmsID := os.Getenv("OVH_OKMS")
	if okmsID == "" {
		checkEnvOrSkip(t, "OVH_OKMS")
	}
	path := fmt.Sprintf("tfacc-%s", acctest.RandString(6))

	configV1 := fmt.Sprintf(testAccOkmsSecretResourceConfigV1, okmsID, path)
	configV2 := fmt.Sprintf(testAccOkmsSecretResourceConfigV2, okmsID, path)
	configV3 := fmt.Sprintf(testAccOkmsSecretResourceConfigV3, okmsID, path)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOkms(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configV1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "version.id", "1"),
					resource.TestCheckResourceAttrSet("ovh_okms_secret.test", "version.data"),
				),
			},
			{
				Config: configV2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "version.id", "2"),
					resource.TestCheckResourceAttrSet("ovh_okms_secret.test", "version.data"),
				),
			},
			{
				Config: configV3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "version.id", "3"),
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "metadata.max_versions", "5"),
					resource.TestCheckResourceAttrSet("ovh_okms_secret.test", "version.data"),
				),
			},
		},
	})
}

// Verify CAS mismatch returns an error when wrong cas is provided.
func TestAccOkmsSecretResource_casMismatch(t *testing.T) {
	okmsID := os.Getenv("OVH_OKMS")
	if okmsID == "" {
		checkEnvOrSkip(t, "OVH_OKMS")
	}
	path := fmt.Sprintf("tfacc-%s", acctest.RandString(6))

	// Create v1
	configCreate := fmt.Sprintf(testAccOkmsSecretResourceConfigV1, okmsID, path)
	// Attempt update with wrong cas (expect failure)
	wrongCasConfig := fmt.Sprintf(testAccOkmsSecretResourceConfigWrongCas, okmsID, path)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOkms(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: configCreate},
			{
				Config:      wrongCasConfig,
				ExpectError: regexp.MustCompile(`(?i)cas`),
			},
		},
	})
}

// Verify updating only metadata (cas required) also creates a new version when version data changes
// and that changing path forces recreation.
func TestAccOkmsSecretResource_pathRecreateAndMetadata(t *testing.T) {
	okmsID := os.Getenv("OVH_OKMS")
	if okmsID == "" {
		checkEnvOrSkip(t, "OVH_OKMS")
	}
	path1 := fmt.Sprintf("tfacc-%s", acctest.RandString(6))
	path2 := path1 + "-b"

	// initial create v1
	cfg1 := fmt.Sprintf(testAccOkmsSecretResourceConfigV1, okmsID, path1)
	// update to v2 with cas=1 and set cas_required true via metadata
	cfg2 := fmt.Sprintf(testAccOkmsSecretResourceConfigV2CasRequired, okmsID, path1)
	// change path should force recreation (id reset to 1 again) because RequiresReplace
	cfg3 := fmt.Sprintf(testAccOkmsSecretResourceConfigV1, okmsID, path2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOkms(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "version.id", "1"),
				),
			},
			{
				Config: cfg2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "version.id", "2"),
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "metadata.cas_required", "true"),
				),
			},
			{
				Config: cfg3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "version.id", "1"),
				),
			},
		},
	})
}

// Test metadata-only update does not create a new version
func TestAccOkmsSecretResource_metadataOnlyNoNewVersion(t *testing.T) {
	okmsID := os.Getenv("OVH_OKMS")
	if okmsID == "" {
		checkEnvOrSkip(t, "OVH_OKMS")
	}
	path := fmt.Sprintf("tfacc-%s", acctest.RandString(6))

	cfg1 := fmt.Sprintf(testAccOkmsSecretResourceConfigMetadataOnlyNoNewVersion, okmsID, path)
	cfg2 := fmt.Sprintf(testAccOkmsSecretResourceConfigMetadataOnlyNoNewVersionStep2, okmsID, path)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOkms(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "version.id", "1"),
				),
			},
			{
				Config: cfg2,
				Check: resource.ComposeTestCheckFunc(
					// version should remain 1 because data unchanged
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "version.id", "1"),
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "metadata.max_versions", "9"),
				),
			},
		},
	})
}

// Test creation with metadata supplied at creation time (no subsequent update needed)
func TestAccOkmsSecretResource_createWithMetadata(t *testing.T) {
	okmsID := os.Getenv("OVH_OKMS")
	if okmsID == "" {
		checkEnvOrSkip(t, "OVH_OKMS")
	}
	path := fmt.Sprintf("tfacc-%s", acctest.RandString(6))

	cfg := fmt.Sprintf(testAccOkmsSecretResourceConfigCreateWithMetadata, okmsID, path)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOkms(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "version.id", "1"),
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "metadata.cas_required", "true"),
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "metadata.max_versions", "7"),
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "metadata.deactivate_version_after", "0s"),
					resource.TestCheckResourceAttr("ovh_okms_secret.test", "metadata.custom_metadata.env", "acc"),
				),
			},
		},
	})
}
