package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

// managedDatabaseUserKind is the model of a kind of user: the attributes every user has, and the fields
// of its engines, which it converts.
type managedDatabaseUserKind[M any] interface {
	*M
	user() *managedDatabaseUserModel
	// creation returns the body of the POST.
	creation(ctx context.Context, diags *diag.Diagnostics) any
	// update returns the body of the PUT, with every field of the plan, nil when none changes.
	update(ctx context.Context, state *M, diags *diag.Diagnostics) map[string]any
	// mergeFields reads the fields of the engines from an answer of the API.
	mergeFields(ctx context.Context, answer json.RawMessage, diags *diag.Diagnostics)
}

var (
	_ resource.ResourceWithConfigure   = (*managedDatabaseUserResource[managedDatabasePlainUserModel, *managedDatabasePlainUserModel])(nil)
	_ resource.ResourceWithImportState = (*managedDatabaseUserResource[managedDatabasePlainUserModel, *managedDatabasePlainUserModel])(nil)
	_ resource.ResourceWithModifyPlan  = (*managedDatabaseUserResource[managedDatabasePlainUserModel, *managedDatabasePlainUserModel])(nil)
)

func NewCloudManagedDatabaseUserResource() resource.Resource {
	return &managedDatabaseUserResource[managedDatabasePlainUserModel, *managedDatabasePlainUserModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"mysql"}, suffix: "_user"},
	}
}

func NewCloudManagedDatabasePostgresqlUserResource() resource.Resource {
	return &managedDatabaseUserResource[managedDatabaseRolesUserModel, *managedDatabaseRolesUserModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"postgresql"}, suffix: "_postgresql_user"},
		fields:                        managedDatabaseRolesUserAttributes(true),
	}
}

func NewCloudManagedDatabaseMongodbUserResource() resource.Resource {
	return &managedDatabaseUserResource[managedDatabaseRolesUserModel, *managedDatabaseRolesUserModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"mongodb"}, suffix: "_mongodb_user"},
		fields:                        managedDatabaseMongodbUserAttributes(),
	}
}

func NewCloudManagedDatabaseValkeyUserResource() resource.Resource {
	return &managedDatabaseUserResource[managedDatabaseValkeyUserModel, *managedDatabaseValkeyUserModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedDatabaseFamily, engines: []string{"valkey"}, suffix: "_valkey_user"},
		fields:                        managedDatabaseValkeyUserAttributes(),
	}
}

func NewCloudManagedAnalyticsUserResource() resource.Resource {
	return &managedDatabaseUserResource[managedDatabasePlainUserModel, *managedDatabasePlainUserModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedAnalyticsFamily, engines: []string{"kafka", "kafkaConnect"}, suffix: "_user"},
	}
}

func NewCloudManagedAnalyticsClickhouseUserResource() resource.Resource {
	return &managedDatabaseUserResource[managedDatabaseRolesUserModel, *managedDatabaseRolesUserModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedAnalyticsFamily, engines: []string{"clickhouse"}, suffix: "_clickhouse_user"},
		fields:                        managedDatabaseRolesUserAttributes(false),
	}
}

func NewCloudManagedAnalyticsOpensearchUserResource() resource.Resource {
	return &managedDatabaseUserResource[managedDatabaseOpensearchUserModel, *managedDatabaseOpensearchUserModel]{
		managedDatabaseObjectResource: managedDatabaseObjectResource{category: managedAnalyticsFamily, engines: []string{"opensearch"}, suffix: "_opensearch_user"},
		fields:                        managedDatabaseOpensearchUserAttributes(),
	}
}

// managedDatabaseRolesUserAttributes returns the roles of a user, which the engine takes, or only
// returns for ClickHouse.
func managedDatabaseRolesUserAttributes(settable bool) map[string]schema.Attribute {
	roles := schema.SetAttribute{
		ElementType:   types.StringType,
		Computed:      true,
		Description:   "Roles of the user.",
		PlanModifiers: []planmodifier.Set{UseStateForUnknownSet()},
	}
	if settable {
		roles.Optional = true
		roles.Description = "Roles assigned to the user. Available roles are returned by the roles endpoint of the service. Removing it from the configuration keeps the current roles: set [] to remove them all."
	}
	return map[string]schema.Attribute{"roles": roles}
}

