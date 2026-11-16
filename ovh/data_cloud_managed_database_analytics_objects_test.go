package ovh

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// managedDatabaseDataSourceTestRead reads d, configured with the given string arguments, from the fake
// API, and returns its state. It fails when an attribute and the decoded model don't match.
func managedDatabaseDataSourceTestRead(t *testing.T, d datasource.DataSource, arguments map[string]string, answers map[string][]managedDatabaseTestAnswer) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	providerConfig, _ := newManagedDatabaseTestConfig(t, answers)
	d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: providerConfig}, &datasource.ConfigureResponse{})

	var schemaResp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	if diags := schemaResp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("invalid schema: %v", diags)
	}
	// The tests read PostgreSQL services: the engine is postgresql wherever it is an argument.
	if engine, ok := schemaResp.Schema.Attributes["engine"]; ok && engine.IsRequired() && arguments["engine"] == "" {
		arguments["engine"] = "postgresql"
	}
	config := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)}
	for name, value := range arguments {
		if diags := config.SetAttribute(ctx, path.Root(name), types.StringValue(value)); diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
	}

	resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: config.Raw.Copy()}}
	d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: config.Raw}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("expected a fully known state, got %v", resp.State.Raw)
	}
	return resp.State
}

func TestCloudManagedDatabasePostgresqlObjectDataSources(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name, path, key, value, object, list string
		singular, plural                     func() datasource.DataSource
	}{
		{"user", "/user", "id", "u1", managedDatabaseUserTestReady, "users",
			NewCloudManagedDatabasePostgresqlUserDataSource, NewCloudManagedDatabasePostgresqlUsersDataSource},
		{"database", "/database", "id", "d1", managedDatabaseDatabaseTestAnswer, "databases",
			NewCloudManagedDatabaseDatabaseDataSource, NewCloudManagedDatabaseDatabasesDataSource},
		{"connection pool", "/connectionPool", "id", "p1", managedDatabaseConnectionPoolTestAnswer(10, "null"), "connection_pools",
			NewCloudManagedDatabasePostgresqlConnectionPoolDataSource, NewCloudManagedDatabasePostgresqlConnectionPoolsDataSource},
		{"role", "/roles", "id", "r1", `{"id": "r1", "name": "replication", "checksum": "c1", "createdAt": "2026-11-12T10:00:00Z", "updatedAt": "2026-11-12T10:00:00Z"}`, "roles",
			NewCloudManagedDatabaseRoleDataSource, NewCloudManagedDatabaseRolesDataSource},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := map[string]string{"service_name": "pj", "service_id": managedDatabasePostgresqlTestId}
			endpoint := managedDatabasePostgresqlTestItemPath + tc.path

			arguments := map[string]string{tc.key: tc.value}
			for name, value := range service {
				arguments[name] = value
			}
			state := managedDatabaseDataSourceTestRead(t, tc.singular(), arguments, map[string][]managedDatabaseTestAnswer{
				"GET " + endpoint + "/" + tc.value: {{http.StatusOK, tc.object}},
			})
			var key types.String
			state.GetAttribute(ctx, path.Root(tc.key), &key)
			if key.ValueString() != tc.value {
				t.Errorf("unexpected %s %s", tc.key, key)
			}

			if tc.plural == nil {
				return
			}
			state = managedDatabaseDataSourceTestRead(t, tc.plural(), service, map[string][]managedDatabaseTestAnswer{
				"GET " + endpoint: {{http.StatusOK, "[" + tc.object + "]"}},
			})
			var list types.List
			state.GetAttribute(ctx, path.Root(tc.list), &list)
			if len(list.Elements()) != 1 {
				t.Errorf("expected 1 element in %s, got %s", tc.list, list)
			}
		})
	}
}

