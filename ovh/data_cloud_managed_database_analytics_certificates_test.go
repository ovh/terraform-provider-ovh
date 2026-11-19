package ovh

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCloudManagedDatabaseCertificatesRead(t *testing.T) {
	state := managedDatabaseDataSourceTestRead(t, NewCloudManagedDatabaseCertificatesDataSource(),
		map[string]string{"service_name": "pj", "engine": "postgresql", "service_id": managedDatabasePostgresqlTestId},
		map[string][]managedDatabaseTestAnswer{
			"GET " + managedDatabasePostgresqlTestItemPath + "/certificates": {{http.StatusOK, `{"ca": "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----\n"}`}},
		})
	var ca types.String
	state.GetAttribute(context.Background(), path.Root("ca"), &ca)
	if !regexp.MustCompile(`^-----BEGIN CERTIFICATE-----`).MatchString(ca.ValueString()) {
		t.Errorf("unexpected ca %s", ca)
	}
}

func TestAccDataSourceCloudManagedDatabasePostgresqlCertificates_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedDatabasePostgresql(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedDatabasePostgresqlConfig(name, fmt.Sprintf("  description = %q", name), "") + `
data "ovh_cloud_managed_database_certificates" "certificates" {
  service_name = ovh_cloud_managed_database.db.service_name
  engine       = ovh_cloud_managed_database.db.engine
  service_id   = ovh_cloud_managed_database.db.id
}
`,
				Check: resource.TestMatchResourceAttr("data.ovh_cloud_managed_database_certificates.certificates", "ca", regexp.MustCompile(`^-----BEGIN CERTIFICATE-----`)),
			},
		},
	})
}
