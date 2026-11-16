package ovh

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// managedDatabaseUserModel holds the attributes every user has. The model of each kind of user embeds
// it, with the fields of its engines.
type managedDatabaseUserModel struct {
	ServiceName   types.String `tfsdk:"service_name"`
	Engine        types.String `tfsdk:"engine"`
	ServiceId     types.String `tfsdk:"service_id"`
	Name          types.String `tfsdk:"name"`
	PasswordReset types.String `tfsdk:"password_reset"`

	Id                       types.String `tfsdk:"id"`
	Username                 types.String `tfsdk:"username"`
	Password                 types.String `tfsdk:"password"`
	Admin                    types.Bool   `tfsdk:"admin"`
	Status                   types.String `tfsdk:"status"`
	Checksum                 types.String `tfsdk:"checksum"`
	CreatedAt                types.String `tfsdk:"created_at"`
	CredentialsLastUpdatedAt types.String `tfsdk:"credentials_last_updated_at"`
}

// managedDatabaseUser is the publicCloud.database.User model, which the users of every engine extend.
// Its tfsdk tags map it onto the attributes of the user data sources, which don't expose the password.
type managedDatabaseUser struct {
	Admin                    bool    `json:"admin" tfsdk:"admin"`
	Checksum                 string  `json:"checksum" tfsdk:"checksum"`
	CreatedAt                string  `json:"createdAt" tfsdk:"created_at"`
	CredentialsLastUpdatedAt string  `json:"credentialsLastUpdatedAt" tfsdk:"credentials_last_updated_at"`
	Id                       string  `json:"id" tfsdk:"id"`
	Password                 *string `json:"password" tfsdk:"-"`
	Status                   string  `json:"status" tfsdk:"status"`
	Username                 string  `json:"username" tfsdk:"username"`
}

// managedDatabaseUserWithRoles is the publicCloud.database.UserWithRoles model.
type managedDatabaseUserWithRoles struct {
	managedDatabaseUser
	Roles []string `json:"roles" tfsdk:"roles"`
}

// managedDatabaseValkeyUser is the publicCloud.database.valkey.User model.
type managedDatabaseValkeyUser struct {
	managedDatabaseUser
	Categories []string `json:"categories" tfsdk:"categories"`
	Channels   []string `json:"channels" tfsdk:"channels"`
	Commands   []string `json:"commands" tfsdk:"commands"`
	Keys       []string `json:"keys" tfsdk:"keys"`
}

type managedDatabaseOpensearchUserAcl struct {
	Pattern    string `json:"pattern" tfsdk:"pattern"`
	Permission string `json:"permission" tfsdk:"permission"`
}

// managedDatabaseOpensearchUser is the publicCloud.database.opensearch.User model.
type managedDatabaseOpensearchUser struct {
	managedDatabaseUser
	Acls []managedDatabaseOpensearchUserAcl `json:"acls" tfsdk:"acls"`
}

// mergeWith copies the attributes every user has. The API returns the password only at the creation
// and after a reset: otherwise the model keeps its own.
func (m *managedDatabaseUserModel) mergeWith(user *managedDatabaseUser) {
	m.Id = types.StringValue(user.Id)
	m.Username = types.StringValue(user.Username)
	if m.Name.IsNull() || m.Name.IsUnknown() {
		name := user.Username
		if m.Engine.ValueString() == "mongodb" {
			// The API names a MongoDB user <name>@admin, and refuses an @ in the name it creates.
			name, _, _ = strings.Cut(name, "@")
		}
		m.Name = types.StringValue(name)
	}
	switch {
	case user.Password != nil:
		m.Password = types.StringValue(*user.Password)
	case m.Password.IsUnknown():
		m.Password = types.StringNull()
	}
	m.Admin = types.BoolValue(user.Admin)
	m.Status = types.StringValue(user.Status)
	m.Checksum = types.StringValue(user.Checksum)
	m.CreatedAt = types.StringValue(user.CreatedAt)
	m.CredentialsLastUpdatedAt = types.StringValue(user.CredentialsLastUpdatedAt)
}

// knownStrings returns the elements of a known list or set, nil otherwise.
func knownStrings(ctx context.Context, value interface {
	IsNull() bool
	IsUnknown() bool
	ElementsAs(context.Context, interface{}, bool) diag.Diagnostics
}, diags *diag.Diagnostics) []string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	elements := []string{}
	diags.Append(value.ElementsAs(ctx, &elements, false)...)
	return elements
}

// The kinds of users. Each one tells what the POST and the PUT send, and reads its fields from an
// answer of the API.

// managedDatabasePlainUserModel is a user without fields of its own: MySQL, Kafka and Kafka Connect.
type managedDatabasePlainUserModel struct {
	managedDatabaseUserModel
}

func (m *managedDatabasePlainUserModel) user() *managedDatabaseUserModel {
	return &m.managedDatabaseUserModel
}

func (m *managedDatabasePlainUserModel) creation(ctx context.Context, diags *diag.Diagnostics) any {
	return map[string]any{"name": m.Name.ValueString()}
}

func (m *managedDatabasePlainUserModel) update(ctx context.Context, state *managedDatabasePlainUserModel, diags *diag.Diagnostics) map[string]any {
	return nil
}

func (m *managedDatabasePlainUserModel) mergeFields(ctx context.Context, answer json.RawMessage, diags *diag.Diagnostics) {
}

// managedDatabaseRolesUserModel is a user with roles: PostgreSQL and MongoDB, which take them, and
// ClickHouse, which only returns them.
type managedDatabaseRolesUserModel struct {
	managedDatabaseUserModel
	Roles types.Set `tfsdk:"roles"`
}

func (m *managedDatabaseRolesUserModel) user() *managedDatabaseUserModel {
	return &m.managedDatabaseUserModel
}

