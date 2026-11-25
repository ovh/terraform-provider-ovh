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
	_ datasource.DataSourceWithConfigure = (*managedDatabaseNamesDataSource)(nil)
)

// managedDatabaseRole is the publicCloud.database.Role model.
type managedDatabaseRole struct {
	Checksum  string `json:"checksum" tfsdk:"checksum"`
	CreatedAt string `json:"createdAt" tfsdk:"created_at"`
	Id        string `json:"id" tfsdk:"id"`
	Name      string `json:"name" tfsdk:"name"`
	UpdatedAt string `json:"updatedAt" tfsdk:"updated_at"`
}

// managedDatabaseLogKind is the dbaas.logs.LogKind model.
type managedDatabaseLogKind struct {
	AdditionalReturnedFields []string `json:"additionalReturnedFields" tfsdk:"additional_returned_fields"`
	CreatedAt                string   `json:"createdAt" tfsdk:"created_at"`
	DisplayName              string   `json:"displayName" tfsdk:"display_name"`
	KindId                   string   `json:"kindId" tfsdk:"kind_id"`
	Name                     string   `json:"name" tfsdk:"name"`
	UpdatedAt                string   `json:"updatedAt" tfsdk:"updated_at"`
}

var (
	managedDatabaseAllEngines  = managedDatabaseFamily.engines
	managedAnalyticsAllEngines = managedAnalyticsFamily.engines
)

// managedDatabaseObjectKind describes an object of a service for its data sources: the engines that have
// it, the path of its route, its name and its attributes.
type managedDatabaseObjectKind struct {
	category   managedDatabaseCategory
	engines    []string
	suffix     string
	path       string
	object     string
	objects    string
	key        string
	attributes func() map[string]schema.Attribute
}

func newManagedDatabaseObjectDataSources[T any](kind managedDatabaseObjectKind) (func() datasource.DataSource, func() datasource.DataSource) {
	if kind.key == "" {
		kind.key = "id"
	}
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
	NewCloudManagedDatabaseDatabaseDataSource, NewCloudManagedDatabaseDatabasesDataSource = newManagedDatabaseObjectDataSources[managedDatabaseDatabase](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: []string{"mysql", "postgresql"}, suffix: "_database", path: "/database", object: "database", objects: "databases", attributes: managedDatabaseDatabaseDataAttributes})
	NewCloudManagedDatabasePostgresqlConnectionPoolDataSource, NewCloudManagedDatabasePostgresqlConnectionPoolsDataSource = newManagedDatabaseObjectDataSources[managedDatabaseConnectionPool](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: []string{"postgresql"}, suffix: "_postgresql_connection_pool", path: "/connectionPool", object: "connection pool", objects: "connection_pools", attributes: managedDatabaseConnectionPoolDataAttributes})
	NewCloudManagedDatabaseRoleDataSource, NewCloudManagedDatabaseRolesDataSource = newManagedDatabaseObjectDataSources[managedDatabaseRole](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: []string{"mongodb", "postgresql"}, suffix: "_role", path: "/roles", object: "role", objects: "roles", attributes: managedDatabaseRoleDataAttributes})
	NewCloudManagedDatabaseLogSubscriptionDataSource, _ = newManagedDatabaseObjectDataSources[managedDatabaseLogSubscription](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: managedDatabaseAllEngines, suffix: "_log_subscription", path: "/log/subscription", object: "log subscription", attributes: managedDatabaseLogSubscriptionDataAttributes})
	NewCloudManagedDatabaseLogKindDataSource, _ = newManagedDatabaseObjectDataSources[managedDatabaseLogKind](managedDatabaseObjectKind{
		category: managedDatabaseFamily, engines: managedDatabaseAllEngines, suffix: "_log_kind", path: "/log/kind", object: "log kind", key: "name", attributes: managedDatabaseLogKindDataAttributes})

	NewCloudManagedAnalyticsUserDataSource, NewCloudManagedAnalyticsUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseUser](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: []string{"grafana", "kafka", "kafkaConnect"}, suffix: "_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseUserDataAttributes})
	NewCloudManagedAnalyticsClickhouseUserDataSource, NewCloudManagedAnalyticsClickhouseUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseUserWithRoles](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: []string{"clickhouse"}, suffix: "_clickhouse_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseRolesUserDataAttributes})
	NewCloudManagedAnalyticsOpensearchUserDataSource, NewCloudManagedAnalyticsOpensearchUsersDataSource = newManagedDatabaseObjectDataSources[managedDatabaseOpensearchUser](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: []string{"opensearch"}, suffix: "_opensearch_user", path: "/user", object: "user", objects: "users", attributes: managedDatabaseOpensearchUserDataAttributes})
	NewCloudManagedAnalyticsClickhouseDatabaseDataSource, NewCloudManagedAnalyticsClickhouseDatabasesDataSource = newManagedDatabaseObjectDataSources[managedDatabaseDatabase](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: []string{"clickhouse"}, suffix: "_clickhouse_database", path: "/database", object: "database", objects: "databases", attributes: managedDatabaseDatabaseDataAttributes})
	NewCloudManagedAnalyticsClickhouseRoleDataSource, NewCloudManagedAnalyticsClickhouseRolesDataSource = newManagedDatabaseObjectDataSources[managedDatabaseRole](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: []string{"clickhouse"}, suffix: "_clickhouse_role", path: "/roles", object: "role", objects: "roles", attributes: managedDatabaseRoleDataAttributes})
	NewCloudManagedAnalyticsLogSubscriptionDataSource, _ = newManagedDatabaseObjectDataSources[managedDatabaseLogSubscription](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: managedAnalyticsAllEngines, suffix: "_log_subscription", path: "/log/subscription", object: "log subscription", attributes: managedDatabaseLogSubscriptionDataAttributes})
	NewCloudManagedAnalyticsLogKindDataSource, _ = newManagedDatabaseObjectDataSources[managedDatabaseLogKind](managedDatabaseObjectKind{
		category: managedAnalyticsFamily, engines: managedAnalyticsAllEngines, suffix: "_log_kind", path: "/log/kind", object: "log kind", key: "name", attributes: managedDatabaseLogKindDataAttributes})
)