// managedDatabaseMongodbUserAttributes returns the roles of a MongoDB user, each with its
// authentication database: the API adds @admin to a role without one, which the plan then doesn't match.
func managedDatabaseMongodbUserAttributes() map[string]schema.Attribute {
	attributes := managedDatabaseRolesUserAttributes(true)
	roles := attributes["roles"].(schema.SetAttribute)
	roles.Description = "Roles assigned to the user, each with its authentication database, such as readAnyDatabase@admin. Available roles are returned by the roles endpoint of the service. Removing it from the configuration keeps the current roles: set [] to remove them all."
	roles.Validators = []validator.Set{setvalidator.ValueStringsAre(stringvalidator.RegexMatches(
		regexp.MustCompile(`^[\w-]+@[\w-]+$`), "must be a role and its authentication database, such as readAnyDatabase@admin"))}
	attributes["roles"] = roles
	return attributes
}

func managedDatabaseValkeyUserAttributes() map[string]schema.Attribute {
	list := func(description string) schema.ListAttribute {
		return schema.ListAttribute{
			ElementType:   types.StringType,
			Optional:      true,
			Computed:      true,
			Description:   description + " Removing it from the configuration keeps the current list: set [] to empty it.",
			PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
		}
	}
	return map[string]schema.Attribute{
		"categories": list("Command categories the user is allowed or denied, such as +@all."),
		"channels":   list("Pub/Sub channels the user can access."),
		"commands":   list("Commands the user is allowed or denied, such as +get."),
		"keys":       list("Key patterns the user can access."),
	}
}

func managedDatabaseOpensearchUserAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"acls": schema.SetNestedAttribute{
			Optional:      true,
			Computed:      true,
			Description:   "Access rules of the user. Removing it from the configuration keeps the current rules: set [] to remove them all.",
			PlanModifiers: []planmodifier.Set{UseStateForUnknownSet()},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"pattern": schema.StringAttribute{
						Required:    true,
						Description: "Pattern of the indexes the rule applies to.",
					},
					"permission": schema.StringAttribute{
						Required:    true,
						Description: "Permission on the indexes: ADMIN, DENY, READ, READ_WRITE or WRITE.",
						Validators:  []validator.String{stringvalidator.OneOf("ADMIN", "DENY", "READ", "READ_WRITE", "WRITE")},
					},
				},
			},
		},
	}
}

// managedDatabaseUserResource is a user of a service. Its kind M holds the fields of its engines.
type managedDatabaseUserResource[M any, P managedDatabaseUserKind[M]] struct {
	managedDatabaseObjectResource
	fields map[string]schema.Attribute
}

func (r *managedDatabaseUserResource[M, P]) endpoint(user *managedDatabaseUserModel) string {
	return r.serviceEndpoint(user.ServiceName, user.Engine, user.ServiceId) + "/user"
}

func (r *managedDatabaseUserResource[M, P]) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := r.objectAttributes("user")
	attributes["name"] = schema.StringAttribute{
		Required:      true,
		Description:   "Login name of the user. Changing it recreates the user.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	}
	attributes["password_reset"] = schema.StringAttribute{
		Optional:    true,
		Description: "Arbitrary string to change to reset the password. Removing it keeps the password.",
	}
	attributes["username"] = schema.StringAttribute{
		Computed:      true,
		Description:   "Login name of the user.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	attributes["password"] = schema.StringAttribute{
		Computed:      true,
		Sensitive:     true,
		Description:   "Password of the user. The API returns it only at the creation and after a reset.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	attributes["admin"] = schema.BoolAttribute{
		Computed:      true,
		Description:   "Whether the user is the administrator of the service, created with it, which can't be deleted.",
		PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
	}
	attributes["status"] = schema.StringAttribute{
		Computed:    true,
		Description: "Readiness of the user: CREATING, UPDATING, READY, ERROR or DELETING.",
	}
	attributes["checksum"] = schema.StringAttribute{
		Computed:    true,
		Description: "Hash of the user, computed by the API.",
	}
	attributes["created_at"] = schema.StringAttribute{
		Computed:      true,
		Description:   "Creation date of the user.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	attributes["credentials_last_updated_at"] = schema.StringAttribute{
		Computed:      true,
		Description:   "Date of the last change of the password.",
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
	for name, attribute := range r.fields {
		attributes[name] = attribute
	}

	resp.Schema = schema.Schema{
		Description: r.describe("a user"),
		Attributes:  attributes,
	}
}

// managedDatabasePasswordResets tells whether the planned password_reset resets the password: a value
// other than the current one. Removing it keeps the password, which clients use. An unknown value may
// reset it.
func managedDatabasePasswordResets(planned, current types.String) bool {
	return !planned.IsNull() && !planned.Equal(current)
}

// ModifyPlan announces a new password only when password_reset resets it, not for a change of the other
// fields.
func (r *managedDatabaseUserResource[M, P]) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var planned, current types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("password_reset"), &planned)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("password_reset"), &current)...)
	if !managedDatabasePasswordResets(planned, current) {
		return
	}
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("password"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("credentials_last_updated_at"), types.StringUnknown())...)
}

