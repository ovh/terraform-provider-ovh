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
	testEnterpriseCloudDBConfig = `

data "ovh_enterprise_cloud_db" "db" {
	cluster_id = "%s"
}
	
resource "ovh_enterprise_cloud_db_security_group" "sg" {
  cluster_id = data.ovh_enterprise_cloud_db.db.id
  name = "%s"
}
	
resource "ovh_enterprise_cloud_db_security_group_rule" "rule" {
  cluster_id = data.ovh_enterprise_cloud_db.db.id
  security_group_id = ovh_enterprise_cloud_db_security_group.sg.id
  source = "%s"
}
`
	testEnterpriseCloudDBRuleName1   = "test-foo-bar"
	testEnterpriseCloudDBRuleName2   = "test-bar-foo"
	testEnterpriseCloudDBRuleSource1 = "51.51.51.0/24"
	testEnterpriseCloudDBRuleSource2 = "52.51.51.0/24"
)

func init() {
	resource.AddTestSweepers("ovh_enterprise_cloud_db_security_group", &resource.Sweeper{
		Name: "ovh_enterprise_cloud_db_security_group",
		F:    testSweepEnterpriseCloudDBSecurityGroup,
	})
}

func testSweepEnterpriseCloudDBSecurityGroup(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	clusterId := os.Getenv("OVH_ENTERPRISE_CLOUD_DB")
	if clusterId == "" {
		log.Print("[DEBUG] No OVH_ENTERPRISE_CLOUD_DB envvar specified. nothing to sweep")
		return nil
	}

	var Ids []string
	endpoint := fmt.Sprintf(EnterpriseCloudDBSecurityGroupBaseUrl, url.PathEscape(clusterId))

	if err := client.Get(endpoint, &Ids); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t%q", endpoint, err)
	}

	if len(Ids) == 0 {
		log.Printf("[DEBUG] No SG to sweep on enterprise cloud db %s", clusterId)
		return nil
	}

	for _, id := range Ids {
		urlGet := fmt.Sprintf(EnterpriseCloudDBSecurityGroupBaseUrl+"/%s", clusterId, id)
		var sgResp EnterpriseCloudDBSecurityGroup
		err := client.Get(urlGet, &sgResp)
		if err != nil {
			return err
		}

		if sgResp.Name == testEnterpriseCloudDBRuleName1 || sgResp.Name == testEnterpriseCloudDBRuleName2 {
			continue
		}

		log.Printf(
			"[INFO] Deleting Security Group %s %v on cluster %v",
			id,
			sgResp.Name,
			clusterId,
		)
		endpoint = fmt.Sprintf(EnterpriseCloudDBSecurityGroupBaseUrl+"/%s", url.PathEscape(clusterId), id)
		if err := client.Delete(endpoint, nil); err != nil {
			return fmt.Errorf("Error calling DELETE %s:\n\t%q", endpoint, err)
		}
	}
	return nil
}

func testEnterpriseCloudDBSecurityGroupExist(n string, securityGroup *EnterpriseCloudDBSecurityGroup) resource.TestCheckFunc {
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
		urlGet := fmt.Sprintf(EnterpriseCloudDBSecurityGroupBaseUrl+"/%s", clusterId, currentId)
		err = client.Get(urlGet, securityGroup)

		if err != nil {
			return err
		}
		return nil
	}
}

func testEntrepriseCloudDBSecurityGroupRuleExist(n string, securityGroupRule *EnterpriseCloudDBSecurityGroupRule, securityGroup *EnterpriseCloudDBSecurityGroup) resource.TestCheckFunc {
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
		urlGet := fmt.Sprintf(EnterpriseCloudDBSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, currentId)
		err = client.Get(urlGet, securityGroupRule)

		if err != nil {
			return err
		}
		return nil
	}
}

func testEnterpriseCloudDBNoChangeId(n string, sg *EnterpriseCloudDBSecurityGroup) resource.TestCheckFunc {
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

func TestAccEnterpriseCloudDBSecurityGroup(t *testing.T) {
	cluster := os.Getenv("OVH_ENTERPRISE_CLOUD_DB")
	var securityGroup EnterpriseCloudDBSecurityGroup
	var securityGroupRule EnterpriseCloudDBSecurityGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnterpriseCloudDB(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testEnterpriseCloudDBConfig,
					cluster,
					testEnterpriseCloudDBRuleName1,
					testEnterpriseCloudDBRuleSource1),
				Check: resource.ComposeTestCheckFunc(
					testEnterpriseCloudDBSecurityGroupExist(
						"ovh_enterprise_cloud_db_security_group.sg", &securityGroup),
					testEntrepriseCloudDBSecurityGroupRuleExist(
						"ovh_enterprise_cloud_db_security_group_rule.rule", &securityGroupRule, &securityGroup),
					resource.TestCheckResourceAttr(
						"ovh_enterprise_cloud_db_security_group.sg", "name", testEnterpriseCloudDBRuleName1),
					resource.TestCheckResourceAttr(
						"ovh_enterprise_cloud_db_security_group_rule.rule", "source", testEnterpriseCloudDBRuleSource1),
					resource.TestCheckResourceAttr(
						"data.ovh_enterprise_cloud_db.db", "cluster_id", cluster),
				),
			},
			{
				Config: fmt.Sprintf(
					testEnterpriseCloudDBConfig,
					cluster,
					testEnterpriseCloudDBRuleName2,
					testEnterpriseCloudDBRuleSource2),
				Check: resource.ComposeTestCheckFunc(
					testEnterpriseCloudDBNoChangeId(
						"ovh_enterprise_cloud_db_security_group.sg", &securityGroup),
					resource.TestCheckResourceAttr(
						"ovh_enterprise_cloud_db_security_group.sg", "name", testEnterpriseCloudDBRuleName2),
					resource.TestCheckResourceAttr(
						"ovh_enterprise_cloud_db_security_group_rule.rule", "source", testEnterpriseCloudDBRuleSource2),
				),
			},
		},
	})
}
