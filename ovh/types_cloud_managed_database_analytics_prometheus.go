package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// managedDatabasePrometheusModel holds the attributes of every Prometheus endpoint. Its ID is the ID of
// the service. The model of each kind of endpoint embeds it, with the fields of its engines.
type managedDatabasePrometheusModel struct {
	ServiceName   types.String `tfsdk:"service_name"`
	Engine        types.String `tfsdk:"engine"`
	ServiceId     types.String `tfsdk:"service_id"`
	PasswordReset types.String `tfsdk:"password_reset"`

	Id       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Timeouts types.Object `tfsdk:"timeouts"`
}

func managedDatabasePrometheusTimeoutsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"create": types.StringType,
		"delete": types.StringType,
	}
}

type managedDatabasePrometheusTarget struct {
	Host string `json:"host" tfsdk:"host"`
	Port int64  `json:"port" tfsdk:"port"`
}

// managedDatabasePrometheusEndpoint is the publicCloud.database.PrometheusEndpoint model, and the
// publicCloud.database.mongodb.PrometheusEndpoint one, which has a SRV domain instead of targets.
type managedDatabasePrometheusEndpoint struct {
	SrvDomain string                            `json:"srvDomain"`
	Targets   []managedDatabasePrometheusTarget `json:"targets"`
	Username  string                            `json:"username"`
}

// managedDatabasePrometheusAccess is the publicCloud.database.PrometheusAccess model, returned by a
// credentials reset only: no other call returns the password.
type managedDatabasePrometheusAccess struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func managedDatabasePrometheusTargetAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"host": types.StringType,
		"port": types.Int64Type,
	}
}

// CloudManagedDatabasePrometheusModel is an endpoint with scraping targets.
type CloudManagedDatabasePrometheusModel struct {
	managedDatabasePrometheusModel
	Targets types.List `tfsdk:"targets"`
}

func (m *CloudManagedDatabasePrometheusModel) prometheus() *managedDatabasePrometheusModel {
	return &m.managedDatabasePrometheusModel
}

func (m *CloudManagedDatabasePrometheusModel) mergeEndpoint(ctx context.Context, endpoint *managedDatabasePrometheusEndpoint) diag.Diagnostics {
	var diags diag.Diagnostics
	m.Targets, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: managedDatabasePrometheusTargetAttrTypes()}, endpoint.Targets)
	return diags
}

// cloudManagedDatabaseMongodbPrometheusModel is the endpoint of MongoDB, found through a SRV domain.
type cloudManagedDatabaseMongodbPrometheusModel struct {
	managedDatabasePrometheusModel
	SrvDomain types.String `tfsdk:"srv_domain"`
}

func (m *cloudManagedDatabaseMongodbPrometheusModel) prometheus() *managedDatabasePrometheusModel {
	return &m.managedDatabasePrometheusModel
}

func (m *cloudManagedDatabaseMongodbPrometheusModel) mergeEndpoint(ctx context.Context, endpoint *managedDatabasePrometheusEndpoint) diag.Diagnostics {
	m.SrvDomain = types.StringValue(endpoint.SrvDomain)
	return nil
}

// mergeWith copies the endpoint into the attributes every endpoint has, whose password it leaves as is.
func (m *managedDatabasePrometheusModel) mergeWith(endpoint *managedDatabasePrometheusEndpoint) {
	m.Id = m.ServiceId
	m.Username = types.StringValue(endpoint.Username)
}
