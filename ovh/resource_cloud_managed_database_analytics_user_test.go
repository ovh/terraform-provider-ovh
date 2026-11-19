package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fwpath "github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const (
	testAccCloudManagedDatabasePostgresqlUserResource = "ovh_cloud_managed_database_postgresql_user.user"
	managedDatabaseUserTestPath                       = managedDatabasePostgresqlTestItemPath + "/user"
	managedDatabaseUserTestItemPath                   = managedDatabaseUserTestPath + "/u1"
	managedDatabaseUserTestReady                      = `{"id": "u1", "username": "app", "roles": ["replication"], "status": "READY", "admin": false, "checksum": "c1", "createdAt": "2026-11-12T10:00:00Z", "credentialsLastUpdatedAt": "2026-11-12T10:00:00Z"}`
)

func newManagedDatabaseUserTestResource(t *testing.T, answers map[string][]managedDatabaseTestAnswer) (*managedDatabaseUserResource[managedDatabaseRolesUserModel, *managedDatabaseRolesUserModel], func(string) []string) {
	t.Helper()
	config, bodies := newManagedDatabaseTestConfig(t, answers)
	r := NewCloudManagedDatabasePostgresqlUserResource().(*managedDatabaseUserResource[managedDatabaseRolesUserModel, *managedDatabaseRolesUserModel])
	r.config = config
	return r, bodies
}

// managedDatabaseUserTestModel returns the state of the user of managedDatabaseUserTestReady,
// created with the password "secret".
func managedDatabaseUserTestModel() managedDatabaseRolesUserModel {
	return managedDatabaseRolesUserModel{
		ServiceName:              types.StringValue("pj"),
		Engine:                   types.StringValue("postgresql"),
		ServiceId:                types.StringValue(managedDatabasePostgresqlTestId),
		Name:                     types.StringValue("app"),
		Roles:                    types.SetValueMust(types.StringType, []attr.Value{types.StringValue("replication")}),
		PasswordReset:            types.StringNull(),
		Id:                       types.StringValue("u1"),
		Username:                 types.StringValue("app"),
		Password:                 types.StringValue("secret"),
		Admin:                    types.BoolValue(false),
		Status:                   types.StringValue("READY"),
		Checksum:                 types.StringValue("c1"),
		CreatedAt:                types.StringValue("2026-11-12T10:00:00Z"),
		CredentialsLastUpdatedAt: types.StringValue("2026-11-12T10:00:00Z"),
	}
}

func TestCloudManagedDatabasePostgresqlUserCreate(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseUserTestResource(t, map[string][]managedDatabaseTestAnswer{
		"POST " + managedDatabaseUserTestPath:    {{http.StatusOK, `{"id": "u1", "username": "app", "password": "secret", "roles": ["replication"], "status": "CREATING", "checksum": "c0", "createdAt": "2026-11-12T10:00:00Z", "credentialsLastUpdatedAt": "2026-11-12T10:00:00Z"}`}},
		"GET " + managedDatabaseUserTestItemPath: {{http.StatusOK, `{"id": "u1", "status": "CREATING"}`}, {http.StatusOK, managedDatabaseUserTestReady}},
	})

	plan := managedDatabaseUserTestModel()
	plan.Id, plan.Username, plan.Password, plan.Status = types.StringUnknown(), types.StringUnknown(), types.StringUnknown(), types.StringUnknown()
	plan.Admin, plan.Checksum, plan.CreatedAt, plan.CredentialsLastUpdatedAt = types.BoolUnknown(), types.StringUnknown(), types.StringUnknown(), types.StringUnknown()
	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	assertManagedDatabaseJSON(t, json.RawMessage(bodies("POST " + managedDatabaseUserTestPath)[0]), `{"name": "app", "roles": ["replication"]}`)
	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("expected a fully known state, got %v", resp.State.Raw)
	}
	var got managedDatabaseRolesUserModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	// The password is only in the answer to the POST.
	if got.Password.ValueString() != "secret" || got.Status.ValueString() != "READY" {
		t.Errorf("unexpected state: password %s, status %s", got.Password, got.Status)
	}
}

