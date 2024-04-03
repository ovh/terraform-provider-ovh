package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudProjects_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "ovh_cloud_projects" "projects" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_projects.projects", "projects.0.access", "full"),
					resource.TestCheckResourceAttr("data.ovh_cloud_projects.projects", "projects.0.plan_code", "project.2018"),
					resource.TestCheckResourceAttr("data.ovh_cloud_projects.projects", "projects.0.status", "ok"),
				),
			},
		},
	})
}