func testAccDataSourceCloudManagedDatabasePostgresqlObjectsConfig(name string) string {
	return testAccCloudManagedDatabasePostgresqlConnectionPoolConfig(name, `
  mode    = "TRANSACTION"
  size    = 10
  user_id = ovh_cloud_managed_database_postgresql_user.user.id`) + `
locals {
  service = {
    service_name = ovh_cloud_managed_database.db.service_name
    service_id   = ovh_cloud_managed_database.db.id
    engine       = ovh_cloud_managed_database.db.engine
  }
}

data "ovh_cloud_managed_database_postgresql_users" "users" {
  service_name = local.service.service_name
  service_id   = local.service.service_id
  depends_on   = [ovh_cloud_managed_database_postgresql_user.user]
}

data "ovh_cloud_managed_database_postgresql_user" "user" {
  service_name = local.service.service_name
  service_id   = local.service.service_id
  id           = ovh_cloud_managed_database_postgresql_user.user.id
}

data "ovh_cloud_managed_database_databases" "databases" {
  service_name = local.service.service_name
  engine       = local.service.engine
  service_id   = local.service.service_id
  depends_on   = [ovh_cloud_managed_database_database.database]
}

data "ovh_cloud_managed_database_database" "database" {
  service_name = local.service.service_name
  engine       = local.service.engine
  service_id   = local.service.service_id
  id           = ovh_cloud_managed_database_database.database.id
}

data "ovh_cloud_managed_database_postgresql_connection_pools" "pools" {
  service_name = local.service.service_name
  service_id   = local.service.service_id
  depends_on   = [ovh_cloud_managed_database_postgresql_connection_pool.pool]
}

data "ovh_cloud_managed_database_postgresql_connection_pool" "pool" {
  service_name = local.service.service_name
  service_id   = local.service.service_id
  id           = ovh_cloud_managed_database_postgresql_connection_pool.pool.id
}

data "ovh_cloud_managed_database_roles" "roles" {
  service_name = local.service.service_name
  engine       = local.service.engine
  service_id   = local.service.service_id
}

data "ovh_cloud_managed_database_role" "role" {
  service_name = local.service.service_name
  engine       = local.service.engine
  service_id   = local.service.service_id
  id           = data.ovh_cloud_managed_database_roles.roles.roles[0].id
}
`
}

func TestAccDataSourceCloudManagedDatabasePostgresqlObjects_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedDatabasePostgresql(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudManagedDatabasePostgresqlObjectsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair("data.ovh_cloud_managed_database_postgresql_users.users", "users.*.id", testAccCloudManagedDatabasePostgresqlUserResource, "id"),
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_database_postgresql_user.user", "username", "tfacc"),
					resource.TestCheckTypeSetElemAttrPair("data.ovh_cloud_managed_database_databases.databases", "databases.*.id", testAccCloudManagedDatabasePostgresqlDatabaseResource, "id"),
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_database_database.database", "name", "tfacc"),
					resource.TestCheckTypeSetElemAttrPair("data.ovh_cloud_managed_database_postgresql_connection_pools.pools", "connection_pools.*.id", testAccCloudManagedDatabasePostgresqlConnectionPoolResource, "id"),
					resource.TestCheckResourceAttr("data.ovh_cloud_managed_database_postgresql_connection_pool.pool", "mode", "TRANSACTION"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_managed_database_roles.roles", "roles.#"),
					resource.TestCheckResourceAttrPair("data.ovh_cloud_managed_database_role.role", "name", "data.ovh_cloud_managed_database_roles.roles", "roles.0.name"),
				),
			},
		},
	})
}

func TestCloudManagedDatabaseUserDataSourceKinds(t *testing.T) {
	ctx := context.Background()
	common := `"id": "u1", "username": "app", "status": "READY", "admin": false, "checksum": "c1", "createdAt": "2026-11-12T10:00:00Z", "credentialsLastUpdatedAt": "2026-11-12T10:00:00Z"`
	for _, tc := range []struct {
		name, path, engine, user string
		new                      func() datasource.DataSource
	}{
		{"valkey", "/managedDatabase/valkey/", "", `{` + common + `, "categories": ["+@all"], "channels": [], "commands": [], "keys": ["*"]}`, NewCloudManagedDatabaseValkeyUsersDataSource},
		{"opensearch", "/managedAnalytics/opensearch/", "", `{` + common + `, "acls": [{"pattern": "logs-*", "permission": "READ"}]}`, NewCloudManagedAnalyticsOpensearchUsersDataSource},
		{"kafka", "/managedAnalytics/kafka/", "kafka", `{` + common + `}`, NewCloudManagedAnalyticsUsersDataSource},
	} {
		t.Run(tc.name, func(t *testing.T) {
			arguments := map[string]string{"service_name": "pj", "service_id": managedDatabasePostgresqlTestId}
			if tc.engine != "" {
				arguments["engine"] = tc.engine
			}
			state := managedDatabaseDataSourceTestRead(t, tc.new(), arguments, map[string][]managedDatabaseTestAnswer{
				"GET /v2/publicCloud/project/pj" + tc.path + managedDatabasePostgresqlTestId + "/user": {{http.StatusOK, "[" + tc.user + "]"}},
			})
			var users types.List
			state.GetAttribute(ctx, path.Root("users"), &users)
			if len(users.Elements()) != 1 {
				t.Errorf("expected 1 user, got %s", users)
			}
		})
	}
}