func TestCloudManagedDatabasePostgresqlUserUpdate(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseUserTestResource(t, map[string][]managedDatabaseTestAnswer{
		"PUT " + managedDatabaseUserTestItemPath:                     {{http.StatusConflict, managedDatabasePostgresqlTestBusy}, {http.StatusOK, managedDatabaseUserTestReady}},
		"POST " + managedDatabaseUserTestItemPath + "/resetPassword": {{http.StatusOK, `{"id": "u1", "username": "app", "password": "new-secret", "status": "UPDATING"}`}},
		"GET " + managedDatabaseUserTestItemPath:                     {{http.StatusOK, managedDatabaseUserTestReady}},
	})

	state := managedDatabaseUserTestModel()
	plan := state
	plan.Roles = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("pg_read_all_stats")})
	plan.PasswordReset = types.StringValue("v2")
	plan.Password, plan.CredentialsLastUpdatedAt = types.StringUnknown(), types.StringUnknown()
	plan.Status, plan.Checksum = types.StringUnknown(), types.StringUnknown()
	resp := fwresource.UpdateResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Update(ctx, fwresource.UpdateRequest{
		Plan:  managedDatabaseObjectTestPlan(t, r, &plan),
		State: managedDatabaseObjectTestState(t, r, &state),
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	puts := bodies("PUT " + managedDatabaseUserTestItemPath)
	if len(puts) != 2 {
		t.Fatalf("expected the PUT to be retried once after the 409, got %d PUT", len(puts))
	}
	assertManagedDatabaseJSON(t, json.RawMessage(puts[1]), `{"roles": ["pg_read_all_stats"]}`)
	var got managedDatabaseRolesUserModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &got)...)
	if got.Password.ValueString() != "new-secret" {
		t.Errorf("expected the password of the reset, got %s", got.Password)
	}
}

func TestCloudManagedDatabasePostgresqlUserModifyPlan(t *testing.T) {
	ctx := context.Background()
	r := NewCloudManagedDatabasePostgresqlUserResource().(*managedDatabaseUserResource[managedDatabaseRolesUserModel, *managedDatabaseRolesUserModel])
	state := managedDatabaseUserTestModel()
	state.PasswordReset = types.StringValue("v1")

	for name, change := range map[string]func(*managedDatabaseRolesUserModel){
		"roles":          func(m *managedDatabaseRolesUserModel) { m.Roles = types.SetValueMust(types.StringType, []attr.Value{}) },
		"password reset": func(m *managedDatabaseRolesUserModel) { m.PasswordReset = types.StringValue("v2") },
		// Removing password_reset keeps the password that clients use.
		"password reset removed": func(m *managedDatabaseRolesUserModel) { m.PasswordReset = types.StringNull() },
	} {
		t.Run(name, func(t *testing.T) {
			plan := state
			change(&plan)
			resp := fwresource.ModifyPlanResponse{Plan: managedDatabaseObjectTestPlan(t, r, &plan)}
			r.ModifyPlan(ctx, fwresource.ModifyPlanRequest{
				Plan:  managedDatabaseObjectTestPlan(t, r, &plan),
				State: managedDatabaseObjectTestState(t, r, &state),
			}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			var got managedDatabaseRolesUserModel
			resp.Diagnostics.Append(resp.Plan.Get(ctx, &got)...)
			if want := name == "password reset"; got.Password.IsUnknown() != want {
				t.Errorf("expected an unknown password: %t, got %s", want, got.Password)
			}
		})
	}
}

func TestCloudManagedDatabasePostgresqlUserDelete(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseUserTestResource(t, map[string][]managedDatabaseTestAnswer{
		"DELETE " + managedDatabaseUserTestItemPath: {{http.StatusOK, `{"id": "u1", "status": "DELETING"}`}},
		"GET " + managedDatabaseUserTestItemPath:    {{http.StatusOK, `{"id": "u1", "status": "DELETING"}`}, {http.StatusNotFound, `{"message": "user not found"}`}},
	})

	state := managedDatabaseUserTestModel()
	resp := fwresource.DeleteResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if got := len(bodies("DELETE " + managedDatabaseUserTestItemPath)); got != 1 {
		t.Errorf("expected 1 DELETE, got %d", got)
	}
}

// testAccCloudManagedDatabasePostgresqlUserConfig renders the user tfacc, with password_reset = null,
// as if removed from the configuration, when passwordReset is "".
func testAccCloudManagedDatabasePostgresqlUserConfig(name, passwordReset string) string {
	value := "null"
	if passwordReset != "" {
		value = fmt.Sprintf("%q", passwordReset)
	}
	return testAccCloudManagedDatabasePostgresqlConfig(name, fmt.Sprintf("  description = %q", name), "") + fmt.Sprintf(`
resource "ovh_cloud_managed_database_postgresql_user" "user" {
  service_name   = ovh_cloud_managed_database.db.service_name
  service_id     = ovh_cloud_managed_database.db.id
  name           = "tfacc"
  password_reset = %s
}
`, value)
}

func TestAccCloudManagedDatabasePostgresqlUser_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedDatabasePostgresql(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudManagedDatabasePostgresqlUserConfig(name, "v1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlUserResource, "username", "tfacc"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlUserResource, "status", "READY"),
					resource.TestCheckResourceAttr(testAccCloudManagedDatabasePostgresqlUserResource, "admin", "false"),
					resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlUserResource, "password"),
				),
			},
			{
				Config: testAccCloudManagedDatabasePostgresqlUserConfig(name, "v2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlUserResource, plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue(testAccCloudManagedDatabasePostgresqlUserResource, tfjsonpath.New("password")),
					},
				},
				Check: resource.TestCheckResourceAttrSet(testAccCloudManagedDatabasePostgresqlUserResource, "password"),
			},
			{
				// Removing password_reset keeps the password.
				Config: testAccCloudManagedDatabasePostgresqlUserConfig(name, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccCloudManagedDatabasePostgresqlUserResource, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(testAccCloudManagedDatabasePostgresqlUserResource, tfjsonpath.New("password"), knownvalue.NotNull()),
					},
				},
			},
			{
				ResourceName:            testAccCloudManagedDatabasePostgresqlUserResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "password_reset"},
				ImportStateIdFunc:       testAccCloudManagedImportStateIdFunc(testAccCloudManagedDatabasePostgresqlUserResource),
			},
		},
	})
}

