package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

// Data sources of the objects of the managed database and managed analytics services. They decode the
// API models, whose tfsdk tags map them onto the attributes: a singular data source reads one object, a
// plural one lists them.

var (
	_ datasource.DataSourceWithConfigure = (*managedDatabaseObjectDataSource[managedDatabaseUser])(nil)
	_ datasource.DataSourceWithConfigure = (*managedDatabaseObjectsDataSource[managedDatabaseUser])(nil)
)

// managedDatabaseRole is the publicCloud.database.Role model.
type managedDatabaseRole struct {
	Checksum  string `json:"checksum" tfsdk:"checksum"`
	CreatedAt string `json:"createdAt" tfsdk:"created_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updatedAt" tfsdk:"updated_at"`
}

// managedDatabaseObjectKind describes an object of a service for its data sources: the engines that have
// it, the path of its route, its name and its attributes.
type managedDatabaseObjectKind struct {
	category   managedDatabaseCategory
	engines    []string
	suffix     string
	path       string
	object     string
	objects    string
	attributes func() map[string]schema.Attribute
}

func newManagedDatabaseObjectDataSources[T any](kind managedDatabaseObjectKind) (func() datasource.DataSource, func() datasource.DataSource) {
	return func() datasource.DataSource {
			return &managedDatabaseObjectDataSource[T]{managedDatabaseObjectKind: kind}
		}, func() datasource.DataSource {
			plural := kind
			plural.suffix += "s"
			return &managedDatabaseObjectsDataSource[T]{managedDatabaseObjectKind: plural}
		}
}

var (
	NewCloudManagedDatabaseUserDataSource, NewCloudManagedDatabaseUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseUser](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: []string{"mysql"}, suffix: "_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseUserDataAttributes})
	NewCloudManagedDatabasePostgresqlUserDataSource, NewCloudManagedDatabasePostgresqlUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseUserWithRoles](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: []string{"postgresql"}, suffix: "_postgresql_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseRolesUserDataAttributes})
	NewCloudManagedDatabaseMongodbUserDataSource, NewCloudManagedDatabaseMongodbUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseUserWithRoles](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: []string{"mongodb"}, suffix: "_mongodb_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseRolesUserDataAttributes})
	NewCloudManagedDatabaseValkeyUserDataSource, NewCloudManagedDatabaseValkeyUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseValkeyUser](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: []string{"valkey"}, suffix: "_valkey_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseValkeyUserDataAttributes})
	NewCloudManagedDatabaseRoleDataSource, NewCloudManagedDatabaseRolesDataSource = newManagedDatabaseObjectDataSources[managedDatabaseRole](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: []string{"mongodb", "postgresql"}, suffix: "_role", path: "/roles", object: "role", objects: "roles", attributes: managedDatabaseRoleDataAttributes})

	NewCloudManagedAnalyticsUserDataSource, NewCloudManagedAnalyticsUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseUser](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: []string{"grafana", "kafka", "kafkaConnect"}, suffix: "_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseUserDataAttributes})
	NewCloudManagedAnalyticsClickhouseUserDataSource, NewCloudManagedAnalyticsClickhouseUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseUserWithRoles](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: []string{"clickhouse"}, suffix: "_clickhouse_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseRolesUserDataAttributes})
	NewCloudManagedAnalyticsOpensearchUserDataSource, NewCloudManagedAnalyticsOpensearchUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseOpensearchUser](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: []string{"opensearch"}, suffix: "_opensearch_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseOpensearchUserDataAttributes})
	NewCloudManagedAnalyticsClickhouseRoleDataSource, NewCloudManagedAnalyticsClickhouseRolesDataSource = newManagedDatabaseObjectDataSources[managedDatabaseRole](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: []string{"clickhouse"}, suffix: "_clickhouse_role", path: "/roles", object: "role", objects: "roles", attributes: managedDatabaseRoleDataAttributes})
)

func managedDatabaseUserDataAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":                          managedDatabaseComputedString("ID of the user."),
		"username":                    managedDatabaseComputedString("Login name of the user."),
		"admin":                       schema.BoolAttribute{Computed: true, Description: "Whether the user is the administrator of the service, created with it, which can't be deleted."},
		"status":                      managedDatabaseComputedString("Readiness of the user: CREATING, UPDATING, READY, ERROR or DELETING."),
		"checksum":                    managedDatabaseComputedString("Hash of the user, computed by the API."),
		"created_at":                  managedDatabaseComputedString("Creation date of the user."),
		"credentials_last_updated_at": managedDatabaseComputedString("Date of the last change of the password."),
	}
}

func managedDatabaseRolesUserDataAttributes() map[string]schema.Attribute {
	attributes := managedDatabaseUserDataAttributes()
	attributes["roles"] = schema.SetAttribute{ElementType: types.StringType, Computed: true, Description: "Roles of the user."}
	return attributes
}

func managedDatabaseValkeyUserDataAttributes() map[string]schema.Attribute {
	attributes := managedDatabaseUserDataAttributes()
	for name, description := range map[string]string{
		"categories": "Command categories the user is allowed or denied.",
		"channels":   "Pub/Sub channels the user can access.",
		"commands":   "Commands the user is allowed or denied.",
		"keys":       "Key patterns the user can access.",
	} {
		attributes[name] = schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: description}
	}
	return attributes
}

func managedDatabaseOpensearchUserDataAttributes() map[string]schema.Attribute {
	attributes := managedDatabaseUserDataAttributes()
	attributes["acls"] = schema.SetNestedAttribute{
		Computed:    true,
		Description: "Access rules of the user.",
		NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"pattern":    managedDatabaseComputedString("Pattern of the indexes the rule applies to."),
			"permission": managedDatabaseComputedString("Permission on the indexes: ADMIN, DENY, READ, READ_WRITE or WRITE."),
		}},
	}
	return attributes
}

func managedDatabaseRoleDataAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":         managedDatabaseComputedString("ID of the role."),
		"name":       managedDatabaseComputedString("Name of the role."),
		"checksum":   managedDatabaseComputedString("Hash of the role, computed by the API."),
		"created_at": managedDatabaseComputedString("Creation date of the role."),
		"updated_at": managedDatabaseComputedString("Last update date of the role."),
	}
}

// managedDatabaseDataSourceConfigured gives a data source the provider configuration.
type managedDatabaseDataSourceConfigured struct {
	config *Config
}

func (d *managedDatabaseDataSourceConfigured) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

func managedDatabaseComputedString(description string) schema.StringAttribute {
	return schema.StringAttribute{Computed: true, Description: description}
}

func managedDatabaseAttrTypes(attributes map[string]schema.Attribute) map[string]attr.Type {
	attrTypes := make(map[string]attr.Type, len(attributes))
	for name, attribute := range attributes {
		attrTypes[name] = attribute.GetType()
	}
	return attrTypes
}

func (k *managedDatabaseObjectKind) named() bool {
	return managedDatabaseEngineInName(k.engines, k.suffix)
}

// describe returns the description of a data source of an object of the services of the engines.
func (k *managedDatabaseObjectKind) describe(what string) string {
	return "Use this data source to " + what + " of a managed " + k.category.family + " service (" + strings.Join(k.engines, ", ") + "), using the Cloud Databases API v2."
}

// serviceAttributes returns the arguments that find the service: the project, the engine and the service.
func (k *managedDatabaseObjectKind) serviceAttributes() map[string]schema.Attribute {
	engine := schema.StringAttribute{Computed: true, Description: "Engine of the service: " + k.engines[0] + "."}
	if !k.named() {
		engine = schema.StringAttribute{
			Required:    true,
			Description: "Engine of the service: " + strings.Join(k.engines, ", ") + ".",
			Validators:  []validator.String{stringvalidator.OneOf(k.engines...)},
		}
	}
	return map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "ID of the public cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used.",
		},
		"engine":     engine,
		"service_id": schema.StringAttribute{Required: true, Description: "ID of the service."},
	}
}

// serviceEndpoint returns the path of the service of the configuration, and sets in the state the project,
// which can come from the environment, and the engine when the type name holds it.
func (k *managedDatabaseObjectKind) serviceEndpoint(ctx context.Context, config tfsdk.Config, resp *datasource.ReadResponse) string {
	var serviceName, engine, serviceId types.String
	resp.Diagnostics.Append(config.GetAttribute(ctx, path.Root("service_name"), &serviceName)...)
	resp.Diagnostics.Append(config.GetAttribute(ctx, path.Root("service_id"), &serviceId)...)
	serviceName = cloudProjectServiceNameStringOrEnv(serviceName, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), serviceName)...)
	if k.named() {
		engine = types.StringValue(k.engines[0])
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("engine"), engine)...)
	} else {
		resp.Diagnostics.Append(config.GetAttribute(ctx, path.Root("engine"), &engine)...)
	}
	return k.category.serviceEndpoint(serviceName.ValueString(), engine.ValueString(), serviceId.ValueString())
}

// managedDatabaseObjectDataSource reads one object of a service, found by its id.
type managedDatabaseObjectDataSource[T any] struct {
	managedDatabaseDataSourceConfigured
	managedDatabaseObjectKind
}

func (d *managedDatabaseObjectDataSource[T]) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.category.typeName(d.suffix)
}

func (d *managedDatabaseObjectDataSource[T]) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := d.serviceAttributes()
	for name, attribute := range d.attributes() {
		attributes[name] = attribute
	}
	attributes["id"] = schema.StringAttribute{Required: true, Description: "ID of the " + d.object + "."}
	resp.Schema = schema.Schema{
		Description: d.describe("get a " + d.object),
		Attributes:  attributes,
	}
}

func (d *managedDatabaseObjectDataSource[T]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	service := d.serviceEndpoint(ctx, req.Config, resp)
	var id types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := service + d.path + "/" + url.PathEscape(id.ValueString())
	var object T
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &object); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	value, diags := types.ObjectValueFrom(ctx, managedDatabaseAttrTypes(d.attributes()), object)
	resp.Diagnostics.Append(diags...)
	for name, attribute := range value.Attributes() {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(name), attribute)...)
	}
}

// managedDatabaseObjectsDataSource lists the objects of a service.
type managedDatabaseObjectsDataSource[T any] struct {
	managedDatabaseDataSourceConfigured
	managedDatabaseObjectKind
}

func (d *managedDatabaseObjectsDataSource[T]) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.category.typeName(d.suffix)
}

func (d *managedDatabaseObjectsDataSource[T]) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	noun := strings.ReplaceAll(d.objects, "_", " ")
	attributes := d.serviceAttributes()
	attributes[d.objects] = schema.ListNestedAttribute{
		Computed:     true,
		Description:  "The " + noun + " of the service.",
		NestedObject: schema.NestedAttributeObject{Attributes: d.attributes()},
	}
	resp.Schema = schema.Schema{
		Description: d.describe("list the " + noun),
		Attributes:  attributes,
	}
}

func (d *managedDatabaseObjectsDataSource[T]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	endpoint := d.serviceEndpoint(ctx, req.Config, resp) + d.path
	if resp.Diagnostics.HasError() {
		return
	}

	objects, err := helpers.GetAllPagesV2[T](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	if objects == nil {
		objects = []T{}
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: managedDatabaseAttrTypes(d.attributes())}, objects)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(d.objects), list)...)
}
