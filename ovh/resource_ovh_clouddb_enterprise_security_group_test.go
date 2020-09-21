package ovh

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"net/url"
	"os"
	"testing"
)

const (
	testCloudDBEnterpriseConfig = `

data "ovh_clouddb_enterprise" "db" {
	cluster_id = "%s"
}
	
resource "ovh_clouddb_enterprise_security_group" "sg" {
  cluster_id = data.ovh_clouddb_enterprise.db.id
  name = "%s"
}
	
resource "ovh_clouddb_enterprise_security_group_rule" "rule" {
  cluster_id = data.ovh_clouddb_enterprise.db.id
  security_group_id = ovh_clouddb_enterprise_security_group.sg.id
  source = "%s"
}
`
	testCloudDBEnterpriseRuleName1   = "test-foo-bar"
	testCloudDBEnterpriseRuleName2   = "test-bar-foo"
	testCloudDBEnterpriseRuleSource1 = "51.51.51.0/24"
	testCloudDBEnterpriseRuleSource2 = "52.51.51.0/24"
)

func init() {
	resource.AddTestSweepers("ovh_clouddb_enterprise_security_group", &resource.Sweeper{
		Name: "ovh_clouddb_enterprise_security_group",
		F:    testSweepCloudDBEnterpriseSecurityGroup,
	})
}

func testSweepCloudDBEnterpriseSecurityGroup(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	clusterId := os.Getenv("OVH_CLOUDDB_ENTERPRISE")
	if clusterId == "" {
		log.Print("[DEBUG] No OVH_CLOUDDB_ENTERPRISE envvar specified. nothing to sweep")
		return nil
	}

	var Ids []string
	endpoint := fmt.Sprintf(CloudDBEnterpriseSecurityGroupBaseUrl, url.PathEscape(clusterId))

	if err := client.Get(endpoint, &Ids); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t%q", endpoint, err)
	}

	if len(Ids) == 0 {
		log.Printf("[DEBUG] No SG to sweep on enterprise cloud db %s", clusterId)
		return nil
	}

	for _, id := range Ids {
		urlGet := fmt.Sprintf(CloudDBEnterpriseSecurityGroupBaseUrl+"/%s", clusterId, id)
		var sgResp CloudDBEnterpriseSecurityGroup
		err := client.Get(urlGet, &sgResp)
		if err != nil {
			return err
		}

		if sgResp.Name == testCloudDBEnterpriseRuleName1 || sgResp.Name == testCloudDBEnterpriseRuleName2 {
			continue
		}

		log.Printf(
			"[INFO] Deleting Security Group %s %v on cluster %v",
			id,
			sgResp.Name,
			clusterId,
		)
		endpoint = fmt.Sprintf(CloudDBEnterpriseSecurityGroupBaseUrl+"/%s", url.PathEscape(clusterId), id)
		if err := client.Delete(endpoint, nil); err != nil {
			return fmt.Errorf("Error calling DELETE %s:\n\t%q", endpoint, err)
		}
	}
	return nil
}

func testCloudDBEnterpriseSecurityGroupExist(n string, securityGroup *CloudDBEnterpriseSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		client, err := sharedClientForRegion("")
		if err != nil {
			return fmt.Errorf("error getting client: %s", err)
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		currentId := rs.Primary.Attributes["id"]
		urlGet := fmt.Sprintf(CloudDBEnterpriseSecurityGroupBaseUrl+"/%s", clusterId, currentId)
		err = client.Get(urlGet, securityGroup)

		if err != nil {
			return err
		}
		return nil
	}
}

func testEntrepriseCloudDBSecurityGroupRuleExist(n string, securityGroupRule *CloudDBEnterpriseSecurityGroupRule, securityGroup *CloudDBEnterpriseSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		client, err := sharedClientForRegion("")
		if err != nil {
			return fmt.Errorf("error getting client: %s", err)
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		securityGroupId := rs.Primary.Attributes["security_group_id"]

		if securityGroupId != securityGroup.Id {
			return fmt.Errorf("Security Group don't match: %s / %s", securityGroupId, securityGroup.Id)
		}

		currentId := rs.Primary.Attributes["id"]
		urlGet := fmt.Sprintf(CloudDBEnterpriseSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, currentId)
		err = client.Get(urlGet, securityGroupRule)

		if err != nil {
			return err
		}
		return nil
	}
}

func testCloudDBEnterpriseNoChangeId(n string, sg *CloudDBEnterpriseSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		resourceId := rs.Primary.Attributes["id"]

		if resourceId != sg.Id {
			return fmt.Errorf("Id changed from: %s to %s", sg.Id, resourceId)
		}

		return nil
	}
}

func TestAccCloudDBEnterpriseSecurityGroup(t *testing.T) {
	cluster := os.Getenv("OVH_CLOUDDB_ENTERPRISE")
	var securityGroup CloudDBEnterpriseSecurityGroup
	var securityGroupRule CloudDBEnterpriseSecurityGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDBEnterprise(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testCloudDBEnterpriseConfig,
					cluster,
					testCloudDBEnterpriseRuleName1,
					testCloudDBEnterpriseRuleSource1),
				Check: resource.ComposeTestCheckFunc(
					testCloudDBEnterpriseSecurityGroupExist(
						"ovh_clouddb_enterprise_security_group.sg", &securityGroup),
					testEntrepriseCloudDBSecurityGroupRuleExist(
						"ovh_clouddb_enterprise_security_group_rule.rule", &securityGroupRule, &securityGroup),
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_security_group.sg", "name", testCloudDBEnterpriseRuleName1),
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_security_group_rule.rule", "source", testCloudDBEnterpriseRuleSource1),
					resource.TestCheckResourceAttr(
						"data.ovh_clouddb_enterprise.db", "cluster_id", cluster),
				),
			},
			{
				Config: fmt.Sprintf(
					testCloudDBEnterpriseConfig,
					cluster,
					testCloudDBEnterpriseRuleName2,
					testCloudDBEnterpriseRuleSource2),
				Check: resource.ComposeTestCheckFunc(
					testCloudDBEnterpriseNoChangeId(
						"ovh_clouddb_enterprise_security_group.sg", &securityGroup),
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_security_group.sg", "name", testCloudDBEnterpriseRuleName2),
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_security_group_rule.rule", "source", testCloudDBEnterpriseRuleSource2),
				),
			},
		},
	})
}