// managedDatabaseUserTestPlanned returns the planned attributes of a user named app, its computed
// attributes unknown.
func managedDatabaseUserTestPlanned(engine types.String) managedDatabaseUserModel {
	return managedDatabaseUserModel{
		ServiceName:              types.StringValue("pj"),
		Engine:                   engine,
		ServiceId:                types.StringValue(managedDatabasePostgresqlTestId),
		Name:                     types.StringValue("app"),
		PasswordReset:            types.StringNull(),
		Id:                       types.StringUnknown(),
		Username:                 types.StringUnknown(),
		Password:                 types.StringUnknown(),
		Admin:                    types.BoolUnknown(),
		Status:                   types.StringUnknown(),
		Checksum:                 types.StringUnknown(),
		CreatedAt:                types.StringUnknown(),
		CredentialsLastUpdatedAt: types.StringUnknown(),
	}
}

// testManagedDatabaseUserKindCreate creates a user of the service at path from plan, and checks the
// body of the POST and that the state is fully known. fields are the fields of the kind in the answers.
func testManagedDatabaseUserKindCreate(t *testing.T, r fwresource.Resource, plan any, service, fields, body string) {
	t.Helper()
	ctx := context.Background()
	common := `"id": "u1", "username": "app", "status": "READY", "admin": false, "checksum": "c1", "createdAt": "2026-11-12T10:00:00Z", "credentialsLastUpdatedAt": "2026-11-12T10:00:00Z"`
	config, bodies := newManagedDatabaseTestConfig(t, map[string][]managedDatabaseTestAnswer{
		"POST " + service + "/user":   {{http.StatusOK, `{"password": "secret", ` + common + fields + `}`}},
		"GET " + service + "/user/u1": {{http.StatusOK, `{` + common + fields + `}`}},
	})
	r.(fwresource.ResourceWithConfigure).Configure(ctx, fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})

	resp := fwresource.CreateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.Create(ctx, fwresource.CreateRequest{Plan: managedDatabaseObjectTestPlan(t, r, plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	assertManagedDatabaseJSON(t, json.RawMessage(bodies("POST " + service + "/user")[0]), body)
	if !resp.State.Raw.IsFullyKnown() {
		t.Errorf("expected a fully known state, got %v", resp.State.Raw)
	}
}

func TestCloudManagedDatabaseUserKinds(t *testing.T) {
	database, analytics, service := "/v2/publicCloud/project/pj/managedDatabase/", "/v2/publicCloud/project/pj/managedAnalytics/", "/"+managedDatabasePostgresqlTestId

	t.Run("mysql", func(t *testing.T) {
		plan := managedDatabasePlainUserModel{managedDatabaseUserTestPlanned(types.StringValue("mysql"))}
		testManagedDatabaseUserKindCreate(t, NewCloudManagedDatabaseUserResource(), &plan, database+"mysql"+service, "", `{"name": "app"}`)
	})
	t.Run("valkey", func(t *testing.T) {
		plan := managedDatabaseValkeyUserModel{
			managedDatabaseUserModel: managedDatabaseUserTestPlanned(types.StringValue("valkey")),
			Categories:               types.ListValueMust(types.StringType, []attr.Value{types.StringValue("+@all")}),
			Channels:                 types.ListUnknown(types.StringType),
			Commands:                 types.ListUnknown(types.StringType),
			Keys:                     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("*")}),
		}
		testManagedDatabaseUserKindCreate(t, NewCloudManagedDatabaseValkeyUserResource(), &plan, database+"valkey"+service,
			`, "categories": ["+@all"], "channels": [], "commands": [], "keys": ["*"]`, `{"name": "app", "categories": ["+@all"], "keys": ["*"]}`)
	})
	t.Run("opensearch", func(t *testing.T) {
		acl := managedDatabaseOpensearchUserAttributes()["acls"].GetType().(types.SetType).ElemType.(types.ObjectType)
		plan := managedDatabaseOpensearchUserModel{
			managedDatabaseUserModel: managedDatabaseUserTestPlanned(types.StringValue("opensearch")),
			Acls: types.SetValueMust(acl, []attr.Value{types.ObjectValueMust(acl.AttrTypes, map[string]attr.Value{
				"pattern": types.StringValue("logs-*"), "permission": types.StringValue("READ"),
			})}),
		}
		testManagedDatabaseUserKindCreate(t, NewCloudManagedAnalyticsOpensearchUserResource(), &plan, analytics+"opensearch"+service,
			`, "acls": [{"pattern": "logs-*", "permission": "READ"}]`, `{"name": "app", "acls": [{"pattern": "logs-*", "permission": "READ"}]}`)
	})
	t.Run("clickhouse", func(t *testing.T) {
		// ClickHouse returns the roles of a user, but doesn't take them.
		plan := managedDatabaseRolesUserModel{managedDatabaseUserModel: managedDatabaseUserTestPlanned(types.StringValue("clickhouse")), Roles: types.SetUnknown(types.StringType)}
		testManagedDatabaseUserKindCreate(t, NewCloudManagedAnalyticsClickhouseUserResource(), &plan, analytics+"clickhouse"+service,
			`, "roles": ["default"]`, `{"name": "app"}`)
	})
}