func NewCloudManagedDatabaseLogSubscriptionsDataSource() datasource.DataSource {
	return &managedDatabaseNamesDataSource{managedDatabaseObjectKind: managedDatabaseObjectKind{category: managedDatabaseFamily, engines: managedDatabaseAllEngines, suffix: "_log_subscriptions", path: "/log/subscription", objects: "IDs of the log subscriptions"}, list: "subscription_ids", filter: "kind"}
}

func NewCloudManagedDatabaseLogKindsDataSource() datasource.DataSource {
	return &managedDatabaseNamesDataSource{managedDatabaseObjectKind: managedDatabaseObjectKind{category: managedDatabaseFamily, engines: managedDatabaseAllEngines, suffix: "_log_kinds", path: "/log/kind", objects: "names of the log kinds"}, list: "names"}
}

func NewCloudManagedAnalyticsLogSubscriptionsDataSource() datasource.DataSource {
	return &managedDatabaseNamesDataSource{managedDatabaseObjectKind: managedDatabaseObjectKind{category: managedAnalyticsFamily, engines: managedAnalyticsAllEngines, suffix: "_log_subscriptions", path: "/log/subscription", objects: "IDs of the log subscriptions"}, list: "subscription_ids", filter: "kind"}
}

func NewCloudManagedAnalyticsLogKindsDataSource() datasource.DataSource {
	return &managedDatabaseNamesDataSource{managedDatabaseObjectKind: managedDatabaseObjectKind{category: managedAnalyticsFamily, engines: managedAnalyticsAllEngines, suffix: "_log_kinds", path: "/log/kind", objects: "names of the log kinds"}, list: "names"}
}

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

func managedDatabaseDatabaseDataAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":         managedDatabaseComputedString("ID of the database."),
		"name":       managedDatabaseComputedString("Name of the database."),
		"default":    schema.BoolAttribute{Computed: true, Description: "Whether the database is the default one, created with the service."},
		"checksum":   managedDatabaseComputedString("Hash of the database, computed by the API."),
		"created_at": managedDatabaseComputedString("Creation date of the database."),
		"updated_at": managedDatabaseComputedString("Last update date of the database."),
	}
}

func managedDatabaseConnectionPoolDataAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":          managedDatabaseComputedString("ID of the connection pool."),
		"name":        managedDatabaseComputedString("Name of the connection pool."),
		"database_id": managedDatabaseComputedString("ID of the database the pool connects to."),
		"mode":        managedDatabaseComputedString("PgBouncer pooling mode: SESSION, STATEMENT or TRANSACTION."),
		"size":        schema.Int64Attribute{Computed: true, Description: "Maximum number of server connections in the pool."},
		"user_id":     managedDatabaseComputedString("ID of the user the pool connects with. Null when the pool uses a system user."),
		"port":        schema.Int64Attribute{Computed: true, Description: "Port the connection pool listens on."},
		"ssl_mode":    managedDatabaseComputedString("TLS mode of the connection pool."),
		"uri":         managedDatabaseComputedString("Connection URI of the connection pool."),
		"checksum":    managedDatabaseComputedString("Hash of the connection pool, computed by the API."),
		"created_at":  managedDatabaseComputedString("Creation date of the connection pool."),
		"updated_at":  managedDatabaseComputedString("Last update date of the connection pool."),
	}
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

func managedDatabaseLogSubscriptionDataAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":               managedDatabaseComputedString("ID of the subscription."),
		"kind":             managedDatabaseComputedString("Kind of logs forwarded."),
		"stream_id":        managedDatabaseComputedString("ID of the Logs Data Platform stream receiving the logs."),
		"ldp_service_name": managedDatabaseComputedString("Name of the Logs Data Platform service of the stream."),
		"resource": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Resource the logs come from.",
			Attributes: map[string]schema.Attribute{
				"name": managedDatabaseComputedString("Name of the resource."),
				"type": managedDatabaseComputedString("Type of the resource."),
			},
		},
		"created_at": managedDatabaseComputedString("Creation date of the subscription."),
		"updated_at": managedDatabaseComputedString("Last update date of the subscription."),
	}
}

func managedDatabaseLogKindDataAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name":                       managedDatabaseComputedString("Name of the log kind."),
		"kind_id":                    managedDatabaseComputedString("ID of the log kind."),
		"display_name":               managedDatabaseComputedString("Display name of the log kind."),
		"additional_returned_fields": schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "Log fields the log kind adds."},
		"created_at":                 managedDatabaseComputedString("Creation date of the log kind."),
		"updated_at":                 managedDatabaseComputedString("Last update date of the log kind."),
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

// managedDatabaseObjectDataSource reads one object of a service, found by its key: its id, or the name
// of a log kind.
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
	key := "ID"
	if d.key == "name" {
		key = "Name"
	}
	attributes[d.key] = schema.StringAttribute{Required: true, Description: key + " of the " + d.object + "."}
	resp.Schema = schema.Schema{
		Description: d.describe("get a " + d.object),
		Attributes:  attributes,
	}
}

func (d *managedDatabaseObjectDataSource[T]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	service := d.serviceEndpoint(ctx, req.Config, resp)
	var key types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(d.key), &key)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := service + d.path + "/" + url.PathEscape(key.ValueString())
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

// managedDatabaseLogSubscriptionsOfKind is the query that lists the subscriptions of a kind of logs. It
// takes the publicCloud.database.KindEnum value (CUSTOMER_LOGS), while the subscriptions and the log
// kinds give the name (customer_logs).
func managedDatabaseLogSubscriptionsOfKind(kind string) string {
	return "?kind=" + url.QueryEscape(strings.ToUpper(kind))
}

// managedDatabaseNamesDataSource lists what the API returns as plain strings: the IDs of the log
// subscriptions, which their kind can filter, and the names of the log kinds.
type managedDatabaseNamesDataSource struct {
	managedDatabaseDataSourceConfigured
	managedDatabaseObjectKind
	list   string
	filter string
}

func (d *managedDatabaseNamesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.category.typeName(d.suffix)
}

func (d *managedDatabaseNamesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := d.serviceAttributes()
	if d.filter != "" {
		attributes[d.filter] = schema.StringAttribute{Optional: true, Description: "Lists only the subscriptions of this kind of logs."}
	}
	attributes[d.list] = schema.ListAttribute{ElementType: types.StringType, Computed: true, Description: "The " + d.objects + "."}
	resp.Schema = schema.Schema{
		Description: d.describe("list the " + d.objects),
		Attributes:  attributes,
	}
}

func (d *managedDatabaseNamesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	endpoint := d.serviceEndpoint(ctx, req.Config, resp) + d.path
	if d.filter != "" {
		var filter types.String
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(d.filter), &filter)...)
		if !filter.IsNull() {
			endpoint += managedDatabaseLogSubscriptionsOfKind(filter.ValueString())
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	names, err := helpers.GetAllPagesV2[string](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	if names == nil {
		names = []string{}
	}

	list, diags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(d.list), list)...)
}
