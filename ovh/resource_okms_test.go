package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func kmsDatasourceStateChecks(displayName string) []statecheck.StateCheck {
	return []statecheck.StateCheck{
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("id"),
			"data.ovh_okms_resource.datakms",
			tfjsonpath.New("id"),
			compare.ValuesSame()),
		statecheck.ExpectKnownValue(
			"data.ovh_okms_resource.datakms",
			tfjsonpath.New("iam").AtMapKey("display_name"),
			knownvalue.StringExact(displayName)),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("iam").AtMapKey("urn"),
			"data.ovh_okms_resource.datakms",
			tfjsonpath.New("iam").AtMapKey("urn"),
			compare.ValuesSame()),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("public_ca"),
			"data.ovh_okms_resource.datakms",
			tfjsonpath.New("public_ca"),
			compare.ValuesSame()),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("region"),
			"data.ovh_okms_resource.datakms",
			tfjsonpath.New("region"),
			compare.ValuesSame()),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("rest_endpoint"),
			"data.ovh_okms_resource.datakms",
			tfjsonpath.New("rest_endpoint"),
			compare.ValuesSame()),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("kmip_endpoint"),
			"data.ovh_okms_resource.datakms",
			tfjsonpath.New("kmip_endpoint"),
			compare.ValuesSame()),
		statecheck.CompareValuePairs(
			"ovh_okms.kms",
			tfjsonpath.New("swagger_endpoint"),
			"data.ovh_okms_resource.datakms",
			tfjsonpath.New("swagger_endpoint"),
			compare.ValuesSame()),
	}
}

func kmsResourceStateChecks(displayName string) []statecheck.StateCheck {
	urnRe := regexp.MustCompile("urn:v1:eu:resource:okms:.*")
	return []statecheck.StateCheck{
		statecheck.ExpectKnownValue(
			"ovh_okms.kms",
			tfjsonpath.New("display_name"),
			knownvalue.StringExact(displayName)),
		statecheck.ExpectKnownValue(
			"ovh_okms.kms",
			tfjsonpath.New("iam").AtMapKey("display_name"),
			knownvalue.StringExact(displayName)),
		statecheck.ExpectKnownValue(
			"ovh_okms.kms",
			tfjsonpath.New("iam").AtMapKey("urn"),
			knownvalue.StringRegexp(urnRe)),
		statecheck.ExpectKnownValue(
			"ovh_okms.kms",
			tfjsonpath.New("kmip_endpoint"),
			knownvalue.NotNull()),
		statecheck.ExpectKnownValue(
			"ovh_okms.kms",
			tfjsonpath.New("swagger_endpoint"),
			knownvalue.NotNull()),
		statecheck.ExpectKnownValue(
			"ovh_okms.kms",
			tfjsonpath.New("rest_endpoint"),
			knownvalue.NotNull()),
		statecheck.ExpectKnownValue(
			"ovh_okms.kms",
			tfjsonpath.New("public_ca"),
			knownvalue.NotNull()),
	}
}

const confOkmsResourceOrder = `
resource "ovh_okms" "kms" {
  ovh_subsidiary = "FR"
  display_name = "%s"
  region = "%s"
}
`

const confOkmsDatasource = `
resource "ovh_okms" "kms" {
  ovh_subsidiary = "FR"
  display_name = "%s"
  region = "%s"
}

data "ovh_okms_resource" "datakms" {
  id = ovh_okms.kms.id
}
`

func TestAccOkmsOrder(t *testing.T) {
	compareValuesSame := statecheck.CompareValue(compare.ValuesSame())
	displayName := acctest.RandomWithPrefix(test_prefix)
	region := "EU_WEST_SBG"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(confOkmsResourceOrder, displayName, region),
				ConfigStateChecks: append(
					kmsResourceStateChecks(displayName),
					statecheck.ExpectKnownValue(
						"ovh_okms.kms",
						tfjsonpath.New("id"),
						knownvalue.NotNull()),
					compareValuesSame.AddStateValue(
						"ovh_okms.kms",
						tfjsonpath.New("id")),
				),
			},
			{
				// Test update
				Config: fmt.Sprintf(confOkmsResourceOrder, displayName+"new", region),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"ovh_okms.kms",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: append(
					kmsResourceStateChecks(displayName+"new"),
					compareValuesSame.AddStateValue(
						"ovh_okms.kms",
						tfjsonpath.New("id")),
				),
			},
			{
				// Test datasource
				Config:            fmt.Sprintf(confOkmsDatasource, displayName+"new", region),
				ConfigStateChecks: kmsDatasourceStateChecks(displayName + "new"),
			},
		},
	})
}

func TestAccOkmsImport(t *testing.T) {
	displayName := acctest.RandomWithPrefix(test_prefix)
	kmsId := os.Getenv("OVH_OKMS")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckOkms(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:  "ovh_okms.kms",
				ImportState:   true,
				ImportStateId: kmsId,
				Config:        fmt.Sprintf(confOkmsResourceOrder, displayName, "EU_WEST_SBG"),
			},
		},
	})
}