func TestCloudManagedDatabaseValkeyUserUpdateBody(t *testing.T) {
	ctx := context.Background()
	list := func(values ...string) types.List {
		elements := make([]attr.Value, 0, len(values))
		for _, value := range values {
			elements = append(elements, types.StringValue(value))
		}
		return types.ListValueMust(types.StringType, elements)
	}
	state := managedDatabaseValkeyUserModel{Categories: list("+@all"), Channels: list(), Commands: list("-flushall"), Keys: list("*")}

	var diags diag.Diagnostics
	plan := state
	if body := plan.update(ctx, &state, &diags); body != nil {
		t.Errorf("expected no PUT without a change, got %v", body)
	}

	// The PUT takes every list: the unknown ones keep their value.
	plan.Keys, plan.Channels = list("app:*"), types.ListUnknown(types.StringType)
	assertManagedDatabaseJSON(t, plan.update(ctx, &state, &diags), `{"categories": ["+@all"], "channels": [], "commands": ["-flushall"], "keys": ["app:*"]}`)
	if diags.HasError() {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
}

func TestCloudManagedDatabaseMongodbUserImport(t *testing.T) {
	ctx := context.Background()
	path := "/v2/publicCloud/project/pj/managedDatabase/mongodb/" + managedDatabasePostgresqlTestId + "/user/u1"
	config, _ := newManagedDatabaseTestConfig(t, map[string][]managedDatabaseTestAnswer{
		"GET " + path: {{http.StatusOK, `{"id": "u1", "username": "app@admin", "roles": ["readAnyDatabase@admin"], "status": "READY", "admin": false, "checksum": "c1", "createdAt": "2026-11-12T10:00:00Z", "credentialsLastUpdatedAt": "2026-11-12T10:00:00Z"}`}},
	})
	r := NewCloudManagedDatabaseMongodbUserResource().(*managedDatabaseUserResource[managedDatabaseRolesUserModel, *managedDatabaseRolesUserModel])
	r.config = config

	resp := fwresource.ImportStateResponse{State: managedDatabaseObjectTestState(t, r, nil)}
	r.ImportState(ctx, fwresource.ImportStateRequest{ID: "pj/" + managedDatabasePostgresqlTestId + "/u1"}, &resp)
	read := fwresource.ReadResponse{State: resp.State}
	r.Read(ctx, fwresource.ReadRequest{State: resp.State}, &read)
	if read.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", read.Diagnostics)
	}
	var got managedDatabaseRolesUserModel
	read.State.Get(ctx, &got)
	// The name the user created it with, not the username the API makes of it.
	if got.Name.ValueString() != "app" || got.Username.ValueString() != "app@admin" {
		t.Errorf("unexpected name %s and username %s", got.Name, got.Username)
	}
}

