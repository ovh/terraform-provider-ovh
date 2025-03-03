package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*ovhcloudConnectsDataSource)(nil)

func NewOvhcloudConnectsDataSource() datasource.DataSource {
	return &ovhcloudConnectsDataSource{}
}

type ovhcloudConnectsDataSource struct {
	config *Config
}

func (d *ovhcloudConnectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ovhcloud_connects"
}

func (d *ovhcloudConnectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ovhcloudConnectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OvhcloudConnectsDataSourceSchema(ctx)
}

func (d *ovhcloudConnectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		occsIDs []string
		data    OvhcloudConnectsModel
	)

	// Retrieve list of occs
	if err := d.config.OVHClient.Get("/ovhCloudConnect", &occsIDs); err != nil {
		resp.Diagnostics.AddError("Error calling Get /ovhCloudConnect", err.Error())
		return
	}

	// Fetch each occ data
	for _, occID := range occsIDs {
		var occData OvhcloudConnectModel
		endpoint := "/ovhCloudConnect/" + url.PathEscape(occID)
		if err := d.config.OVHClient.Get(endpoint, &occData); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Get %s", endpoint),
				err.Error(),
			)
			return
		}

		occData.ServiceName = ovhtypes.NewTfStringValue(occID)

		data.Occs = append(data.Occs, occData)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