func (r *managedDatabaseUserResource[M, P]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	r.importState(ctx, req, resp, "id")
}

func (r *managedDatabaseUserResource[M, P]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data M
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	model := P(&data)
	body := model.creation(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := r.endpoint(model.user())
	var answer json.RawMessage
	if err := retryWhileManagedDatabaseServiceBusy(ctx, managedDatabaseObjectTimeout, func() error {
		return r.config.OVHClient.PostWithContext(ctx, endpoint, body, &answer)
	}); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Saved right away: the password is only in this answer.
	r.merge(ctx, model, answer, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint += "/" + url.PathEscape(model.user().Id.ValueString())
	if err := waitForManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, false, managedDatabaseObjectTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for the user to be ready", err.Error())
		return
	}
	r.read(ctx, model, endpoint, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *managedDatabaseUserResource[M, P]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data M
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	model := P(&data)

	endpoint := r.endpoint(model.user()) + "/" + url.PathEscape(model.user().Id.ValueString())
	var answer json.RawMessage
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &answer); err != nil {
		if helpers.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	r.merge(ctx, model, answer, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *managedDatabaseUserResource[M, P]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state M
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	model, current := P(&plan), P(&state)

	endpoint := r.endpoint(current.user()) + "/" + url.PathEscape(current.user().Id.ValueString())
	body := model.update(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if body != nil {
		if err := retryWhileManagedDatabaseObjectBusy(ctx, r.config.OVHClient, endpoint, managedDatabaseObjectTimeout, func() error {
			return r.config.OVHClient.PutWithContext(ctx, endpoint, body, nil)
		}); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
			return
		}
		if err := waitForManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, false, managedDatabaseObjectTimeout); err != nil {
			resp.Diagnostics.AddError("Error waiting for the user to be ready", err.Error())
			return
		}
	}

	if managedDatabasePasswordResets(model.user().PasswordReset, current.user().PasswordReset) {
		var answer json.RawMessage
		if err := retryWhileManagedDatabaseObjectBusy(ctx, r.config.OVHClient, endpoint, managedDatabaseObjectTimeout, func() error {
			return r.config.OVHClient.PostWithContext(ctx, endpoint+"/resetPassword", nil, &answer)
		}); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s/resetPassword", endpoint), err.Error())
			return
		}
		// Saved right away: the new password is only in this answer.
		r.merge(ctx, model, answer, &resp.Diagnostics)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := waitForManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, false, managedDatabaseObjectTimeout); err != nil {
			resp.Diagnostics.AddError("Error waiting for the user to be ready", err.Error())
			return
		}
	} else if model.user().Password.IsUnknown() {
		// A password_reset unknown in the plan, which could have reset the password, turned out null.
		model.user().Password = current.user().Password
	}

	r.read(ctx, model, endpoint, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *managedDatabaseUserResource[M, P]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data M
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	user := P(&data).user()
	if user.Admin.ValueBool() {
		resp.Diagnostics.AddWarning("Administrator not deleted",
			"The API can't delete the administrator of the service, created with it: "+user.Username.ValueString()+" is only removed from the state.")
		return
	}

	endpoint := r.endpoint(user) + "/" + url.PathEscape(user.Id.ValueString())
	if err := deleteManagedDatabaseObject(ctx, r.config.OVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}
	if err := waitForManagedDatabaseObject(ctx, r.config.OVHClient, endpoint, true, managedDatabaseObjectTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for the user to be deleted", err.Error())
	}
}

// merge copies an answer of the API into the model: the attributes every user has, then the fields of
// its kind.
func (r *managedDatabaseUserResource[M, P]) merge(ctx context.Context, model P, answer json.RawMessage, diags *diag.Diagnostics) {
	var user managedDatabaseUser
	if err := json.Unmarshal(answer, &user); err != nil {
		diags.AddError("Error decoding the user", err.Error())
		return
	}
	model.user().mergeWith(&user)
	model.mergeFields(ctx, answer, diags)
}

func (r *managedDatabaseUserResource[M, P]) read(ctx context.Context, model P, endpoint string, diags *diag.Diagnostics) {
	var answer json.RawMessage
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &answer); err != nil {
		diags.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	r.merge(ctx, model, answer, diags)
}