func TestCloudManagedDatabaseMongodbUserRoles(t *testing.T) {
	ctx := context.Background()
	roles := managedDatabaseMongodbUserAttributes()["roles"].(schema.SetAttribute)
	for role, valid := range map[string]bool{"readAnyDatabase@admin": true, "readWrite@app": true, "readAnyDatabase": false} {
		value := types.SetValueMust(types.StringType, []attr.Value{types.StringValue(role)})
		resp := validator.SetResponse{}
		for _, v := range roles.Validators {
			v.ValidateSet(ctx, validator.SetRequest{Path: fwpath.Root("roles"), ConfigValue: value}, &resp)
		}
		if resp.Diagnostics.HasError() == valid {
			t.Errorf("role %s: expected valid %t, got %v", role, valid, resp.Diagnostics)
		}
	}
}

// The reset answers the only copy of the new password: it stays in the state when the wait fails.
func TestCloudManagedDatabasePostgresqlUserResetThenWaitFails(t *testing.T) {
	ctx := context.Background()
	r, _ := newManagedDatabaseUserTestResource(t, map[string][]managedDatabaseTestAnswer{
		"POST " + managedDatabaseUserTestItemPath + "/resetPassword": {{http.StatusOK, `{"id": "u1", "username": "app", "password": "new-secret", "roles": ["replication"], "status": "UPDATING", "admin": false, "checksum": "c2", "createdAt": "2026-11-12T10:00:00Z", "credentialsLastUpdatedAt": "2026-11-13T10:00:00Z"}`}},
		"GET " + managedDatabaseUserTestItemPath:                     {{http.StatusOK, `{"id": "u1", "status": "ERROR"}`}},
	})
	state := managedDatabaseUserTestModel()
	plan := state
	plan.PasswordReset = types.StringValue("v2")
	plan.Password, plan.CredentialsLastUpdatedAt = types.StringUnknown(), types.StringUnknown()
	plan.Status, plan.Checksum = types.StringUnknown(), types.StringUnknown()

	// The framework answers the prior state unless Update sets one.
	resp := fwresource.UpdateResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Update(ctx, fwresource.UpdateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan), State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected the wait to fail")
	}
	var got managedDatabaseRolesUserModel
	resp.State.Get(ctx, &got)
	if got.Password.ValueString() != "new-secret" || got.PasswordReset.ValueString() != "v2" || !resp.State.Raw.IsFullyKnown() {
		t.Errorf("expected the reset in a fully known state, got password %s, password_reset %s", got.Password, got.PasswordReset)
	}
}

func TestCloudManagedDatabasePostgresqlUserDeleteAdmin(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseUserTestResource(t, map[string][]managedDatabaseTestAnswer{})
	state := managedDatabaseUserTestModel()
	state.Admin = types.BoolValue(true)
	resp := fwresource.DeleteResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
	if resp.Diagnostics.HasError() || resp.Diagnostics.WarningsCount() != 1 {
		t.Errorf("expected only a warning, got %v", resp.Diagnostics)
	}
	if requests := bodies(""); len(requests) != 0 {
		t.Errorf("expected no request for the administrator, got %v", requests)
	}
}

