package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedAnalyticsKafkaUserAccessDataSource)(nil)

func NewCloudManagedAnalyticsKafkaUserAccessDataSource() datasource.DataSource {
	return &cloudManagedAnalyticsKafkaUserAccessDataSource{managedDatabaseObjectKind: managedDatabaseObjectKind{category: managedAnalyticsFamily, engines: []string{"kafka"}, suffix: "_kafka_user_access"}}
}

type cloudManagedAnalyticsKafkaUserAccessDataSource struct {
	managedDatabaseDataSourceConfigured
	managedDatabaseObjectKind
}

// managedDatabaseKafkaUserAccess is the publicCloud.database.kafka.UserAccess model.
type managedDatabaseKafkaUserAccess struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func (d *cloudManagedAnalyticsKafkaUserAccessDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.category.typeName(d.suffix)
}

func (d *cloudManagedAnalyticsKafkaUserAccessDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := d.serviceAttributes()
	attributes["user_id"] = schema.StringAttribute{Required: true, Description: "ID of the user."}
	attributes["cert"] = schema.StringAttribute{
		Computed:    true,
		Description: "Client certificate of the user, in PEM format, for the SSL authentication to Kafka.",
	}
	attributes["key"] = schema.StringAttribute{
		Computed:    true,
		Sensitive:   true,
		Description: "Private key of the client certificate, in PEM format.",
	}
	resp.Schema = schema.Schema{
		Description: d.describe("get the client certificate and key of a user"),
		Attributes:  attributes,
	}
}

func (d *cloudManagedAnalyticsKafkaUserAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	service := d.serviceEndpoint(ctx, req.Config, resp)
	var userId types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("user_id"), &userId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := service + "/user/" + url.PathEscape(userId.ValueString()) + "/access"
	var access managedDatabaseKafkaUserAccess
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &access); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cert"), access.Cert)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), access.Key)...)
}
