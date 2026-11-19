package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var _ datasource.DataSourceWithConfigure = (*cloudManagedDatabaseCertificatesDataSource)(nil)

func NewCloudManagedDatabaseCertificatesDataSource() datasource.DataSource {
	return &cloudManagedDatabaseCertificatesDataSource{managedDatabaseObjectKind: managedDatabaseObjectKind{category: managedDatabaseFamily, engines: []string{"mysql", "postgresql"}, suffix: "_certificates"}}
}

func NewCloudManagedAnalyticsCertificatesDataSource() datasource.DataSource {
	return &cloudManagedDatabaseCertificatesDataSource{managedDatabaseObjectKind: managedDatabaseObjectKind{category: managedAnalyticsFamily, engines: []string{"clickhouse", "kafka"}, suffix: "_certificates"}}
}

type cloudManagedDatabaseCertificatesDataSource struct {
	managedDatabaseDataSourceConfigured
	managedDatabaseObjectKind
}

// managedDatabaseCertificates is the publicCloud.database.Certificates model.
type managedDatabaseCertificates struct {
	Ca string `json:"ca"`
}

func (d *cloudManagedDatabaseCertificatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.category.typeName(d.suffix)
}

func (d *cloudManagedDatabaseCertificatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := d.serviceAttributes()
	attributes["ca"] = schema.StringAttribute{
		Computed:    true,
		Description: "CA certificate, in PEM format, that verifies the connections to the service.",
	}
	resp.Schema = schema.Schema{
		Description: d.describe("get the CA certificate"),
		Attributes:  attributes,
	}
}

func (d *cloudManagedDatabaseCertificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	endpoint := d.serviceEndpoint(ctx, req.Config, resp) + "/certificates"
	if resp.Diagnostics.HasError() {
		return
	}

	var certificates managedDatabaseCertificates
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &certificates); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ca"), certificates.Ca)...)
}