// A user still applying a change refuses another one with a 409 until it is READY.
func TestCloudManagedDatabasePostgresqlUserNotReady(t *testing.T) {
	ctx := context.Background()
	notReady := `{"class": "Client::Conflict::UserStatusNotReady", "message": "user \"app\" status is not ready"}`
	r, bodies := newManagedDatabaseUserTestResource(t, map[string][]managedDatabaseTestAnswer{
		"PUT " + managedDatabaseUserTestItemPath:    {{http.StatusConflict, notReady}, {http.StatusOK, managedDatabaseUserTestReady}},
		"DELETE " + managedDatabaseUserTestItemPath: {{http.StatusConflict, notReady}, {http.StatusOK, `{"id": "u1", "status": "DELETING"}`}},
		// The wait before the second PUT, the wait after it, the read, the wait before the second DELETE,
		// and the wait after it.
		"GET " + managedDatabaseUserTestItemPath: {
			{http.StatusOK, `{"id": "u1", "status": "UPDATING"}`}, {http.StatusOK, managedDatabaseUserTestReady},
			{http.StatusOK, managedDatabaseUserTestReady}, {http.StatusOK, managedDatabaseUserTestReady},
			{http.StatusOK, managedDatabaseUserTestReady}, {http.StatusNotFound, `{"message": "user not found"}`},
		},
	})
	state := managedDatabaseUserTestModel()
	plan := state
	plan.Roles = types.SetValueMust(types.StringType, []attr.Value{})
	plan.Status, plan.Checksum = types.StringUnknown(), types.StringUnknown()
	update := fwresource.UpdateResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Update(ctx, fwresource.UpdateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan), State: managedDatabaseObjectTestState(t, r, &state)}, &update)
	if update.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", update.Diagnostics)
	}
	if puts := bodies("PUT " + managedDatabaseUserTestItemPath); len(puts) != 2 {
		t.Errorf("expected the PUT again once the user is READY, got %d PUT", len(puts))
	}

	remove := fwresource.DeleteResponse{State: managedDatabaseObjectTestState(t, r, &state)}
	r.Delete(ctx, fwresource.DeleteRequest{State: managedDatabaseObjectTestState(t, r, &state)}, &remove)
	if remove.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", remove.Diagnostics)
	}
	if deletes := bodies("DELETE " + managedDatabaseUserTestItemPath); len(deletes) != 2 {
		t.Errorf("expected the DELETE again once the user is READY, got %d DELETE", len(deletes))
	}
}

// Removing password_reset sends no reset and keeps the password, also when an unknown password_reset
// in the plan turns out null.
func TestCloudManagedDatabasePostgresqlUserResetRemoved(t *testing.T) {
	ctx := context.Background()
	r, bodies := newManagedDatabaseUserTestResource(t, map[string][]managedDatabaseTestAnswer{
		"GET " + managedDatabaseUserTestItemPath: {{http.StatusOK, managedDatabaseUserTestReady}},
	})
	state := managedDatabaseUserTestModel()
	state.PasswordReset = types.StringValue("v1")
	for name, password := range map[string]types.String{"known": state.Password, "unknown": types.StringUnknown()} {
		t.Run(name, func(t *testing.T) {
			plan := state
			plan.PasswordReset, plan.Password = types.StringNull(), password
			plan.Status, plan.Checksum = types.StringUnknown(), types.StringUnknown()
			resp := fwresource.UpdateResponse{State: managedDatabaseObjectTestState(t, r, &state)}
			r.Update(ctx, fwresource.UpdateRequest{Plan: managedDatabaseObjectTestPlan(t, r, &plan), State: managedDatabaseObjectTestState(t, r, &state)}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			var got managedDatabaseRolesUserModel
			resp.State.Get(ctx, &got)
			if got.Password.ValueString() != "secret" || !got.PasswordReset.IsNull() {
				t.Errorf("expected the password kept and no password_reset, got %s and %s", got.Password, got.PasswordReset)
			}
		})
	}
	if resets := bodies("POST " + managedDatabaseUserTestItemPath + "/resetPassword"); len(resets) != 0 {
		t.Errorf("expected no reset, got %d", len(resets))
	}
}