func (m *managedDatabaseRolesUserModel) creation(ctx context.Context, diags *diag.Diagnostics) any {
	body := map[string]any{"name": m.Name.ValueString()}
	if roles := knownStrings(ctx, m.Roles, diags); roles != nil {
		body["roles"] = roles
	}
	return body
}

func (m *managedDatabaseRolesUserModel) update(ctx context.Context, state *managedDatabaseRolesUserModel, diags *diag.Diagnostics) map[string]any {
	if m.Roles.IsUnknown() || m.Roles.Equal(state.Roles) {
		return nil
	}
	roles := knownStrings(ctx, m.Roles, diags)
	if roles == nil {
		roles = []string{}
	}
	return map[string]any{"roles": roles}
}

func (m *managedDatabaseRolesUserModel) mergeFields(ctx context.Context, answer json.RawMessage, diags *diag.Diagnostics) {
	var user managedDatabaseUserWithRoles
	if err := json.Unmarshal(answer, &user); err != nil {
		diags.AddError("Error decoding the user", err.Error())
		return
	}
	if user.Roles == nil {
		user.Roles = []string{}
	}
	var d diag.Diagnostics
	m.Roles, d = types.SetValueFrom(ctx, types.StringType, user.Roles)
	diags.Append(d...)
}

// managedDatabaseValkeyUserModel is a Valkey user, with its access control fields.
type managedDatabaseValkeyUserModel struct {
	managedDatabaseUserModel
	Categories types.List `tfsdk:"categories"`
	Channels   types.List `tfsdk:"channels"`
	Commands   types.List `tfsdk:"commands"`
	Keys       types.List `tfsdk:"keys"`
}

func (m *managedDatabaseValkeyUserModel) user() *managedDatabaseUserModel {
	return &m.managedDatabaseUserModel
}

func (m *managedDatabaseValkeyUserModel) fields() map[string]types.List {
	return map[string]types.List{"categories": m.Categories, "channels": m.Channels, "commands": m.Commands, "keys": m.Keys}
}

func (m *managedDatabaseValkeyUserModel) creation(ctx context.Context, diags *diag.Diagnostics) any {
	body := map[string]any{"name": m.Name.ValueString()}
	for name, value := range m.fields() {
		if elements := knownStrings(ctx, value, diags); elements != nil {
			body[name] = elements
		}
	}
	return body
}

func (m *managedDatabaseValkeyUserModel) update(ctx context.Context, state *managedDatabaseValkeyUserModel, diags *diag.Diagnostics) map[string]any {
	body, changed := map[string]any{}, false
	current := state.fields()
	for name, value := range m.fields() {
		if value.IsUnknown() {
			value = current[name]
		}
		changed = changed || !value.Equal(current[name])
		elements := knownStrings(ctx, value, diags)
		if elements == nil {
			elements = []string{}
		}
		body[name] = elements
	}
	if !changed {
		return nil
	}
	return body
}

func (m *managedDatabaseValkeyUserModel) mergeFields(ctx context.Context, answer json.RawMessage, diags *diag.Diagnostics) {
	var user managedDatabaseValkeyUser
	if err := json.Unmarshal(answer, &user); err != nil {
		diags.AddError("Error decoding the user", err.Error())
		return
	}
	list := func(elements []string) types.List {
		if elements == nil {
			elements = []string{}
		}
		value, d := types.ListValueFrom(ctx, types.StringType, elements)
		diags.Append(d...)
		return value
	}
	m.Categories, m.Channels, m.Commands, m.Keys = list(user.Categories), list(user.Channels), list(user.Commands), list(user.Keys)
}

// managedDatabaseOpensearchUserModel is an OpenSearch user, with its ACLs.
type managedDatabaseOpensearchUserModel struct {
	managedDatabaseUserModel
	Acls types.Set `tfsdk:"acls"`
}

func (m *managedDatabaseOpensearchUserModel) user() *managedDatabaseUserModel {
	return &m.managedDatabaseUserModel
}

func (m *managedDatabaseOpensearchUserModel) acls(ctx context.Context, diags *diag.Diagnostics) []managedDatabaseOpensearchUserAcl {
	acls := []managedDatabaseOpensearchUserAcl{}
	if !m.Acls.IsNull() && !m.Acls.IsUnknown() {
		diags.Append(m.Acls.ElementsAs(ctx, &acls, false)...)
	}
	return acls
}

func (m *managedDatabaseOpensearchUserModel) creation(ctx context.Context, diags *diag.Diagnostics) any {
	body := map[string]any{"name": m.Name.ValueString()}
	if !m.Acls.IsNull() && !m.Acls.IsUnknown() {
		body["acls"] = m.acls(ctx, diags)
	}
	return body
}

func (m *managedDatabaseOpensearchUserModel) update(ctx context.Context, state *managedDatabaseOpensearchUserModel, diags *diag.Diagnostics) map[string]any {
	if m.Acls.IsUnknown() || m.Acls.Equal(state.Acls) {
		return nil
	}
	return map[string]any{"acls": m.acls(ctx, diags)}
}

func (m *managedDatabaseOpensearchUserModel) mergeFields(ctx context.Context, answer json.RawMessage, diags *diag.Diagnostics) {
	var user managedDatabaseOpensearchUser
	if err := json.Unmarshal(answer, &user); err != nil {
		diags.AddError("Error decoding the user", err.Error())
		return
	}
	if user.Acls == nil {
		user.Acls = []managedDatabaseOpensearchUserAcl{}
	}
	var d diag.Diagnostics
	m.Acls, d = types.SetValueFrom(ctx, managedDatabaseOpensearchUserAttributes()["acls"].GetType().(types.SetType).ElemType, user.Acls)
	diags.Append(d...)
}
